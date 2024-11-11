package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	type customString string
	type customInt int
	type customInt8 int8
	type customInt16 int16
	type customInt32 int32
	type customInt64 int64
	type customUInt uint
	type customUInt8 uint8
	type customUInt16 uint16
	type customUInt32 uint32
	type customUInt64 uint64
	type customFloat32 float32
	type customFloat64 float64
	type customUintptr uintptr
	type customTime time.Time
	type customBytes []byte
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		cType   string
	}{
		{less: customString("a"), greater: customString("b"), cType: "string"},
		{less: "a", greater: "b", cType: "string"},
		{less: customInt(1), greater: customInt(2), cType: "int"},
		{less: int(1), greater: int(2), cType: "int"},
		{less: customInt8(1), greater: customInt8(2), cType: "int8"},
		{less: int8(1), greater: int8(2), cType: "int8"},
		{less: customInt16(1), greater: customInt16(2), cType: "int16"},
		{less: int16(1), greater: int16(2), cType: "int16"},
		{less: customInt32(1), greater: customInt32(2), cType: "int32"},
		{less: int32(1), greater: int32(2), cType: "int32"},
		{less: customInt64(1), greater: customInt64(2), cType: "int64"},
		{less: int64(1), greater: int64(2), cType: "int64"},
		{less: customUInt(1), greater: customUInt(2), cType: "uint"},
		{less: uint8(1), greater: uint8(2), cType: "uint8"},
		{less: customUInt8(1), greater: customUInt8(2), cType: "uint8"},
		{less: uint16(1), greater: uint16(2), cType: "uint16"},
		{less: customUInt16(1), greater: customUInt16(2), cType: "uint16"},
		{less: uint32(1), greater: uint32(2), cType: "uint32"},
		{less: customUInt32(1), greater: customUInt32(2), cType: "uint32"},
		{less: uint64(1), greater: uint64(2), cType: "uint64"},
		{less: customUInt64(1), greater: customUInt64(2), cType: "uint64"},
		{less: float32(1.23), greater: float32(2.34), cType: "float32"},
		{less: customFloat32(1.23), greater: customFloat32(2.23), cType: "float32"},
		{less: float64(1.23), greater: float64(2.34), cType: "float64"},
		{less: customFloat64(1.23), greater: customFloat64(2.34), cType: "float64"},
		{less: uintptr(1), greater: uintptr(2), cType: "uintptr"},
		{less: customUintptr(1), greater: customUintptr(2), cType: "uint64"},
		{less: time.Now(), greater: time.Now().Add(time.Hour), cType: "time.Time"},
		{less: time.Date(2024, 0, 0, 0, 0, 0, 0, time.Local), greater: time.Date(2263, 0, 0, 0, 0, 0, 0, time.Local), cType: "time.Time"},
		{less: customTime(time.Now()), greater: customTime(time.Now().Add(time.Hour)), cType: "time.Time"},
		{less: []byte{1, 1}, greater: []byte{1, 2}, cType: "[]byte"},
		{less: customBytes([]byte{1, 1}), greater: customBytes([]byte{1, 2}), cType: "[]byte"},
	} {
		resLess, isComparable := compare(currCase.less, currCase.greater, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object should be comparable for type " + currCase.cType)
		}

		if resLess != compareLess {
			t.Errorf("object less (%v) should be less than greater (%v) for type "+currCase.cType,
				currCase.less, currCase.greater)
		}

		resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resGreater != compareGreater {
			t.Errorf("object greater should be greater than less for type " + currCase.cType)
		}

		resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resEqual != 0 {
			t.Errorf("objects should be equal for type " + currCase.cType)
		}
	}
}

type outputT struct {
	buf     *bytes.Buffer
	helpers map[string]struct{}
}

// Implements TestingT
func (t *outputT) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	t.buf.WriteString(s)
}

func (t *outputT) Helper() {
	if t.helpers == nil {
		t.helpers = make(map[string]struct{})
	}
	t.helpers[callerName(1)] = struct{}{}
}

// callerName gives the function name (qualified with a package path)
// for the caller after skip frames (where 0 means the current function).
func callerName(skip int) string {
	// Make room for the skip PC.
	var pc [1]uintptr
	n := runtime.Callers(skip+2, pc[:]) // skip + runtime.Callers + callerName
	if n == 0 {
		panic("testing: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return strings.TrimSuffix(frame.Function, "[...]")
}

func TestGreater(t *testing.T) {
	mockT := new(testing.T)

	if !Greater(mockT, 2, 1) {
		t.Error("Greater should return true")
	}

	if Greater(mockT, 1, 1) {
		t.Error("Greater should return false")
	}

	if Greater(mockT, 1, 2) {
		t.Error("Greater should return false")
	}

	// Check error report

	// Old tests
	checkGreater(t, "a", "b", `"a" is not greater than "b"`)
	checkGreater(t, int(1), int(2), `"1" is not greater than "2"`)
	checkGreater(t, int8(1), int8(2), `"1" is not greater than "2"`)
	checkGreater(t, int16(1), int16(2), `"1" is not greater than "2"`)
	checkGreater(t, int32(1), int32(2), `"1" is not greater than "2"`)
	checkGreater(t, int64(1), int64(2), `"1" is not greater than "2"`)
	checkGreater(t, uint8(1), uint8(2), `"1" is not greater than "2"`)
	checkGreater(t, uint16(1), uint16(2), `"1" is not greater than "2"`)
	checkGreater(t, uint32(1), uint32(2), `"1" is not greater than "2"`)
	checkGreater(t, uint64(1), uint64(2), `"1" is not greater than "2"`)
	checkGreater(t, float32(1), float32(2), `"1" is not greater than "2"`)
	checkGreater(t, float64(1), float64(2), `"1" is not greater than "2"`)
	checkGreater(t, uintptr(1), uintptr(2), `"1" is not greater than "2"`)
	checkGreater(t, time.Time{}, time.Time{}.Add(time.Hour), `"0001-01-01 00:00:00 +0000 UTC" is not greater than "0001-01-01 01:00:00 +0000 UTC"`)
	checkGreater(t, []byte{1, 1}, []byte{1, 2}, `"[1 1]" is not greater than "[1 2]"`)

	// New tests
	checkGreater(t, 1, 2, `"1" is not greater than "2"`)
}

func checkGreater[T1, T2 Ordered | []byte | time.Time](t TestingT, less T1, greater T2, msg string) {
	t.Helper()
	out := &outputT{buf: bytes.NewBuffer(nil)}
	False(t, Greater(out, less, greater))
	Contains(t, out.buf.String(), msg)
	Contains(t, out.helpers, "github.com/stretchr/testify/assert.Greater")
}

func TestGreaterOrEqual(t *testing.T) {
	mockT := new(testing.T)

	if !GreaterOrEqual(mockT, 2, 1) {
		t.Error("GreaterOrEqual should return true")
	}

	if !GreaterOrEqual(mockT, 1, 1) {
		t.Error("GreaterOrEqual should return true")
	}

	if GreaterOrEqual(mockT, 1, 2) {
		t.Error("GreaterOrEqual should return false")
	}

	// Check error report

	// Old tests
	checkGreaterOrEqual(t, "a", "b", `"a" is not greater than or equal to "b"`)
	checkGreaterOrEqual(t, int(1), int(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, int8(1), int8(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, int16(1), int16(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, int32(1), int32(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, int64(1), int64(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, uint8(1), uint8(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, uint16(1), uint16(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, uint32(1), uint32(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, uint64(1), uint64(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, float32(1), float32(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, float64(1), float64(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, uintptr(1), uintptr(2), `"1" is not greater than or equal to "2"`)
	checkGreaterOrEqual(t, time.Time{}, time.Time{}.Add(time.Hour), `"0001-01-01 00:00:00 +0000 UTC" is not greater than or equal to "0001-01-01 01:00:00 +0000 UTC"`)
	checkGreaterOrEqual(t, []byte{1, 1}, []byte{1, 2}, `"[1 1]" is not greater than or equal to "[1 2]"`)

	// New tests
	checkGreaterOrEqual(t, 1, 2, `"1" is not greater than or equal to "2"`)
}

func checkGreaterOrEqual[T1, T2 Ordered | []byte | time.Time](t TestingT, less T1, greater T2, msg string) {
	t.Helper()
	out := &outputT{buf: bytes.NewBuffer(nil)}
	False(t, GreaterOrEqual(out, less, greater))
	Contains(t, out.buf.String(), msg)
	Contains(t, out.helpers, "github.com/stretchr/testify/assert.GreaterOrEqual")
}

func TestLess(t *testing.T) {
	mockT := new(testing.T)

	if !Less(mockT, 1, 2) {
		t.Error("Less should return true")
	}

	if Less(mockT, 1, 1) {
		t.Error("Less should return false")
	}

	if Less(mockT, 2, 1) {
		t.Error("Less should return false")
	}

	// Old tests
	checkLess(t, "b", "a", `"b" is not less than "a"`)
	checkLess(t, int(2), int(1), `"2" is not less than "1"`)
	checkLess(t, int8(2), int8(1), `"2" is not less than "1"`)
	checkLess(t, int16(2), int16(1), `"2" is not less than "1"`)
	checkLess(t, int32(2), int32(1), `"2" is not less than "1"`)
	checkLess(t, int64(2), int64(1), `"2" is not less than "1"`)
	checkLess(t, uint8(2), uint8(1), `"2" is not less than "1"`)
	checkLess(t, uint16(2), uint16(1), `"2" is not less than "1"`)
	checkLess(t, uint32(2), uint32(1), `"2" is not less than "1"`)
	checkLess(t, uint64(2), uint64(1), `"2" is not less than "1"`)
	checkLess(t, float32(2), float32(1), `"2" is not less than "1"`)
	checkLess(t, float64(2), float64(1), `"2" is not less than "1"`)
	checkLess(t, uintptr(2), uintptr(1), `"2" is not less than "1"`)
	checkLess(t, time.Time{}.Add(time.Hour), time.Time{}, `"0001-01-01 01:00:00 +0000 UTC" is not less than "0001-01-01 00:00:00 +0000 UTC"`)
	checkLess(t, []byte{1, 2}, []byte{1, 1}, `"[1 2]" is not less than "[1 1]"`)

	// New tests
	checkLess(t, 2, 1, `"2" is not less than "1"`)
}

func checkLess[T1, T2 Ordered | []byte | time.Time](t TestingT, less T1, greater T2, msg string) {
	t.Helper()
	out := &outputT{buf: bytes.NewBuffer(nil)}
	False(t, Less(out, less, greater))
	Contains(t, out.buf.String(), msg)
	Contains(t, out.helpers, "github.com/stretchr/testify/assert.Less")
}

func TestLessOrEqual(t *testing.T) {
	mockT := new(testing.T)

	if !LessOrEqual(mockT, 1, 2) {
		t.Error("LessOrEqual should return true")
	}

	if !LessOrEqual(mockT, 1, 1) {
		t.Error("LessOrEqual should return true")
	}

	if LessOrEqual(mockT, 2, 1) {
		t.Error("LessOrEqual should return false")
	}

	// Old tests
	checkLessOrEqual(t, "b", "a", `"b" is not less than or equal to "a"`)
	checkLessOrEqual(t, int(2), int(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, int8(2), int8(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, int16(2), int16(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, int32(2), int32(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, int64(2), int64(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, uint8(2), uint8(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, uint16(2), uint16(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, uint32(2), uint32(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, uint64(2), uint64(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, float32(2), float32(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, float64(2), float64(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, uintptr(2), uintptr(1), `"2" is not less than or equal to "1"`)
	checkLessOrEqual(t, time.Time{}.Add(time.Hour), time.Time{}, `"0001-01-01 01:00:00 +0000 UTC" is not less than or equal to "0001-01-01 00:00:00 +0000 UTC"`)
	checkLessOrEqual(t, []byte{1, 2}, []byte{1, 1}, `"[1 2]" is not less than or equal to "[1 1]"`)

	// New tests
	checkLessOrEqual(t, 2, 1, `"2" is not less than or equal to "1"`)
}

func checkLessOrEqual[T1, T2 Ordered | []byte | time.Time](t TestingT, less T1, greater T2, msg string) {
	t.Helper()
	out := &outputT{buf: bytes.NewBuffer(nil)}
	False(t, LessOrEqual(out, less, greater))
	Contains(t, out.buf.String(), msg)
	Contains(t, out.helpers, "github.com/stretchr/testify/assert.LessOrEqual")
}

func TestPositive(t *testing.T) {
	mockT := new(testing.T)

	if !Positive(mockT, 1) {
		t.Error("Positive should return true")
	}

	if !Positive(mockT, 1.23) {
		t.Error("Positive should return true")
	}

	if Positive(mockT, -1) {
		t.Error("Positive should return false")
	}

	if Positive(mockT, -1.23) {
		t.Error("Positive should return false")
	}

	// Old tests
	checkPositive(t, int(-1), `"-1" is not positive`)
	checkPositive(t, int8(-1), `"-1" is not positive`)
	checkPositive(t, int16(-1), `"-1" is not positive`)
	checkPositive(t, int32(-1), `"-1" is not positive`)
	checkPositive(t, int64(-1), `"-1" is not positive`)
	checkPositive(t, float32(-1), `"-1" is not positive`)
	checkPositive(t, float64(-1), `"-1" is not positive`)

	// New tests
	checkPositive(t, -1, `"-1" is not positive`)
}

func checkPositive[T Number](t TestingT, e T, msg string) {
	t.Helper()
	out := &outputT{buf: bytes.NewBuffer(nil)}
	False(t, Positive(out, e, msg))
	Contains(t, out.buf.String(), msg)
	Contains(t, out.helpers, "github.com/stretchr/testify/assert.Positive")
}

func TestNegative(t *testing.T) {
	mockT := new(testing.T)

	if !Negative(mockT, -1) {
		t.Error("Negative should return true")
	}

	if !Negative(mockT, -1.23) {
		t.Error("Negative should return true")
	}

	if Negative(mockT, 1) {
		t.Error("Negative should return false")
	}

	if Negative(mockT, 1.23) {
		t.Error("Negative should return false")
	}

	// Check error report

	// Old tests
	checkNegative(t, int(1), `"1" is not negative`)
	checkNegative(t, int8(1), `"1" is not negative`)
	checkNegative(t, int16(1), `"1" is not negative`)
	checkNegative(t, int32(1), `"1" is not negative`)
	checkNegative(t, int64(1), `"1" is not negative`)
	checkNegative(t, float32(1), `"1" is not negative`)
	checkNegative(t, float64(1), `"1" is not negative`)

	// New tests
	checkNegative(t, 1, `"1" is not negative`)
}

func checkNegative[T Number](t TestingT, e T, msg string) {
	t.Helper()
	out := &outputT{buf: bytes.NewBuffer(nil)}
	False(t, Negative(out, e, msg))
	Contains(t, out.buf.String(), msg)
	Contains(t, out.helpers, "github.com/stretchr/testify/assert.Negative")
}

func Test_compareTwoValuesDifferentValuesTypes(t *testing.T) {
	mockT := new(testing.T)

	for _, currCase := range []struct {
		v1            interface{}
		v2            interface{}
		compareResult bool
	}{
		{v1: 123, v2: "abc"},
		{v1: "abc", v2: 123456},
		{v1: float64(12), v2: "123"},
		{v1: "float(12)", v2: float64(1)},
	} {
		result := compareTwoValues(mockT, currCase.v1, currCase.v2, []compareResult{compareLess, compareEqual, compareGreater}, "testFailMessage")
		False(t, result)
	}
}

func Test_compareTwoValuesNotComparableValues(t *testing.T) {
	mockT := new(testing.T)

	type CompareStruct struct {
	}

	for _, currCase := range []struct {
		v1 interface{}
		v2 interface{}
	}{
		{v1: CompareStruct{}, v2: CompareStruct{}},
		{v1: map[string]int{}, v2: map[string]int{}},
		{v1: make([]int, 5), v2: make([]int, 5)},
	} {
		result := compareTwoValues(mockT, currCase.v1, currCase.v2, []compareResult{compareLess, compareEqual, compareGreater}, "testFailMessage")
		False(t, result)
	}
}

func Test_compareTwoValuesCorrectCompareResult(t *testing.T) {
	mockT := new(testing.T)

	for _, currCase := range []struct {
		v1             interface{}
		v2             interface{}
		allowedResults []compareResult
	}{
		{v1: 1, v2: 2, allowedResults: []compareResult{compareLess}},
		{v1: 1, v2: 2, allowedResults: []compareResult{compareLess, compareEqual}},
		{v1: 2, v2: 2, allowedResults: []compareResult{compareGreater, compareEqual}},
		{v1: 2, v2: 2, allowedResults: []compareResult{compareEqual}},
		{v1: 2, v2: 1, allowedResults: []compareResult{compareEqual, compareGreater}},
		{v1: 2, v2: 1, allowedResults: []compareResult{compareGreater}},
	} {
		result := compareTwoValues(mockT, currCase.v1, currCase.v2, currCase.allowedResults, "testFailMessage")
		True(t, result)
	}
}

func Test_containsValue(t *testing.T) {
	for _, currCase := range []struct {
		values []compareResult
		value  compareResult
		result bool
	}{
		{values: []compareResult{compareGreater}, value: compareGreater, result: true},
		{values: []compareResult{compareGreater, compareLess}, value: compareGreater, result: true},
		{values: []compareResult{compareGreater, compareLess}, value: compareLess, result: true},
		{values: []compareResult{compareGreater, compareLess}, value: compareEqual, result: false},
	} {
		result := containsValue(currCase.values, currCase.value)
		Equal(t, currCase.result, result)
	}
}

func TestComparingMsgAndArgsForwarding(t *testing.T) {
	msgAndArgs := []interface{}{"format %s %x", "this", 0xc001}
	expectedOutput := "format this c001\n"
	funcs := []func(t TestingT){
		func(t TestingT) { Greater(t, 1, 2, msgAndArgs...) },
		func(t TestingT) { GreaterOrEqual(t, 1, 2, msgAndArgs...) },
		func(t TestingT) { Less(t, 2, 1, msgAndArgs...) },
		func(t TestingT) { LessOrEqual(t, 2, 1, msgAndArgs...) },
		func(t TestingT) { Positive(t, 0, msgAndArgs...) },
		func(t TestingT) { Negative(t, 0, msgAndArgs...) },
	}
	for _, f := range funcs {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		f(out)
		Contains(t, out.buf.String(), expectedOutput)
	}
}
