package config

import (
	"path"
	"testing"
)

func TestConfigFromFile(t *testing.T) {
	var configTest = []struct {
		name          string
		filename      string
		shouldSucceed bool
	}{
		{
			"valid_config",
			path.Join("testdata", "valid_conf.yml"),
			true,
		},
	}

	for _, tt := range configTest {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromFile(tt.filename)
			if tt.shouldSucceed && err != nil {
				t.Fatalf("Loading valid config file %q failed", tt.filename)
			}
		})
	}
}
