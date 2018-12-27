package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	tt := []LogLevel{FATAL, ERROR, WARN, INFO, DEBUG}

	for _, tc := range tt {
		for _, ll := range tt {
			var buffer bytes.Buffer
			writer := io.Writer(&buffer)
			l := NewStdLog(tc, log.New(writer, "", 0))
			l.ExitOnFatal = false

			l.Log(ll, "foo")
			got := buffer.String()
			if ll > tc {
				if got != "" {
					t.Errorf("%s:%s expected be empty, but got %s", tc, ll, got)
				}
			} else {
				expected := ll.String()
				if !strings.Contains(got, expected) {
					t.Errorf("%s:%s expected to contain %s, but got %s", tc, ll, expected, got)
				}
			}

		}
	}
}

func TestKeyValuePair(t *testing.T) {
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	logger := NewStdLog(DEBUG, log.New(writer, "", 0))

	logger.Log(DEBUG, "foo", "bar")
	got := buffer.String()
	if expected := "foo: bar"; !strings.Contains(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}

	if expected := "[DEBUG]"; !strings.Contains(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}
}

func TestSingleLine(t *testing.T) {
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	logger := NewStdLog(DEBUG, log.New(writer, "", 0))

	logger.Log(DEBUG, "foo")
	got := buffer.String()
	if expected := "foo"; !strings.Contains(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}

	if expected := "[DEBUG]"; !strings.Contains(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}
}

func TestCombination(t *testing.T) {
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	logger := NewStdLog(DEBUG, log.New(writer, "", 0))

	logger.Log(DEBUG, "foo", "bar", "test")
	got := buffer.String()
	if expected := "foo: bar test"; !strings.Contains(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}

	if expected := "[DEBUG]"; !strings.Contains(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}
}

func TestString(t *testing.T) {
	ll := DEBUG
	expected := "DEBUG"

	if ll.String() != expected {
		t.Errorf("expected %s, but got %s", expected, ll.String())
	}

	got := fmt.Sprintf("%s", ll)
	if got != expected {
		t.Errorf("expected %s, but got %s", expected, got)
	}

}
