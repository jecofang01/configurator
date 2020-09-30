package configurator

import (
	"encoding/base64"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type flagProvider struct {
	flags map[string]func()
}

func NewFlagProvider() *flagProvider {
	return &flagProvider{
		flags: make(map[string]func()),
	}
}

func (p *flagProvider) Provide(v interface{}, si StructInfo) error {
	for _, fi := range si.Fields() {
		k := fi.FlagKey()
		if k == "" {
			continue
		}
		if _, ok := p.flags[k]; ok {
			return fmt.Errorf("flagProvider/Provide: %w [%s]", ErrConflictKey, k)
		}
		fn, err := createVarSetFunc(k, fi.Value(), fi.StructField().Type)
		if err != nil {
			return err
		}
		p.flags[k] = fn
	}
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if fn, ok := p.flags[f.Name]; ok {
			fn()
		}
	})

	return nil
}

var (
	durationType    = reflect.TypeOf(time.Duration(0))
	durationPtrType = reflect.TypeOf((*time.Duration)(nil))
)

func createVarSetFunc(k string, val reflect.Value, typ reflect.Type) (func(), error) {
	switch typ.Kind() {
	case reflect.Bool:
		v := flag.Bool(k, false, "")
		return func() { val.SetBool(*v) }, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		v := flag.Int(k, 0, "")
		return func() { val.SetInt(int64(*v)) }, nil
	case reflect.Int64:
		if typ == durationType {
			v := flag.Duration(k, time.Duration(0), "")
			return func() { val.SetInt(int64(*v)) }, nil
		} else {
			v := flag.Int64(k, 0, "")
			return func() { val.SetInt(*v) }, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		v := flag.Uint(k, 0, "")
		return func() { val.SetUint(uint64(*v)) }, nil
	case reflect.Uint64:
		v := flag.Uint64(k, 0, "")
		return func() { val.SetUint(*v) }, nil
	case reflect.Float32, reflect.Float64:
		v := flag.Float64(k, 0, "")
		return func() { val.SetFloat(*v) }, nil
	case reflect.String:
		v := flag.String(k, "", "")
		return func() { val.SetString(*v) }, nil
	case reflect.Ptr:
		return createPtrSetFunc(k, val, typ)
	case reflect.Slice:
		return createSliceSetFunc(k, val, typ)
	case reflect.Struct:
		if typ == timeType {
			var v timeValue
			flag.Var(&v, k, "")
			return func() {
				t := time.Time(v)
				val.Set(reflect.ValueOf(t))
			}, nil
		}
		return nil, fmt.Errorf("flagProvider/createVarSetFunc: %w type [%s]", ErrUnsupported, typ.Kind().String())
	default:
		return nil, fmt.Errorf("flagProvider/createVarSetFunc: %w type [%s]", ErrUnsupported, typ.Kind().String())
	}
}

func createPtrSetFunc(k string, val reflect.Value, typ reflect.Type) (func(), error) {
	switch typ.Elem().Kind() {
	case reflect.Bool:
		v := flag.Bool(k, false, "")
		return func() {
			val.Set(reflect.ValueOf(v))
		}, nil
	case reflect.Int:
		v := flag.Int(k, 0, "")
		return func() {
			val.Set(reflect.ValueOf(v))
		}, nil
	case reflect.Int8:
		v := flag.Int(k, 0, "")
		return func() {
			i8 := int8(*v)
			val.Set(reflect.ValueOf(&i8))
		}, nil
	case reflect.Int16:
		v := flag.Int(k, 0, "")
		return func() {
			i16 := int16(*v)
			val.Set(reflect.ValueOf(&i16))
		}, nil
	case reflect.Int32:
		v := flag.Int(k, 0, "")
		return func() {
			i32 := int32(*v)
			val.Set(reflect.ValueOf(&i32))
		}, nil
	case reflect.Int64:
		if typ == durationPtrType {
			v := flag.Duration(k, time.Duration(0), "")
			return func() {
				val.Set(reflect.ValueOf(v))
			}, nil
		} else {
			v := flag.Int64(k, 0, "")
			return func() {
				val.Set(reflect.ValueOf(v))
			}, nil
		}
	case reflect.Uint:
		v := flag.Uint(k, 0, "")
		return func() {
			val.Set(reflect.ValueOf(v))
		}, nil
	case reflect.Uint8:
		v := flag.Uint(k, 0, "")
		return func() {
			u8 := uint8(*v)
			val.Set(reflect.ValueOf(&u8))
		}, nil
	case reflect.Uint16:
		v := flag.Uint(k, 0, "")
		return func() {
			u16 := uint16(*v)
			val.Set(reflect.ValueOf(&u16))
		}, nil
	case reflect.Uint32:
		v := flag.Uint(k, 0, "")
		return func() {
			u32 := uint32(*v)
			val.Set(reflect.ValueOf(&u32))
		}, nil
	case reflect.Uint64:
		v := flag.Uint64(k, 0, "")
		return func() {
			val.Set(reflect.ValueOf(v))
		}, nil
	case reflect.Float32:
		v := flag.Float64(k, 0, "")
		return func() {
			f32 := float32(*v)
			val.Set(reflect.ValueOf(&f32))
		}, nil
	case reflect.Float64:
		v := flag.Float64(k, 0, "")
		return func() {
			val.Set(reflect.ValueOf(v))
		}, nil
	case reflect.String:
		v := flag.String(k, "", "")
		return func() {
			val.Set(reflect.ValueOf(v))
		}, nil
	case reflect.Struct:
		if typ == timePtrType {
			var v timeValue
			flag.Var(&v, k, "")
			return func() {
				t := time.Time(v)
				val.Set(reflect.ValueOf(&t))
			}, nil
		}
		return nil, fmt.Errorf("flagProvider/createPtrSetFunc: %w type [%s]", ErrUnsupported, typ.Kind().String())
	default:
		return nil, fmt.Errorf("flagProvider/createPtrSetFunc: %w type [%s]", ErrUnsupported, typ.Kind().String())
	}
}

func createSliceSetFunc(k string, val reflect.Value, typ reflect.Type) (func(), error) {
	switch typ.Elem().Kind() {
	case reflect.Bool:
		var v boolSliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.Int:
		var v intSliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.Int64:
		if typ.Elem() == durationType {
			var v durationSliceValue
			flag.Var(&v, k, "")
			return func() { val.Set(reflect.ValueOf(v)) }, nil
		} else {
			var v int64SliceValue
			flag.Var(&v, k, "")
			return func() { val.Set(reflect.ValueOf(v)) }, nil
		}
	case reflect.Uint:
		var v uintSliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.Uint8:
		var v base64StringValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.Uint64:
		var v uint64SliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.Float32:
		var v float32SliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.Float64:
		var v float64SliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	case reflect.String:
		var v stringSliceValue
		flag.Var(&v, k, "")
		return func() { val.Set(reflect.ValueOf(v)) }, nil
	default:
		return nil, fmt.Errorf("flagProvider/createSliceSetFunc: %w type [%s]", ErrUnsupported, typ.Kind().String())
	}
}

type timeValue time.Time

func (t *timeValue) Get() interface{} { return time.Time(*t) }

func (t *timeValue) String() string { return time.Time(*t).String() }

func (t *timeValue) Set(v string) error {
	tv, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	*t = timeValue(tv)
	return nil
}

type boolSliceValue []bool

func (b *boolSliceValue) Get() interface{} { return []bool(*b) }

func (b *boolSliceValue) String() string { return fmt.Sprintf("%v", []bool(*b)) }

func (b *boolSliceValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*b = append(*b, v)
	return nil
}

type intSliceValue []int

func (i *intSliceValue) Get() interface{} { return []int(*i) }

func (i *intSliceValue) String() string { return fmt.Sprintf("%v", []int(*i)) }

func (i *intSliceValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*i = append(*i, int(v))
	return nil
}

type int64SliceValue []int64

func (i *int64SliceValue) Get() interface{} { return []int64(*i) }

func (i *int64SliceValue) String() string { return fmt.Sprintf("%v", []int64(*i)) }

func (i *int64SliceValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	*i = append(*i, v)
	return nil
}

type durationSliceValue []time.Duration

func (d *durationSliceValue) Get() interface{} { return []time.Duration(*d) }

func (d *durationSliceValue) String() string { return fmt.Sprintf("%v", []time.Duration(*d)) }

func (d *durationSliceValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = append(*d, v)
	return nil
}

type uintSliceValue []uint

func (u *uintSliceValue) Get() interface{} { return []uint(*u) }

func (u *uintSliceValue) String() string { return fmt.Sprintf("%v", []uint(*u)) }

func (u *uintSliceValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		return err
	}
	*u = append(*u, uint(v))
	return nil
}

type uint64SliceValue []uint64

func (u *uint64SliceValue) Get() interface{} { return []uint64(*u) }

func (u *uint64SliceValue) String() string { return fmt.Sprintf("%v", []uint64(*u)) }

func (u *uint64SliceValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*u = append(*u, v)
	return nil
}

type float32SliceValue []float32

func (f *float32SliceValue) Get() interface{} { return []float32(*f) }

func (f *float32SliceValue) String() string { return fmt.Sprintf("%v", []float32(*f)) }

func (f *float32SliceValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = append(*f, float32(v))
	return nil
}

type float64SliceValue []float64

func (f *float64SliceValue) Get() interface{} { return []float64(*f) }

func (f *float64SliceValue) String() string { return fmt.Sprintf("%v", []float64(*f)) }

func (f *float64SliceValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = append(*f, v)
	return nil
}

type stringSliceValue []string

func (s *stringSliceValue) Get() interface{} { return []string(*s) }

func (s *stringSliceValue) String() string { return fmt.Sprintf("%v", []string(*s)) }

func (s *stringSliceValue) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type base64StringValue []byte

func (b *base64StringValue) Get() interface{} { return []byte(*b) }

func (b *base64StringValue) String() string { return base64.StdEncoding.EncodeToString([]byte(*b)) }

func (b *base64StringValue) Set(s string) error {
	bb, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	*b = bb
	return nil
}
