package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {}

func requireEnv(names ...string) (map[string]string, error) {
	vals := make(map[string]string, len(names))
	var missing []string
	for _, name := range names {
		v := os.Getenv(name)
		if v == "" {
			missing = append(missing, name)
		} else {
			vals[name] = v
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("bipline: missing required env vars: %s", strings.Join(missing, ", "))
	}
	return vals, nil
}
