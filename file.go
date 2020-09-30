package configurator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func NewFileProvider(filename string) *fileProvider {
	return &fileProvider{filename: filename}
}

type fileProvider struct {
	filename string
}

func (p fileProvider) Provide(v interface{}, _ StructInfo) error {
	f, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	var d decoder
	switch strings.ToLower(filepath.Ext(p.filename)) {
	case ".json":
		d = json.NewDecoder(f)
	case ".yaml", ".yml":
		d = yaml.NewDecoder(f)
	default:
		return fmt.Errorf("the specified file %s is %w", p.filename, ErrUnsupported)
	}
	return d.Decode(v)
}

type decoder interface {
	Decode(interface{}) error
}
