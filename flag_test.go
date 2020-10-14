package configurator

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlagProvider(t *testing.T) {
	resetForTesting()
	type example struct {
		B      bool            `config:"flag"`
		Bptr   *bool           `config:"flag"`
		I      int             `config:"flag"`
		Iptr   *int            `config:"flag"`
		I8     int8            `config:"flag"`
		I8ptr  *int8           `config:"flag"`
		I16    int16           `config:"flag"`
		I16ptr *int16          `config:"flag"`
		I32    int32           `config:"flag"`
		I32ptr *int32          `config:"flag"`
		I64    int64           `config:"flag"`
		I64ptr *int64          `config:"flag"`
		D      time.Duration   `config:"flag"`
		Dptr   *time.Duration  `config:"flag"`
		U      uint            `config:"flag"`
		Uptr   *uint           `config:"flag"`
		U8     uint8           `config:"flag"`
		U8ptr  *uint8          `config:"flag"`
		U16    uint16          `config:"flag"`
		U16ptr *uint16         `config:"flag"`
		U32    uint32          `config:"flag"`
		U32ptr *uint32         `config:"flag"`
		U64    uint64          `config:"flag"`
		U64ptr *uint64         `config:"flag"`
		F32    float32         `config:"flag"`
		F32ptr *float32        `config:"flag"`
		F64    float64         `config:"flag"`
		F64ptr *float64        `config:"flag"`
		S      string          `config:"flag"`
		Sptr   *string         `config:"flag"`
		T      time.Time       `config:"flag"`
		Tptr   *time.Time      `config:"flag"`
		Bs     []bool          `config:"flag"`
		Bytes  []byte          `config:"flag"`
		Is     []int           `config:"flag"`
		I64s   []int64         `config:"flag"`
		Ds     []time.Duration `config:"flag"`
		Us     []uint          `config:"flag"`
		U64s   []uint64        `config:"flag"`
		F32s   []float32       `config:"flag"`
		F64s   []float64       `config:"flag"`
		Ss     []string        `config:"flag"`
		Ts     []time.Time     `config:"flag"`
	}

	tt := &example{}
	os.Args = []string{"jhon",
		"-b=True", "-bptr=1", "-i=1", "-iptr=2", "-i8=3", "-i8ptr=4", "-i16=5", "-i16ptr=6", "-i32=7", "-i32ptr=8", "-i64=9", "-i64ptr=10",
		"-d=2s", "-dptr=3ms", "-u=1", "-uptr=2", "-u8=3", "-u8ptr=4", "-u16=5", "-u16ptr=6", "-u32=7", "-u32ptr=8", "-u64=9", "-u64ptr=10",
		"-f32=11", "-f32ptr=12", "-f64=13", "-f64ptr=14", "-s=a", "-sptr=b", "-bs=1", "-bs=0", "-bs=1", "-bytes=AQIDBAoL",
		"-t=2020-09-30T22:51:49-08:00", "-tptr=2020-09-30T22:51:49-08:00", "-is=2", "-is=3", "-i64s=4", "-i64s=5", "-us=5", "-us=6",
		"-u64s=6", "-u64s=7", "-f32s=7", "-f32s=8", "-f64s=8", "-f64s=9", "-ss=abc", "-ss=defg", "-ts=2020-09-30T22:51:49-08:00",
		"-ts=2021-09-30T22:51:49-08:00", "-ds=2s", "-ds=5m",
	}

	si, err := getStructInfo(tt, nil)
	assert.NoError(t, err)

	fp := NewFlagProvider()
	err = fp.Provide(tt, si)
	assert.NoError(t, err)

	t.Logf("%+v", tt)
	assert.Equal(t, &example{
		B:      true,
		Bptr:   b(true),
		I:      1,
		Iptr:   i(2),
		I8:     int8(3),
		I8ptr:  i8(int8(4)),
		I16:    int16(5),
		I16ptr: i16(int16(6)),
		I32:    int32(7),
		I32ptr: i32(int32(8)),
		I64:    int64(9),
		I64ptr: i64(int64(10)),
		D:      2 * time.Second,
		Dptr:   tptr(time.Millisecond * 3),
		U:      uint(1),
		Uptr:   u(uint(2)),
		U8:     uint8(3),
		U8ptr:  u8(uint8(4)),
		U16:    uint16(5),
		U16ptr: u16(uint16(6)),
		U32:    uint32(7),
		U32ptr: u32(uint32(8)),
		U64:    uint64(9),
		U64ptr: u64(uint64(10)),
		F32:    float32(11),
		F32ptr: f32(float32(12)),
		F64:    float64(13),
		F64ptr: f64(float64(14)),
		S:      "a",
		Sptr:   s("b"),
		T:      time.Date(2020, 9, 30, 22, 51, 49, 0, time.FixedZone("", -28800)),
		Tptr:   timePtr(time.Date(2020, 9, 30, 22, 51, 49, 0, time.FixedZone("", -28800))),
		Bs:     []bool{true, false, true},
		Bytes:  []byte{0x01, 0x02, 0x03, 0x04, 0x0a, 0x0b},
		Is:     []int{2, 3},
		I64s:   []int64{int64(4), int64(5)},
		Ds:     []time.Duration{2 * time.Second, 5 * time.Minute},
		Us:     []uint{uint(5), uint(6)},
		U64s:   []uint64{uint64(6), uint64(7)},
		F32s:   []float32{float32(7), float32(8)},
		F64s:   []float64{float64(8), float64(9)},
		Ss:     []string{"abc", "defg"},
		Ts:     []time.Time{time.Date(2020, 9, 30, 22, 51, 49, 0, time.FixedZone("", -28800)), time.Date(2021, 9, 30, 22, 51, 49, 0, time.FixedZone("", -28800))},
	}, tt)
}

func resetForTesting() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func b(v bool) *bool { return &v }

func i(v int) *int { return &v }

func i8(v int8) *int8 { return &v }

func i16(v int16) *int16 { return &v }

func i32(v int32) *int32 { return &v }

func i64(v int64) *int64 { return &v }

func u(v uint) *uint { return &v }

func u8(v uint8) *uint8 { return &v }

func u16(v uint16) *uint16 { return &v }

func u32(v uint32) *uint32 { return &v }

func u64(v uint64) *uint64 { return &v }

func f32(v float32) *float32 { return &v }

func f64(v float64) *float64 { return &v }

func s(v string) *string { return &v }

func tptr(v time.Duration) *time.Duration { return &v }

func timePtr(v time.Time) *time.Time { return &v }
