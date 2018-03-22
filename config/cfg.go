package config

// https://github.com/ardanlabs/kit

import (
	"sync"
	"bytes"
	"strings"
	"fmt"
	"strconv"
	"time"
	"net/url"
)

type Config struct {
	*sync.RWMutex
	m map[string]string
}

type Provider interface {
	Provide() (map[string]string, error)
}

func New(p Provider) (*Config, error) {
	m, err := p.Provide()
	if err != nil {
		return nil, err
	}

	c := &Config{
		&sync.RWMutex{},
		m,
	}

	return c, nil
}

func (c *Config) Log() string {
	c.RLock()
	defer c.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("\n")
	for k, v := range c.m {
		if strings.Contains(k, "pass") {
			buf.WriteString(k + "\t=\t********\n")
		} else {
			buf.WriteString(k + "\t=\t" + v + "\n")
		}
	}

	return buf.String()
}

// String returns the value of the given key as a string. It will return an
// error if key was not found.
func (c *Config) String(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		return "", fmt.Errorf("unknown key %s", key)
	}

	return value, nil
}

// MustString returns the value of the given key as a string. It will panic if
// the key was not found.
func (c *Config) MustString(key string) string {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown key %s !", key))
	}

	return value
}

// SetString adds or modifies the configuration for the specified key and
// value.
func (c *Config) SetString(key string, value string) {
	c.Lock()
	{
		c.m[key] = value
	}
	c.Unlock()
}

// Int returns the value of the given key as an int. It will return an error if
// the key was not found or the value can't be converted to an int.
func (c *Config) Int(key string) (int, error) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		return 0, fmt.Errorf("unknown key %s", key)
	}

	iv, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return iv, nil
}

// MustInt returns the value of the given key as an int. It will panic if the
// key was not found or the value can't be converted to an int.
func (c *Config) MustInt(key string) int {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("Unknown key %s !", key))
	}

	iv, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("key %q value is not an int", key))
	}

	return iv
}

// SetInt adds or modifies the configuration for the specified key and value.
func (c *Config) SetInt(key string, value int) {
	c.Lock()
	{
		c.m[key] = strconv.Itoa(value)
	}
	c.Unlock()
}

// Time returns the value of the given key as a Time. It will return an error
// if the key was not found or the value can't be converted to a Time.
func (c *Config) Time(key string) (time.Time, error) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		return time.Time{}, fmt.Errorf("unknown key %s", key)
	}

	tv, err := time.Parse(time.UnixDate, value)
	if err != nil {
		return tv, err
	}

	return tv, nil
}

// MustTime returns the value of the given key as a Time. It will panic if the
// key was not found or the value can't be converted to a Time.
func (c *Config) MustTime(key string) time.Time {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("unknown key %s", key))
	}

	tv, err := time.Parse(time.UnixDate, value)
	if err != nil {
		panic(fmt.Sprintf("key %q value is not a Time", key))
	}

	return tv
}

// SetTime adds or modifies the configuration for the specified key and value.
func (c *Config) SetTime(key string, value time.Time) {
	c.Lock()
	{
		c.m[key] = value.Format(time.UnixDate)
	}
	c.Unlock()
}

// Bool returns the bool value of a given key as a bool. It will return an
// error if the key was not found or the value can't be converted to a bool.
func (c *Config) Bool(key string) (bool, error) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		return false, fmt.Errorf("unknown key %s", key)
	}

	if value == "on" || value == "yes" {
		value = "true"
	} else if value == "off" || value == "no" {
		value = "false"
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return val, nil
}

// MustBool returns the bool value of a given key as a bool. It will panic if
// the key was not found or the value can't be converted to a bool.
func (c *Config) MustBool(key string) bool {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("unknown key %s", key))
	}

	if value == "on" || value == "yes" {
		value = "true"
	} else if value == "off" || value == "no" {
		value = "false"
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return val
}

// SetBool adds or modifies the configuration for the specified key and value.
func (c *Config) SetBool(key string, value bool) {
	str := "false"
	if value {
		str = "true"
	}

	c.Lock()
	{
		c.m[key] = str
	}
	c.Unlock()
}

// URL returns the value of the given key as a URL. It will return an error if
// the key was not found or the value can't be converted to a URL.
func (c *Config) URL(key string) (*url.URL, error) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		return nil, fmt.Errorf("unknown key %s", key)
	}

	u, err := url.Parse(value)
	if err != nil {
		return u, err
	}

	return u, nil
}

// MustURL returns the value of the given key as a URL. It will panic if the
// key was not found or the value can't be converted to a URL.
func (c *Config) MustURL(key string) *url.URL {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Sprintf("unknown key %s", key))
	}

	u, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("key %q value is not a URL", key))
	}

	return u
}

// SetURL adds or modifies the configuration for the specified key and value.
func (c *Config) SetURL(key string, value *url.URL) {
	c.Lock()
	{
		c.m[key] = value.String()
	}
	c.Unlock()
}

// Duration returns the value of the given key as a Duration. It will return an
// error if the key was not found or the value can't be converted to a Duration.
func (c *Config) Duration(key string) (time.Duration, error) {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		return time.Duration(0), fmt.Errorf("unknown key %s", key)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return d, err
	}

	return d, nil
}

// MustDuration returns the value of the given key as a Duration. It will panic
// if the key was not found or the value can't be converted into a Duration.
func (c *Config) MustDuration(key string) time.Duration {
	c.RLock()
	defer c.RUnlock()

	value, found := c.m[key]
	if !found {
		panic(fmt.Errorf("unknown key %s", key))
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("key %q value is not a Duration", key))
	}

	return d
}

// SetDuration adds or modifies the configuration for a given duration at a
// specific key.
func (c *Config) SetDuration(key string, value time.Duration) {
	c.Lock()
	{
		c.m[key] = value.String()
	}
	c.Unlock()
}
