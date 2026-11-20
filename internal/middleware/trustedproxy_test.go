package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mustCIDR(t *testing.T, s string) *net.IPNet {
	t.Helper()
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestTrustedProxy(t *testing.T) {
	trusted := []*net.IPNet{mustCIDR(t, "10.0.0.0/8")}

	seen := ""
	h := TrustedProxy(trusted)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.RemoteAddr
	}))

	for _, tt := range []struct {
		name, peer, xff, want string
	}{
		{"peer is not a proxy: header ignored", "203.0.113.9:4000", "198.51.100.1", "203.0.113.9:4000"},
		{"proxy with one client", "10.0.0.2:4000", "198.51.100.1", "198.51.100.1:0"},
		{"right-most untrusted wins over a forged left entry", "10.0.0.2:4000", "1.2.3.4, 198.51.100.1", "198.51.100.1:0"},
		{"chained trusted proxies are skipped", "10.0.0.2:4000", "198.51.100.1, 10.0.0.3", "198.51.100.1:0"},
		{"garbage in the header keeps the peer", "10.0.0.2:4000", "not-an-ip", "10.0.0.2:4000"},
		{"no header keeps the peer", "10.0.0.2:4000", "", "10.0.0.2:4000"},
		{"all trusted takes the left-most", "10.0.0.2:4000", "10.0.0.5, 10.0.0.6", "10.0.0.5:0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.peer
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}
			h.ServeHTTP(httptest.NewRecorder(), req)

			if seen != tt.want {
				t.Errorf("RemoteAddr = %q, want %q", seen, tt.want)
			}
		})
	}
}

func TestTrustedProxyWithoutNetworksIsANoOp(t *testing.T) {
	h := TrustedProxy(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RemoteAddr != "10.0.0.2:4000" {
			t.Errorf("RemoteAddr changed to %q", r.RemoteAddr)
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.2:4000"
	req.Header.Set("X-Forwarded-For", "198.51.100.1")
	h.ServeHTTP(httptest.NewRecorder(), req)
}
