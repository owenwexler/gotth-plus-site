package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	GoEnv string
	Port  string
	// DatabasePath is the SQLite file. Defaults to app.db.
	DatabasePath string
	// TrustedProxies are the networks of the reverse proxies whose X-Forwarded-For header is believed.
	// Empty means the peer address is the client (no proxy in front of the app).
	TrustedProxies []*net.IPNet
}

// IsProduction reports whether the app runs in production mode.
func (e Env) IsProduction() bool {
	return e.GoEnv == "production"
}

func GetEnv() Env {
	// .env is for local development; in production the platform provides the environment, so a
	// missing file is fine. Anything else (e.g. a malformed file) is not.
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(err)
	}

	proxies, err := ParseTrustedProxies(os.Getenv("TrustedProxies"))
	if err != nil {
		panic(err)
	}

	env := Env{
		GoEnv:          os.Getenv("GoEnv"),
		Port:           os.Getenv("Port"),
		DatabasePath:   os.Getenv("DatabasePath"),
		TrustedProxies: proxies,
	}

	// fail fast: an unknown mode must never silently behave like development
	if env.GoEnv != "development" && env.GoEnv != "production" {
		panic(fmt.Sprintf(`GoEnv must be "development" or "production", got %q`, env.GoEnv))
	}

	if env.Port == "" {
		panic("Port must be set")
	}

	if env.DatabasePath == "" {
		env.DatabasePath = "app.db"
	}

	return env
}

// ParseTrustedProxies reads a comma-separated list of CIDR networks or single addresses, such as
// "127.0.0.1/32, 172.18.0.0/16". Empty input means no trusted proxies.
func ParseTrustedProxies(s string) ([]*net.IPNet, error) {
	var nets []*net.IPNet

	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		if !strings.Contains(item, "/") {
			ip := net.ParseIP(item)
			if ip == nil {
				return nil, fmt.Errorf("TrustedProxies: %q is not an IP address or CIDR network", item)
			}
			if ip.To4() != nil {
				item += "/32"
			} else {
				item += "/128"
			}
		}

		_, n, err := net.ParseCIDR(item)
		if err != nil {
			return nil, fmt.Errorf("TrustedProxies: %q: %w", item, err)
		}
		nets = append(nets, n)
	}

	return nets, nil
}
