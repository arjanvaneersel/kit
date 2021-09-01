package envprovider_test

import (
	"os"
	"testing"

	"github.com/arjanvaneersel/kit/cfg"
	env "github.com/arjanvaneersel/kit/cfg/providers/env"
	th "github.com/arjanvaneersel/kit/cfg/testhelpers"
)

func TestEnvProvider(t *testing.T) {
	// Set values for testing
	os.Setenv("CFGTEST_STR", th.StrVal)
	os.Setenv("CFGTEST_INT", th.IntStr)
	os.Setenv("CFGTEST_FLOAT", th.FloatStr)
	os.Setenv("CFGTEST_BOOL", th.BoolStr)
	os.Setenv("CFGTEST_URL", th.UrlStr)
	os.Setenv("CFGTEST_DURATION", th.DurationStr)
	os.Setenv("CFGTEST_TIME", th.TimeStr)
	os.Setenv("CFGTEST_SLICE", th.SliceStr)

	// Remove test values after performing the test
	defer func() {
		os.Unsetenv("CFGTEST_STR")
		os.Unsetenv("CFGTEST_INT")
		os.Unsetenv("CFGTEST_FLOAT")
		os.Unsetenv("CFGTEST_BOOL")
		os.Unsetenv("CFGTEST_URL")
		os.Unsetenv("CFGTEST_DURATION")
		os.Unsetenv("CFGTEST_TIME")
		os.Unsetenv("CFGTEST_SLICE")
	}()

	p := env.EnvProvider{Prefix: "CFGTEST"}
	cfg, err := cfg.Parse(p)
	if err != nil {
		t.Fatalf("Expected to pass, but got error: %v\n", err)
	}

	th.TestConfig(cfg, t)
}
