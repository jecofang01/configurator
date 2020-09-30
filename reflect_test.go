package configurator

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	type testStruct struct {
		Name     string `config:"env=am,flag=Me,default=Foo"`
		Empty    string `config:"default=omitempty"`
		Value    string `config:"env,flag"`
		NoTag    string
		EmptyKey string `config:"default=Bar"`
	}
	testObj := testStruct{}
	tests := []struct {
		name  string
		field reflect.StructField
		tag   *tagInfo
	}{
		{
			name:  "flag with value",
			field: reflect.TypeOf(&testObj).Elem().Field(0),
			tag:   &tagInfo{flag: "Me", hasFlag: true, env: "am", hasENV: true, defVal: "Foo", hasDefault: true},
		},
		{
			name:  "omitempty",
			field: reflect.TypeOf(&testObj).Elem().Field(1),
			tag:   &tagInfo{defVal: "omitempty", hasDefault: true},
		},
		{
			name:  "flag without value",
			field: reflect.TypeOf(&testObj).Elem().Field(2),
			tag:   &tagInfo{flag: "", hasFlag: true, env: "", hasENV: true},
		},
		{
			name:  "no tag",
			field: reflect.TypeOf(&testObj).Elem().Field(3),
			tag:   &tagInfo{},
		},
		{
			name:  "empty key",
			field: reflect.TypeOf(&testObj).Elem().Field(4),
			tag:   &tagInfo{hasDefault: true, defVal: "Bar"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ti, err := parseTag(tt.field)
			assert.NoError(t, err)
			assert.Equal(t, tt.tag, ti)
		})
	}
}

func TestGetStructInfo(t *testing.T) {
	type Embedded struct {
		Port int `config:"env=MYSQL_PORT,flag,default=3306"`
	}
	type named struct {
		Host string `config:"env,flag"`
	}
	type namedPtr struct {
		Enable string `config:"env,flag=d"`
	}
	type testStruct struct {
		Embedded
		Name       string `config:"env,flag,default=Foo"`
		MySQL      named
		Debug      *namedPtr
		unexported string
		StartAt    time.Time  `config:"env,flag"`
		Expire     *time.Time `config:"env,flag"`
	}

	cfg := &testStruct{}

	si, err := getStructInfo(cfg, nil)

	assert.NoError(t, err)
	assert.Len(t, si.Fields(), 6)

	assert.Equal(t, "MYSQL_PORT", si.Fields()[0].ENVKey())
	assert.Equal(t, "port", si.Fields()[0].FlagKey())
	assert.Equal(t, "3306", si.Fields()[0].DefVal())

	assert.Equal(t, "NAME", si.Fields()[1].ENVKey())
	assert.Equal(t, "name", si.Fields()[1].FlagKey())
	assert.Equal(t, "Foo", si.Fields()[1].DefVal())

	assert.Equal(t, "MYSQL_HOST", si.Fields()[2].ENVKey())
	assert.Equal(t, "mysql-host", si.Fields()[2].FlagKey())
	assert.Equal(t, "", si.Fields()[2].DefVal())

	assert.Equal(t, "DEBUG_ENABLE", si.fields[3].ENVKey())
	assert.Equal(t, "d", si.fields[3].FlagKey())
	assert.Equal(t, "", si.Fields()[3].DefVal())

	assert.Equal(t, "STARTAT", si.fields[4].ENVKey())
	assert.Equal(t, "startat", si.fields[4].FlagKey())
	assert.Equal(t, "", si.Fields()[4].DefVal())
	assert.True(t, si.Fields()[4].StructField().Type == timeType)

	assert.Equal(t, "EXPIRE", si.fields[5].ENVKey())
	assert.Equal(t, "expire", si.fields[5].FlagKey())
	assert.Equal(t, "", si.Fields()[5].DefVal())
	assert.True(t, si.Fields()[5].StructField().Type == timePtrType)
}
