package main

import (
	"strings"
	"testing"
)

func TestRequireEnv(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		vars        []string
		wantErr     bool
		errContains string
	}{
		{
			name: "all present",
			env:  map[string]string{"FOO": "foo", "BAR": "bar"},
			vars: []string{"FOO", "BAR"},
		},
		{
			name:        "one missing",
			env:         map[string]string{"FOO": "foo"},
			vars:        []string{"FOO", "BAR"},
			wantErr:     true,
			errContains: "BAR",
		},
		{
			name:        "all missing",
			env:         map[string]string{},
			vars:        []string{"FOO", "BAR"},
			wantErr:     true,
			errContains: "FOO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got, err := requireEnv(tt.vars...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, name := range tt.vars {
				if got[name] != tt.env[name] {
					t.Errorf("got[%q] = %q, want %q", name, got[name], tt.env[name])
				}
			}
		})
	}
}
