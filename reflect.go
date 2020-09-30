package configurator

import (
	"fmt"
	"reflect"
	"strings"
)

type StructInfo interface {
	Fields() []FieldInfo
}

type structInfo struct {
	fields []FieldInfo
}

var _ StructInfo = &structInfo{}

func (s *structInfo) Fields() []FieldInfo {
	return s.fields
}

type FieldInfo interface {
	StructField() reflect.StructField
	Value() reflect.Value
	Parent() FieldInfo
	Name() string
	ENVKey() string
	FlagKey() string
	DefVal() string
}

type fieldInfo struct {
	parent *fieldInfo
	field  reflect.StructField
	val    reflect.Value
	tag    tagInfo
}

var _ FieldInfo = &fieldInfo{}

func (f *fieldInfo) StructField() reflect.StructField {
	return f.field
}

func (f *fieldInfo) Value() reflect.Value {
	return f.val
}

func (f *fieldInfo) Parent() FieldInfo {
	return f.parent
}

func (f *fieldInfo) Name() string {
	return f.field.Name
}

func (f *fieldInfo) path() []string {
	path := []string{f.Name()}
	for p := f.parent; p != nil; p = p.parent {
		path = append([]string{p.Name()}, path...)
	}
	return path
}

func (f *fieldInfo) ENVKey() string {
	if f.tag.hasENV {
		if f.tag.env == "" {
			return strings.ToUpper(strings.Join(f.path(), "_"))
		}
		return f.tag.env
	}
	return ""
}

func (f *fieldInfo) FlagKey() string {
	if f.tag.hasFlag {
		if f.tag.flag == "" {
			return strings.ToLower(strings.Join(f.path(), "-"))
		}
		return f.tag.flag
	}
	return ""
}

func (f *fieldInfo) DefVal() string {
	if f.tag.hasDefault {
		return f.tag.defVal
	}
	return ""
}

func getStructInfo(i interface{}, parent *fieldInfo) (*structInfo, error) {
	v := reflect.ValueOf(i)
	for v.Kind() != reflect.Ptr {
		return nil, ErrInvalidConfig
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, ErrInvalidConfig
	}

	si := &structInfo{}
	typ := v.Type()
	if v.Kind() == reflect.Struct {
		n := v.NumField()
		for i := 0; i < n; i++ {
			fv := v.Field(i)
			ft := typ.Field(i)

			// unexported fields
			if !fv.CanSet() {
				continue
			}

			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					if fv.Type().Elem().Kind() != reflect.Struct {
						// nil pointer to a non-struct: leave it alone
						break
					}
					// nil pointer to struct: create a zero instance
					fv.Set(reflect.New(fv.Type().Elem()))
				}
				fv = fv.Elem()
			}

			fi, err := getFieldInfo(fv, ft, parent)
			if err != nil {
				return nil, err
			}

			if fv.Kind() == reflect.Struct {
				p := fi
				// embedded structs
				if ft.Anonymous {
					p = parent
				}
				inner, err := getStructInfo(fv.Addr().Interface(), p)
				if err != nil {
					return nil, err
				}
				si.fields = append(si.fields, inner.fields...)
				continue
			}

			si.fields = append(si.fields, fi)
		}
	}
	return si, nil
}

func getFieldInfo(v reflect.Value, t reflect.StructField, p *fieldInfo) (*fieldInfo, error) {
	fi := &fieldInfo{
		field:  t,
		val:    v,
		parent: p,
	}

	tag, err := parseTag(t)
	if err != nil {
		return nil, err
	}
	fi.tag = *tag
	return fi, nil
}

const (
	tagName              = "config"
	tagSeparator         = ","
	flagFlag             = "flag"
	flagFlagWithValue    = "flag="
	envFlag              = "env"
	envFlagWithValue     = "env="
	defaultFlag          = "default"
	defaultFlagWithValue = "default="
)

type tagInfo struct {
	flag       string
	hasFlag    bool
	env        string
	hasENV     bool
	defVal     string
	hasDefault bool
}

func parseTag(field reflect.StructField) (*tagInfo, error) {
	t := tagInfo{}
	val := field.Tag.Get(tagName)
	tags := strings.Split(val, tagSeparator)
	for _, s := range tags {
		switch {
		case strings.HasPrefix(s, envFlag):
			if err := parseENV(field, &t, s); err != nil {
				return nil, err
			}
		case strings.HasPrefix(s, flagFlag):
			if err := parseFlag(field, &t, s); err != nil {
				return nil, err
			}
		case strings.HasPrefix(s, defaultFlag):
			if err := parseDefault(field, &t, s); err != nil {
				return nil, err
			}
		}
	}

	return &t, nil
}

func parseENV(field reflect.StructField, t *tagInfo, v string) error {
	t.hasENV = true
	if strings.HasPrefix(v, envFlagWithValue) {
		t.env = strings.TrimPrefix(v, envFlagWithValue)
		if t.env == "" {
			return fmt.Errorf("%w, either `env` or `env=ENV_KEY` is valid", ErrInvalidTagFormat)
		}
	}
	return nil
}

func parseFlag(field reflect.StructField, t *tagInfo, v string) error {
	t.hasFlag = true
	if strings.HasPrefix(v, flagFlagWithValue) {
		t.flag = strings.TrimPrefix(v, flagFlagWithValue)
		if t.flag == "" {
			return fmt.Errorf("%w, either `flag` or `flag=flag-key` is valid", ErrInvalidTagFormat)
		}
	}
	return nil
}

func parseDefault(field reflect.StructField, t *tagInfo, v string) error {
	t.hasDefault = true
	if strings.HasPrefix(v, defaultFlagWithValue) {
		t.defVal = strings.TrimPrefix(v, defaultFlagWithValue)
		if t.defVal == "" {
			return fmt.Errorf("%w, either `default` or `default=value` is valid", ErrInvalidTagFormat)
		}
	}
	return nil
}
