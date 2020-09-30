package configurator

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type example struct {
	Name string   `json:"name" yaml:"name"`
	Age  *int64   `json:"age" yaml:"age"`
	Tags []string `json:"tags" yaml:"tags"`
}

func TestFileProvider(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		content string
		expect  *example
	}{
		{
			name:    "json file",
			file:    "*.json",
			content: `{"name":"Tom","age":24,"tags":["foo","bar","baz"]}`,
			expect: &example{
				Name: "Tom",
				Age:  int64ptr(int64(24)),
				Tags: []string{"foo", "bar", "baz"},
			},
		},
		{
			name:    "yaml file",
			file:    "*.yaml",
			content: `{"name":"Tom","tags":["foo","bar","baz"]}`,
			expect: &example{
				Name: "Tom",
				Tags: []string{"foo", "bar", "baz"},
			},
		},
		{
			name:    "yml file",
			file:    "*.yml",
			content: `{"name":"Tom"}`,
			expect: &example{
				Name: "Tom",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfg := example{}
			f, err := ioutil.TempFile("", tt.file)
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.WriteString(tt.content)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())

			fp := NewFileProvider(f.Name())
			err = fp.Provide(&cfg, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, &cfg)
		})
	}
}

func TestFileLoader_PathError(t *testing.T) {
	var cfg example
	fp := NewFileProvider("not_exist_file.json")
	err := fp.Provide(&cfg, nil)
	assert.Error(t, err)
	var pathErr *os.PathError
	assert.True(t, errors.As(err, &pathErr))
}

func TestFileLoader_UnsupportError(t *testing.T) {
	var cfg example
	f, err := ioutil.TempFile("", "*.ini")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	fp := NewFileProvider(f.Name())

	err = fp.Provide(&cfg, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnsupported))
}

func int64ptr(i int64) *int64 {
	return &i
}
