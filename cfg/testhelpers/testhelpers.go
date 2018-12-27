package testhelpers

import (
	"net/url"
	"testing"
	"time"

	"github.com/arjanvaneersel/kit/cfg"
)

var (
	StrVal      string        = "TEST"
	IntStr      string        = "2"
	IntVal      int           = 2
	FloatStr    string        = "1.2345E+00"
	FloatVal    float64       = 1.2345
	BoolStr     string        = "true"
	BoolVal     bool          = true
	UrlStr      string        = "http://www.example.com"
	UrlVal, _                 = url.Parse(UrlStr)
	DurationStr string        = "5s"
	DurationVal time.Duration = 5 * time.Second
	TimeStr     string        = "2012-11-01T22:08:41+00:00"
	TimeVal, _                = time.Parse(time.RFC3339, TimeStr)
	SliceStr    string        = "val1,val2,val3"
	SliceVal    []string      = []string{"val1", "val2", "val3"}
)

func MockConfig() *cfg.Config {
	// Set default values for testing
	c := cfg.NewConfig()
	c.SetString("STR", StrVal)
	c.SetInt("INT", IntVal)
	c.SetFloat("FLOAT", FloatVal)
	c.SetBool("BOOL", BoolVal)
	c.SetURL("URL", UrlVal)
	c.SetDuration("DURATION", DurationVal)
	c.SetTime("TIME", TimeVal)
	c.SetSlice("SLICE", SliceVal)

	return c
}

func TestConfig(c *cfg.Config, t *testing.T) {
	t.Run("STR", func(t *testing.T) {
		expected := StrVal
		got, err := c.GetString("STR")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got != expected {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("INT", func(t *testing.T) {
		expected := IntVal
		got, err := c.GetInt("INT")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got != expected {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("FLOAT", func(t *testing.T) {
		expected := FloatVal
		got, err := c.GetFloat("FLOAT")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got != expected {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("BOOL", func(t *testing.T) {
		expected := BoolVal
		got, err := c.GetBool("BOOL")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got != expected {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("URL", func(t *testing.T) {
		expected := UrlVal
		got, err := c.GetURL("URL")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got.String() != expected.String() {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("DURATION", func(t *testing.T) {
		expected := DurationVal
		got, err := c.GetDuration("DURATION")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got != expected {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("TIME", func(t *testing.T) {
		expected := TimeVal
		got, err := c.GetTime("TIME")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if got.String() != expected.String() {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})

	t.Run("SLICE", func(t *testing.T) {
		expected := SliceVal
		got, err := c.GetSlice("SLICE")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		equal := true
		for _, e := range expected {
			found := false
			for _, i := range got {
				if i == e {
					found = true
					break
				}
			}
			if !found {
				equal = false
				break
			}
		}
		if !equal {
			t.Errorf("expected %v, but got %v", expected, got)
		}
	})
}
