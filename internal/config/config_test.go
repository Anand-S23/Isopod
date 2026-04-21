package config

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "Returns set value",
			key:          "TEST_VAR",
			value:        "high_performance",
			defaultValue: "default",
			expected:     "high_performance",
		},
		{
			name:         "Returns default when unset",
			key:          "UNSET_VAR",
			value:        "",
			defaultValue: "brutalist_default",
			expected:     "brutalist_default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	t.Run("Production Mode", func(t *testing.T) {
		os.Setenv("MODE", "production")
		os.Setenv("PORT", "443")
		os.Setenv("DB_URI", "postgres://localhost:5432/test")
		defer os.Clearenv()

		cfg, _ := LoadEnv()

		if !cfg.PRODUCTION {
			t.Error("Expected PRODUCTION to be true")
		}
		if cfg.PORT != "443" {
			t.Errorf("Expected PORT 443, got %s", cfg.PORT)
		}
	})

	t.Run("Development Default", func(t *testing.T) {
		os.Unsetenv("MODE")
		os.Unsetenv("PORT")
		os.Setenv("DB_URI", "postgres://localhost:5432/test")
		defer os.Unsetenv("DB_URI")

		cfg, _ := LoadEnv()

		if cfg.PRODUCTION {
			t.Error("Expected PRODUCTION to be false")
		}
		if cfg.PORT != "8080" {
			t.Errorf("Expected default PORT 8080, got %s", cfg.PORT)
		}
		if cfg.DB_URI != "postgres://localhost:5432/test" {
			t.Errorf("DB_URI = %q", cfg.DB_URI)
		}
	})

	t.Run("DB_URI from env", func(t *testing.T) {
		os.Setenv("DB_URI", "postgres://localhost:5432/app")
		defer os.Unsetenv("DB_URI")

		cfg, err := LoadEnv()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.DB_URI != "postgres://localhost:5432/app" {
			t.Errorf("DB_URI = %q", cfg.DB_URI)
		}
	})
}

func TestGetEnvOrPanic_Fails(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		GetEnvOrPanic("NON_EXISTENT_VAR")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestGetEnvOrPanic_Fails$", "-test.count=1")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected subprocess to exit non-zero; output:\n%s", out)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("unexpected error: %v\noutput:\n%s", err, out)
	}
	if code := exitErr.ExitCode(); code != 1 {
		t.Fatalf("expected exit code 1 from log.Fatal, got %d\noutput:\n%s", code, out)
	}
}
