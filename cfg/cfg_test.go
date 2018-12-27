package cfg_test

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/arjanvaneersel/kit/cfg"
)

var (
	Success = "\u2713"
	Failed  = "\u2717"
)

var (
	STRVAL      string        = "TEST"
	INTVAL      int           = 2
	FLOATVAL    float64       = 1.2345
	BOOLVAL     bool          = true
	urlStr      string        = "http://www.example.com"
	URLVAL, _                 = url.Parse(urlStr)
	DURATIONVAL time.Duration = 5 * time.Second
	timeStr     string        = "2012-11-01T22:08:41+00:00"
	TIMEVAL, _                = time.Parse(time.RFC3339, timeStr)
	SLICEVAL    []string      = []string{"val1", "val2", "val3"}
)

func mockConfig() *cfg.Config {
	// Set default values for testing
	c := cfg.NewConfig()
	c.SetString("STR", STRVAL)
	c.SetInt("INT", INTVAL)
	c.SetFloat("FLOAT", FLOATVAL)
	c.SetBool("BOOL", BOOLVAL)
	c.SetURL("URL", URLVAL)
	c.SetDuration("DURATION", DURATIONVAL)
	c.SetTime("TIME", TIMEVAL)
	c.SetSlice("SLICE", SLICEVAL)

	return c
}

func TestEnvProvider(t *testing.T) {
	// Set values for testing
	os.Setenv("CFGTEST_FOO", "Bar")

	// Remove test values after performing the test
	defer func() {
		os.Unsetenv("CFGTEST_FOO")
	}()

	t.Log("Testing EnvProvider")
	{
		t.Logf("\tTesting Provider setup")

		p := cfg.EnvProvider{"CFGTEST"}
		cfg, err := cfg.Parse(p)
		if err != nil {
			t.Fatalf("\t\t[%s] Expected to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected Provide() to pass\n", Success)

		expLength := 1
		l := cfg.Len()
		if l != expLength {
			t.Fatalf("\t\t[%s] Expected a length of %d, but got %d:\n%s\n", Failed, expLength, l, cfg)
		}
		t.Logf("\t\t[%s] Expected a length of %d\n", Success, expLength)

		//testValues(t, cfg)
	}
}

func TestTxtProvider(t *testing.T) {
	c := mockConfig()

	p := cfg.TxtProvider{"config.txt", ","}
	err := p.Save(c)
	if err != nil {
		t.Fatalf("\t[%s]Expected to be able to write test file, but got error: %v\n", Failed, err)
	}
	t.Logf("\t[%s]Expected to be able to write test file\n", Success)

	// Remove test file after performing the test
	defer func() {
		os.Remove("config.txt")
	}()

	t.Log("Testing TxtProvider")
	{
		t.Logf("\tTesting Provider setup")

		p := cfg.TxtProvider{"config.txt", ","}
		cfg, err := cfg.Parse(p)
		if err != nil {
			t.Fatalf("\t\t[%s] Expected to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected Provide() to pass\n", Success)

		expLen := c.Len()
		l := cfg.Len()
		if l != expLen {
			t.Fatalf("\t\t[%s] Expected a length of %d, but got %d:\n%s\n", Failed, expLen, l, cfg)
		}
		t.Logf("\t\t[%s] Expected a length of %d\n", Success, expLen)

		testValues(t, cfg)
	}
}

func TestGobProvider(t *testing.T) {
	c := mockConfig()

	p := cfg.GobProvider{"config.gob"}
	err := p.Save(c)
	if err != nil {
		t.Fatalf("\t[%s]Expected to be able to write test file, but got error: %v\n", Failed, err)
	}
	t.Logf("\t[%s]Expected to be able to write test file\n", Success)

	// Remove test file after performing the test
	defer func() {
		os.Remove("config.gob")
	}()

	t.Log("Testing GobProvider")
	{
		t.Logf("\tTesting Provider setup")

		p := cfg.GobProvider{"config.gob"}
		cfg, err := cfg.Parse(p)
		if err != nil {
			t.Fatalf("\t\t[%s] Expected to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected Provide() to pass\n", Success)

		expLen := c.Len()
		l := cfg.Len()
		if l != expLen {
			t.Fatalf("\t\t[%s] Expected a length of %d, but got %d:\n%s\n", Failed, expLen, l, cfg)
		}
		t.Logf("\t\t[%s] Expected a length of %d\n", Success, expLen)

		testValues(t, cfg)
	}
}

func TestJSONProvider(t *testing.T) {
	c := mockConfig()

	p := cfg.JSONProvider{"config.json"}
	err := p.Save(c)
	if err != nil {
		t.Fatalf("\t[%s]Expected to be able to write test file, but got error: %v\n", Failed, err)
	}
	t.Logf("\t[%s]Expected to be able to write test file\n", Success)

	// Remove test file after performing the test
	/*defer func() {
		os.Remove("config.json")
	}()*/

	t.Log("Testing JSONProvider")
	{
		t.Logf("\tTesting Provider setup")

		p := cfg.JSONProvider{"config.json"}
		cfg, err := cfg.Parse(p)
		if err != nil {
			t.Fatalf("\t\t[%s] Expected to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected Provide() to pass\n", Success)

		expLen := c.Len()
		l := cfg.Len()
		if l != expLen {
			t.Fatalf("\t\t[%s] Expected a length of %d, but got %d:\n%s\n", Failed, expLen, l, cfg)
		}
		t.Logf("\t\t[%s] Expected a length of %d\n", Success, expLen)

		testValues(t, cfg)
	}
}

func testValues(t *testing.T, cfg *cfg.Config) {
	t.Logf("\tTesting STRING")
	{
		str, err := cfg.GetString("STR")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetString to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetString to pass\n", Success)

		if str != STRVAL {
			t.Fatalf("\t\t[%s] Expected a GetString to return \"Test\", but got: %s\n", Failed, str)
		}
		t.Logf("\t\t[%s] Expected a GetString to return \"Test\"\n", Success)
	}

	t.Logf("\tTesting INT")
	{
		i, err := cfg.GetInt("INT")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetInt to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetInt to pass\n", Success)

		if i != INTVAL {
			t.Fatalf("\t\t[%s] Expected a GetInt to return 2, but got error: %d\n", Failed, i)
		}
		t.Logf("\t\t[%s] Expected a GetInt to return 2\n", Success)
	}

	t.Logf("\tTesting FLOAT")
	{
		f, err := cfg.GetFloat("FLOAT")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetFloat to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetFloat to pass\n", Success)

		if f != FLOATVAL {
			t.Fatalf("\t\t[%s] Expected a GetFloat to return 1.2345, but got: %f\n", Failed, f)
		}
		t.Logf("\t\t[%s] Expected a GetFloat to return 1.2345\n", Success)
	}

	t.Logf("\tTesting BOOL")
	{
		b, err := cfg.GetBool("BOOL")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetBool to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetBool to pass\n", Success)

		if b != BOOLVAL {
			t.Fatalf("\t\t[%s] Expected a GetBool to return true, but got: %v\n", Failed, b)
		}
		t.Logf("\t\t[%s] Expected a GetBool to return true\n", Success)
	}

	t.Logf("\tTesting URL")
	{
		u, err := cfg.GetURL("URL")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetURL to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetURL to pass\n", Success)

		if u.String() != urlStr {
			t.Fatalf("\t\t[%s] Expected a GetURL to return \"http://www.example.com\", but got: %v\n", Failed, u)
		}
		t.Logf("\t\t[%s] Expected a GetBoolURL to return \"http://www.example.com\"\n", Success)
	}

	t.Logf("\tTesting DURATION")
	{
		d, err := cfg.GetDuration("DURATION")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetDuration to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetDuration to pass\n", Success)

		if d != DURATIONVAL {
			t.Fatalf("\t\t[%s] Expected GetDuration to return %v, but got: %v\n", Failed, DURATIONVAL, d)
		}
		t.Logf("\t\t[%s] Expected GetDuration to return %v\n", Success, DURATIONVAL)
	}

	t.Logf("\tTesting TIME")
	{
		tme, err := cfg.GetTime("TIME")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetTime to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetTime to pass\n", Success)

		if tme.Format(time.RFC3339) != TIMEVAL.Format(time.RFC3339) {
			t.Fatalf("\t\t[%s] Expected GetTime to return %s, but got: %s\n", Failed, TIMEVAL.Format(time.RFC3339), tme.Format(time.RFC3339))
		}
		t.Logf("\t\t[%s] Expected GetTime to return %v\n", Success, TIMEVAL)
	}

	t.Logf("\tTesting SLICE")
	{
		s, err := cfg.GetSlice("SLICE")
		if err != nil {
			t.Fatalf("\t\t[%s] Expected a GetSlice to pass, but got error: %v\n", Failed, err)
		}
		t.Logf("\t\t[%s] Expected a GetSlice to pass\n", Success)

		if len(s) != 3 || s[0] != SLICEVAL[0] {
			t.Fatalf("\t\t[%s] Expected GetSlice to return %v, but got: %v\n", Failed, SLICEVAL, s)
		}
		t.Logf("\t\t[%s] Expected GetSlice to return %v\n", Success, SLICEVAL)
	}
}
