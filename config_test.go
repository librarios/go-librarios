package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestLoadConfigBytes(t *testing.T) {
	data := `
port: 80
plugins:
  foo:
    enabled: true
    int-value: 10
  bar:
    enabled: false
    float-value: 10.2
  foobar:
    enabled: false
`
	config, err := LoadConfigBytes([]byte(data))
	if err != nil {
		t.Error(err)
	}
	expected := &Config{
		Port: 80,
		Plugins: map[string]Map{
			"foo": Map{
				"enabled":   true,
				"int-value": 10,
			},
			"bar": Map{
				"enabled":     false,
				"float-value": 10.2,
			},
			"foobar": Map{
				"enabled": false,
			},
		},
	}

	if diff := cmp.Diff(expected, config); diff != "" {
		t.Errorf("TestLoadConfigBytes() mismatch (-expected +actual):\n%s", diff)
	}
}
