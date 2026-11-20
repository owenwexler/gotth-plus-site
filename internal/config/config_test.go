package config

import "testing"

func TestIsProduction(t *testing.T) {
	for goEnv, want := range map[string]bool{"production": true, "development": false, "": false, "Production": false} {
		if got := (Env{GoEnv: goEnv}).IsProduction(); got != want {
			t.Errorf("GoEnv=%q: IsProduction() = %v, want %v", goEnv, got, want)
		}
	}
}

func TestParseTrustedProxies(t *testing.T) {
	nets, err := ParseTrustedProxies(" 127.0.0.1 , 172.18.0.0/16,, ::1 ")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"127.0.0.1/32", "172.18.0.0/16", "::1/128"}
	if len(nets) != len(want) {
		t.Fatalf("got %d networks, want %d", len(nets), len(want))
	}
	for i, n := range nets {
		if n.String() != want[i] {
			t.Errorf("network %d = %s, want %s", i, n, want[i])
		}
	}

	if nets, err := ParseTrustedProxies(""); err != nil || len(nets) != 0 {
		t.Errorf("empty: got %v, %v; want no networks and no error", nets, err)
	}

	for _, bad := range []string{"localhost", "10.0.0.0/33", "not/a/net"} {
		if _, err := ParseTrustedProxies(bad); err == nil {
			t.Errorf("ParseTrustedProxies(%q): want an error", bad)
		}
	}
}
