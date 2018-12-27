package cfg

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ErrKeyNotFound struct {
	Key string
}

func (err ErrKeyNotFound) Error() string {
	return fmt.Sprintf("%s: key not found", err.Key)
}

func IsNotFoundErr(err error) bool {
	_, ok := err.(ErrKeyNotFound)
	return !ok
}

var (
	ErrInvalidCommaString = errors.New("Invalid comma separated string")
)

// Config is a go routine safe configuration store structure which can be accessed via providers
type Config struct {
	mu sync.RWMutex
	v  map[string]string
}

// NewConfig returns a pointer to an initialised Config
func NewConfig() *Config {
	return &Config{v: make(map[string]string)}
}

// Provider is an interface to provide the corresponding configuration as a map
type Provider interface {
	Provide() (map[string]string, error)
}

// Parse parses the configuration from a provider and returns a pointer to a configuration store
func Parse(p Provider) (*Config, error) {
	v, err := p.Provide()
	if err != nil {
		return nil, err
	}

	return &Config{v: v}, nil
}

// Map returns the map with configuration values
func (c *Config) Map() map[string]string {
	return c.v
}

// Len returns the length of the value map
func (c *Config) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.v)
}

// String returns the configuration map as a string
func (c *Config) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var buf bytes.Buffer
	for key, val := range c.v {
		buf.WriteString(key + ": " + val + "\n")
	}

	return buf.String()
}

// GetString gets the requested value from the configuration map and returns it as a string.
// Returns an error if the key can't be found.
func (c *Config) GetString(k string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.v[k]
	if !ok {
		return "", ErrKeyNotFound{k}
	}

	return v, nil
}

// MustString gets the requested value from the configuration map and returns it as a string.
// Will panic if the key can't be found.
func (c *Config) MustString(k string) string {
	v, err := c.GetString(k)
	if err != nil {
		panic(err)
	}
	return v
}

// SetString updates the configuration map with the provided value for the provided key
func (c *Config) SetString(k, v string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.v[k] = v
}

// GetInt gets the requested value from the configuration map and returns it as an integer.
// Returns an error if the key can't be found.
func (c *Config) GetInt(k string) (int, error) {
	v, err := c.GetString(k)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// MustInt gets the requested value from the configuration map and returns it as an integer.
// Panics if the key can't be found.
func (c *Config) MustInt(k string) int {
	i, err := c.GetInt(k)
	if err != nil {
		panic(err)
	}
	return i
}

// SetInt updates the configuration map with the provided value for the provided key
func (c *Config) SetInt(k string, i int) {
	c.SetString(k, strconv.Itoa(i))
}

// GetFloat gets the requested value from the configuration map and returns it as a float64.
// Returns an error if the key can't be found.
func (c *Config) GetFloat(k string) (float64, error) {
	v, err := c.GetString(k)
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

// MustFloat gets the requested value from the configuration map and returns it as a float64.
// Panics if the key can't be found.
func (c *Config) MustFloat(k string) float64 {
	i, err := c.GetFloat(k)
	if err != nil {
		panic(err)
	}
	return i
}

// SetFloat updates the configuration map with the provided value for the provided key
func (c *Config) SetFloat(k string, f float64) {
	c.SetString(k, strconv.FormatFloat(f, 'E', -1, 64))
}

// GetBool gets the requested value from the configuration map and returns it as a boolean.
// Returns an error if the key can't be found.
func (c *Config) GetBool(k string) (bool, error) {
	v, err := c.GetString(k)
	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, err
	}

	return b, nil
}

// MustBool gets the requested value from the configuration map and returns it as a boolean.
// Panics if the key can't be found.
func (c *Config) MustBool(k string) bool {
	b, err := c.GetBool(k)
	if err != nil {
		panic(err)
	}
	return b
}

// SetBool updates the configuration map with the provided value for the provided key
func (c *Config) SetBool(k string, b bool) {
	c.SetString(k, strconv.FormatBool(b))
}

// GetURL gets the requested value from the configuration map and returns it as a pointer to a URL.
// Returns an error if the key can't be found.
func (c *Config) GetURL(k string) (*url.URL, error) {
	v, err := c.GetString(k)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(v)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// MustURL gets the requested value from the configuration map and returns it as a pointer to a URL.
// Panics if the key can't be found.
func (c *Config) MustURL(k string) *url.URL {
	u, err := c.GetURL(k)
	if err != nil {
		panic(err)
	}
	return u
}

// SetURL updates the configuration map with the provided value for the provided key
func (c *Config) SetURL(k string, u *url.URL) {
	c.SetString(k, u.String())
}

// GetDuration gets the requested value from the configuration map and returns it as a Duration.
// Returns an error if the key can't be found.
func (c *Config) GetDuration(k string) (time.Duration, error) {
	v, err := c.GetString(k)
	if err != nil {
		return 0, err
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, err
	}

	return d, nil
}

// MustDuration gets the requested value from the configuration map and returns it as a Duration.
// Panics if the key can't be found.
func (c *Config) MustDuration(k string) time.Duration {
	d, err := c.GetDuration(k)
	if err != nil {
		panic(err)
	}
	return d
}

// SetDuration updates the configuration map with the provided value for the provided key
func (c *Config) SetDuration(k string, d time.Duration) {
	c.SetString(k, d.String())
}

// GetTime gets the requested value from the configuration map and returns it as a Time.
// Returns an error if the key can't be found.
func (c *Config) GetTime(k string) (time.Time, error) {
	v, err := c.GetString(k)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// MustTime gets the requested value from the configuration map and returns it as a Duration.
// Panics if the key can't be found.
func (c *Config) MustTime(k string) time.Time {
	t, err := c.GetTime(k)
	if err != nil {
		panic(err)
	}
	return t
}

// SetTime updates the configuration map with the provided value for the provided key
func (c *Config) SetTime(k string, t time.Time) {
	c.SetString(k, t.Format(time.RFC3339))
}

// GetSlice gets the requested value, which must be a comma separated
// string, from the configuration map and returns it as a slice.
// Returns an error if the key can't be found.
func (c *Config) GetSlice(k string) ([]string, error) {
	v, err := c.GetString(k)
	if err != nil {
		return nil, err
	}

	s := strings.Split(v, ",")
	if len(s) == 0 {
		// If v is a valid, but not comma separated, string
		//then return an array with v as the only element
		if len(v) > 0 {
			return []string{v}, nil
		}
		return nil, ErrInvalidCommaString
	}

	return s, nil
}

// MustSlice gets the requested value, which must be a comma separated
// string, from the configuration map and returns it as a slice.
// Panics if the key can't be found.
func (c *Config) MustSlice(k string) []string {
	s, err := c.GetSlice(k)
	if err != nil {
		panic(err)
	}
	return s
}

// SetSlice updates the configuration map with the provided value for the provided key
func (c *Config) SetSlice(k string, s []string) {
	c.SetString(k, strings.Join(s, ","))
}
