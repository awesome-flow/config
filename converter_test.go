package config

import (
	"fmt"
	"reflect"
	"testing"
)

type convAct struct {
	res Value
	ok  bool
}

type testConverter struct {
	conv func(kv *KeyValue) (*KeyValue, bool)
}

var _ Converter = (*testConverter)(nil)

func newTestConverter(act convAct) *testConverter {
	return &testConverter{
		conv: func(kv *KeyValue) (*KeyValue, bool) {
			if act.ok {
				return &KeyValue{Key: kv.Key, Value: act.res}, act.ok
			}
			return nil, false
		},
	}
}

func (tc *testConverter) Convert(kv *KeyValue) (*KeyValue, bool) {
	return tc.conv(kv)
}

func TestIdentityConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, 1, true},
		{nil, nil, true},
		{'a', 'a', true},
		{"asdf", "asdf", true},
		{struct{}{}, struct{}{}, true},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := Identity.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestIntToStrConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, "1", true},
		{-1, "-1", true},
		{nil, nil, false},
		{'a', nil, false},
		{"asdf", nil, false},
		{struct{}{}, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := IntToStr.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestIntPtrToIntConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, nil, false},
		{intptr(42), 42, true},
		{nil, nil, false},
		{'a', nil, false},
		{"asdf", nil, false},
		{struct{}{}, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := IntPtrToInt.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestStrPtrToStrConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, nil, false},
		{nil, nil, false},
		{'a', nil, false},
		{"asdf", nil, false},
		{strptr("asdf"), "asdf", true},
		{struct{}{}, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := StrPtrToStr.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestBoolPtrToBoolConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{true, nil, false},
		{boolptr(true), true, true},
		{boolptr(false), false, true},
		{nil, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := BoolPtrToBool.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestStrToIntConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, nil, false},
		{"1", 1, true},
		{"-1", -1, true},
		{"1234567890", 1234567890, true},
		{"asdf", nil, false},
		{'1', nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := StrToInt.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestIfIntConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, 1, true},
		{-1, -1, true},
		{0, 0, true},
		{"asdf", nil, false},
		{intptr(1), nil, false},
		{"1", nil, false},
		{nil, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := IfInt.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestIfStrConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{1, nil, false},
		{"asdf", "asdf", true},
		{strptr("asdf"), nil, false},
		{'a', nil, false},
		{nil, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := IfStr.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestIfBoolConverter(t *testing.T) {
	tests := []struct {
		inVal   interface{}
		outVal  interface{}
		outFlag bool
	}{
		{true, true, true},
		{false, false, true},
		{nil, nil, false},
		{"true", nil, false},
		{1, nil, false},
		{0, nil, false},
	}

	t.Parallel()

	for ix, testCase := range tests {
		t.Run(fmt.Sprintf("Test #%d", ix), func(t *testing.T) {
			in := &KeyValue{Key: nil, Value: testCase.inVal}
			out, ok := IfBool.Convert(in)
			if ok != testCase.outFlag {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.outFlag, ok)
			}
			if !ok {
				return
			}
			if out == nil && testCase.outVal != nil {
				t.Errorf("Expected a non-nil result, got nil")
			}
			if !reflect.DeepEqual(testCase.outVal, out.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.outVal, out.Value)
			}
		})
	}
}

func TestCompositeConverter_CompAnd(t *testing.T) {
	tests := []struct {
		name   string
		chain  []convAct
		expVal Value
		expOk  bool
	}{
		{
			"Empty chain",
			[]convAct{},
			nil, false,
		},
		{
			"1 positive",
			[]convAct{
				convAct{1, true},
			},
			1, true,
		},
		{
			"2 positive",
			[]convAct{
				convAct{1, true},
				convAct{2, true},
			},
			2, true,
		},
		{
			"1 negative",
			[]convAct{
				convAct{nil, false},
			},
			nil, false,
		},
		{
			"1 positive 1 negative",
			[]convAct{
				convAct{1, true},
				convAct{nil, false},
			},
			nil, false,
		},
		{
			"1 negative 1 positive",
			[]convAct{
				convAct{nil, false},
				convAct{1, true},
			},
			nil, false,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			convChain := make([]Converter, 0, len(testCase.chain))
			for _, act := range testCase.chain {
				convChain = append(convChain, newTestConverter(act))
			}

			comp := NewCompositeConverter(CompAnd, convChain...)
			// None of the converters react to the input kv, so
			// passing a nil value
			got, gotOk := comp.Convert(&KeyValue{Key: nil, Value: nil})
			if gotOk != testCase.expOk {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.expOk, gotOk)
			}
			if !gotOk {
				return
			}
			if !reflect.DeepEqual(testCase.expVal, got.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.expVal, got.Value)
			}
		})
	}
}

func TestCompositeConverterCompOr(t *testing.T) {
	tests := []struct {
		name   string
		chain  []convAct
		expVal Value
		expOk  bool
	}{
		{
			"Empty chain",
			[]convAct{},
			nil, false,
		},
		{
			"1 positive",
			[]convAct{
				convAct{1, true},
			},
			1, true,
		},
		{
			"2 positive",
			[]convAct{
				convAct{1, true},
				convAct{2, true},
			},
			1, true,
		},
		{
			"1 negative",
			[]convAct{
				convAct{nil, false},
			},
			nil, false,
		},
		{
			"1 positive 1 negative",
			[]convAct{
				convAct{1, true},
				convAct{nil, false},
			},
			1, true,
		},
		{
			"1 negative 1 positive",
			[]convAct{
				convAct{nil, false},
				convAct{2, true},
			},
			2, true,
		},
		{
			"1 negative 2 positives",
			[]convAct{
				convAct{nil, false},
				convAct{1, true},
				convAct{2, true},
			},
			1, true,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			convChain := make([]Converter, 0, len(testCase.chain))
			for _, act := range testCase.chain {
				convChain = append(convChain, newTestConverter(act))
			}

			comp := NewCompositeConverter(CompOr, convChain...)
			// None of the converters react to the input kv, so
			// passing a nil value
			got, gotOk := comp.Convert(&KeyValue{Key: nil, Value: nil})
			if gotOk != testCase.expOk {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.expOk, gotOk)
			}
			if !gotOk {
				return
			}
			if !reflect.DeepEqual(testCase.expVal, got.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.expVal, got.Value)
			}
		})
	}
}

func TestCompositeConverterCompFirst(t *testing.T) {
	tests := []struct {
		name   string
		chain  []convAct
		expVal Value
		expOk  bool
	}{
		{
			"Empty chain",
			[]convAct{},
			nil, false,
		},
		{
			"1 positive",
			[]convAct{
				convAct{1, true},
			},
			1, true,
		},
		{
			"2 positive",
			[]convAct{
				convAct{1, true},
				convAct{2, true},
			},
			1, true,
		},
		{
			"1 negative",
			[]convAct{
				convAct{nil, false},
			},
			nil, false,
		},
		{
			"1 positive 1 negative",
			[]convAct{
				convAct{1, true},
				convAct{nil, false},
			},
			1, true,
		},
		{
			"1 negative 1 positive",
			[]convAct{
				convAct{nil, false},
				convAct{2, true},
			},
			2, true,
		},
		{
			"1 negative 2 positives",
			[]convAct{
				convAct{nil, false},
				convAct{1, true},
				convAct{2, true},
			},
			1, true,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			convChain := make([]Converter, 0, len(testCase.chain))
			for _, act := range testCase.chain {
				convChain = append(convChain, newTestConverter(act))
			}

			comp := NewCompositeConverter(CompFirst, convChain...)
			// None of the converters react to the input kv, so
			// passing a nil value
			got, gotOk := comp.Convert(&KeyValue{Key: nil, Value: nil})
			if gotOk != testCase.expOk {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.expOk, gotOk)
			}
			if !gotOk {
				return
			}
			if !reflect.DeepEqual(testCase.expVal, got.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.expVal, got.Value)
			}
		})
	}
}

func TestCompositeConverterCompLast(t *testing.T) {
	tests := []struct {
		name   string
		chain  []convAct
		expVal Value
		expOk  bool
	}{
		{
			"Empty chain",
			[]convAct{},
			nil, false,
		},
		{
			"1 positive",
			[]convAct{
				convAct{1, true},
			},
			1, true,
		},
		{
			"2 positive",
			[]convAct{
				convAct{1, true},
				convAct{2, true},
			},
			2, true,
		},
		{
			"1 negative",
			[]convAct{
				convAct{nil, false},
			},
			nil, false,
		},
		{
			"1 positive 1 negative",
			[]convAct{
				convAct{1, true},
				convAct{nil, false},
			},
			1, true,
		},
		{
			"1 negative 1 positive",
			[]convAct{
				convAct{nil, false},
				convAct{2, true},
			},
			2, true,
		},
		{
			"1 negative 2 positives",
			[]convAct{
				convAct{nil, false},
				convAct{1, true},
				convAct{2, true},
			},
			2, true,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			convChain := make([]Converter, 0, len(testCase.chain))
			for _, act := range testCase.chain {
				convChain = append(convChain, newTestConverter(act))
			}

			comp := NewCompositeConverter(CompLast, convChain...)
			// None of the converters react to the input kv, so
			// passing a nil value
			got, gotOk := comp.Convert(&KeyValue{Key: nil, Value: nil})
			if gotOk != testCase.expOk {
				t.Errorf("Unexpected Convert flag: want: %t, got: %t", testCase.expOk, gotOk)
			}
			if !gotOk {
				return
			}
			if !reflect.DeepEqual(testCase.expVal, got.Value) {
				t.Errorf("Unexpected Convert value: want: %#v, got: %#v", testCase.expVal, got.Value)
			}
		})
	}
}
