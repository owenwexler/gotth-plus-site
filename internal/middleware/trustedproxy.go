package middleware

import (
	"net"
	"net/http"
	"strings"
)

// TrustedProxy replaces r.RemoteAddr with the client's address when the request comes from one of the
// trusted proxy networks. The proxy appends the address it received the request from to X-Forwarded-For,
// so the right-most entry that is not itself trusted is the client. Entries to its left were written by
// the client and are ignored. Requests from any other peer are left alone, so a forged header changes nothing.
//
// The rate limiter and the request log read RemoteAddr, so they see the client rather than the proxy.
func TrustedProxy(trusted []*net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if len(trusted) == 0 {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if peer := hostIP(r.RemoteAddr); peer != nil && contains(trusted, peer) {
				if client := forwardedClient(r.Header.Get("X-Forwarded-For"), trusted); client != nil {
					r.RemoteAddr = net.JoinHostPort(client.String(), "0")
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// forwardedClient walks X-Forwarded-For from the right and returns the first address outside the trusted
// networks. If every address is trusted, it returns the left-most valid one. It returns nil when the header
// holds nothing usable, so the peer address is kept.
func forwardedClient(header string, trusted []*net.IPNet) net.IP {
	if header == "" {
		return nil
	}

	parts := strings.Split(header, ",")

	var leftmost net.IP
	for i := len(parts) - 1; i >= 0; i-- {
		ip := net.ParseIP(strings.TrimSpace(parts[i]))
		if ip == nil {
			return nil
		}
		if !contains(trusted, ip) {
			return ip
		}
		leftmost = ip
	}

	return leftmost
}

func contains(nets []*net.IPNet, ip net.IP) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func hostIP(addr string) net.IP {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	return net.ParseIP(host)
}
