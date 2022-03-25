package assert

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

type mockTB struct {
	*testing.T

	ErrorCalls  []messageParams
	ErrorfCalls []formattedMessageParams
	FatalCalls  []messageParams
	FatalfCalls []formattedMessageParams
	HelperCalls int
}

type messageParams struct {
	args []any
}

type formattedMessageParams struct {
	format string
	args   []any
}

func newMockTB() *mockTB {
	return &mockTB{
		T:           &testing.T{},
		ErrorCalls:  []messageParams{},
		ErrorfCalls: []formattedMessageParams{},
		FatalCalls:  []messageParams{},
		FatalfCalls: []formattedMessageParams{},
		HelperCalls: 0,
	}
}

func (t *mockTB) Reset() {
	t.ErrorCalls = []messageParams{}
	t.ErrorfCalls = []formattedMessageParams{}
	t.FatalCalls = []messageParams{}
	t.FatalfCalls = []formattedMessageParams{}
	t.HelperCalls = 0
}

func (t *mockTB) Error(args ...any) {
	t.ErrorCalls = append(t.ErrorCalls, messageParams{args: args})
}

func (t *mockTB) Errorf(format string, args ...any) {
	t.ErrorfCalls = append(t.ErrorfCalls, formattedMessageParams{format: format, args: args})
}

func (t *mockTB) Fatal(args ...any) {
	t.FatalCalls = append(t.FatalCalls, messageParams{args: args})
}

func (t *mockTB) Fatalf(format string, args ...any) {
	t.FatalfCalls = append(t.FatalfCalls, formattedMessageParams{args: args})
}

func (t *mockTB) Helper() {
	t.HelperCalls++
}

func TestFatalTBCallsFatal(t *testing.T) {
	mockT := newMockTB()
	Fatal(mockT).Error("foo")

	if len(mockT.FatalCalls) != 1 {
		t.Errorf("expected 1 call to Fatal(), got %d", len(mockT.FatalCalls))
	}

	if len(mockT.ErrorCalls) != 0 {
		t.Errorf("expected 0 call to Error(), got %d", len(mockT.ErrorCalls))
	}
}

func TestFatalTBCallsFatalf(t *testing.T) {
	mockT := newMockTB()
	Fatal(mockT).Errorf("foo")

	if len(mockT.FatalfCalls) != 1 {
		t.Errorf("expected 1 call to Fatalf(), got %d", len(mockT.FatalfCalls))
	}

	if len(mockT.ErrorfCalls) != 0 {
		t.Errorf("expected 0 call to Errorf(), got %d", len(mockT.ErrorfCalls))
	}
}

func TestError(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t   *mockTB
		err error
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Non nil error",
			args: args{
				t:   mockT,
				err: errors.New("uh oh"),
			},
			expectedCalls: 0,
		},
		{
			name: "Nil error",
			args: args{
				t:   mockT,
				err: nil,
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			Error(tt.args.t, tt.args.err)
			n := len(tt.args.t.ErrorCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Error(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestNoError(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t   *mockTB
		err error
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Non nil error",
			args: args{
				t:   mockT,
				err: errors.New("uh oh"),
			},
			expectedCalls: 1,
		},
		{
			name: "Nil error",
			args: args{
				t:   mockT,
				err: nil,
			},
			expectedCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			NoError(tt.args.t, tt.args.err)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestIsError(t *testing.T) {
	sentinelErr := errors.New("base error")
	wrapperErr := fmt.Errorf("wrapper: %w", sentinelErr)

	mockT := newMockTB()
	type args struct {
		t      *mockTB
		err    error
		target error
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Sentinel error",
			args: args{
				t:      mockT,
				err:    sentinelErr,
				target: sentinelErr,
			},
			expectedCalls: 0,
		},
		{
			name: "Wrapped error",
			args: args{
				t:      mockT,
				err:    wrapperErr,
				target: sentinelErr,
			},
			expectedCalls: 0,
		},
		{
			name: "No match",
			args: args{
				t:      mockT,
				err:    errors.New("no match"),
				target: sentinelErr,
			},
			expectedCalls: 1,
		},
		{
			name: "Nil error, want nil",
			args: args{
				t:      mockT,
				err:    nil,
				target: nil,
			},
			expectedCalls: 0,
		},
		{
			name: "Non nil error, want nil",
			args: args{
				t:      mockT,
				err:    errors.New("uh oh"),
				target: nil,
			},
			expectedCalls: 1,
		},
		{
			name: "Non error, want non nil",
			args: args{
				t:      mockT,
				err:    nil,
				target: errors.New("uh oh"),
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			ErrorIs(tt.args.t, tt.args.err, tt.args.target)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t        *mockTB
		got      int
		expected int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Are equal",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "Are not equal",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			Equal(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t        *mockTB
		got      int
		expected int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Are equal",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "Are not equal",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			NotEqual(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestDeepEqual(t *testing.T) {
	mockT := newMockTB()

	type Foo struct {
		Bar string
		Baz int
	}

	type args struct {
		t        *mockTB
		got      any
		expected any
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Equal structs",
			args: args{
				t:        mockT,
				got:      Foo{Bar: "bar", Baz: 1},
				expected: Foo{Bar: "bar", Baz: 1},
			},
			expectedCalls: 0,
		},
		{
			name: "Equal structs (pointers)",
			args: args{
				t:        mockT,
				got:      &Foo{Bar: "bar", Baz: 1},
				expected: &Foo{Bar: "bar", Baz: 1},
			},
			expectedCalls: 0,
		},
		{
			name: "Not equal structs",
			args: args{
				t:        mockT,
				got:      Foo{Bar: "bar", Baz: 1},
				expected: Foo{Bar: "blah", Baz: 0},
			},
			expectedCalls: 1,
		},
		{
			name: "Equal slices",
			args: args{
				t:        mockT,
				got:      []int{1, 2, 3},
				expected: []int{1, 2, 3},
			},
			expectedCalls: 0,
		},
		{
			name: "Not equal slices",
			args: args{
				t:        mockT,
				got:      []int{1, 2, 3},
				expected: []int{4, 5, 6},
			},
			expectedCalls: 1,
		},
		{
			name: "Equal primitives",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "Not equal primitives",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			DeepEqual(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestNotDeepEqual(t *testing.T) {
	mockT := newMockTB()

	type Foo struct {
		Bar string
		Baz int
	}

	type args struct {
		t        *mockTB
		got      any
		expected any
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Equal structs",
			args: args{
				t:        mockT,
				got:      Foo{Bar: "bar", Baz: 1},
				expected: Foo{Bar: "bar", Baz: 1},
			},
			expectedCalls: 1,
		},
		{
			name: "Equal structs (pointers)",
			args: args{
				t:        mockT,
				got:      &Foo{Bar: "bar", Baz: 1},
				expected: &Foo{Bar: "bar", Baz: 1},
			},
			expectedCalls: 1,
		},
		{
			name: "Not equal structs",
			args: args{
				t:        mockT,
				got:      Foo{Bar: "bar", Baz: 1},
				expected: Foo{Bar: "blah", Baz: 0},
			},
			expectedCalls: 0,
		},
		{
			name: "Equal slices",
			args: args{
				t:        mockT,
				got:      []int{1, 2, 3},
				expected: []int{1, 2, 3},
			},
			expectedCalls: 1,
		},
		{
			name: "Not equal slices",
			args: args{
				t:        mockT,
				got:      []int{1, 2, 3},
				expected: []int{4, 5, 6},
			},
			expectedCalls: 0,
		},
		{
			name: "Equal primitives",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "Not equal primitives",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			NotDeepEqual(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t        *mockTB
		got      int
		expected int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "2 > 1: true",
			args: args{
				t:        mockT,
				got:      2,
				expected: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "1 > 1: false",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "1 > 2: false",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			GreaterThan(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t        *mockTB
		got      int
		expected int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "2 >= 1: true",
			args: args{
				t:        mockT,
				got:      2,
				expected: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "1 >= 1: true",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "1 >= 2: false",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			GreaterThanOrEqual(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestLessThan(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t        *mockTB
		got      int
		expected int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "2 < 1: false",
			args: args{
				t:        mockT,
				got:      2,
				expected: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "1 < 1: false",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "1 < 2: true",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			LessThan(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestLessThanOrEqual(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t        *mockTB
		got      int
		expected int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "2 <= 1: false",
			args: args{
				t:        mockT,
				got:      2,
				expected: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "1 <= 1: true",
			args: args{
				t:        mockT,
				got:      1,
				expected: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "1 <= 2: true",
			args: args{
				t:        mockT,
				got:      1,
				expected: 2,
			},
			expectedCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			LessThanOrEqual(tt.args.t, tt.args.got, tt.args.expected)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	mockT := newMockTB()
	type args[T comparable] struct {
		t      *mockTB
		slice  []T
		values []T
	}
	tests := []struct {
		name          string
		args          args[int]
		expectedCalls int
	}{
		{
			name: "Slice contains single element",
			args: args[int]{
				t:      mockT,
				slice:  []int{1, 2, 3},
				values: []int{2},
			},
			expectedCalls: 0,
		},
		{
			name: "Slice contains multiple elements",
			args: args[int]{
				t:      mockT,
				slice:  []int{1, 2, 3},
				values: []int{3, 2},
			},
			expectedCalls: 0,
		},
		{
			name: "Slice missing single element",
			args: args[int]{
				t:      mockT,
				slice:  []int{1, 2, 3},
				values: []int{4},
			},
			expectedCalls: 1,
		},
		{
			name: "Slice missing multiple elements",
			args: args[int]{
				t:      mockT,
				slice:  []int{1, 2, 3},
				values: []int{4, 5},
			},
			expectedCalls: 2,
		},
		{
			name: "Slice contains and is missing",
			args: args[int]{
				t:      mockT,
				slice:  []int{1, 2, 3},
				values: []int{4, 3, 2, 0},
			},
			expectedCalls: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			SliceContains(tt.args.t, tt.args.slice, tt.args.values...)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestMapContains(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t     *mockTB
		m     map[string]int
		key   string
		value int
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Contains key/value",
			args: args{
				t:     mockT,
				m:     map[string]int{"foo": 1, "bar": 2},
				key:   "foo",
				value: 1,
			},
			expectedCalls: 0,
		},
		{
			name: "Does not contain key",
			args: args{
				t:     mockT,
				m:     map[string]int{"foo": 1, "bar": 2},
				key:   "baz",
				value: 1,
			},
			expectedCalls: 1,
		},
		{
			name: "Contains key but not value",
			args: args{
				t:     mockT,
				m:     map[string]int{"foo": 1, "bar": 2},
				key:   "foo",
				value: 3,
			},
			expectedCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			MapContains(tt.args.t, tt.args.m, tt.args.key, tt.args.value)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if tt.args.t.HelperCalls != 1 {
				t.Errorf("expected 1 call to Helper(), got %d", tt.args.t.HelperCalls)
			}
		})
	}
}

func TestMapContainsKey(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t    *mockTB
		m    map[string]int
		keys []string
	}
	tests := []struct {
		name          string
		args          args
		expectedCalls int
	}{
		{
			name: "Contains multiple keys",
			args: args{
				t:    mockT,
				m:    map[string]int{"foo": 1, "bar": 2},
				keys: []string{"bar", "foo"},
			},
			expectedCalls: 0,
		},
		{
			name: "Contains key but missing others",
			args: args{
				t:    mockT,
				m:    map[string]int{"foo": 1, "bar": 2},
				keys: []string{"bar", "baz", "blah"},
			},
			expectedCalls: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()

			MapContainsKey(tt.args.t, tt.args.m, tt.args.keys...)
			n := len(tt.args.t.ErrorfCalls)

			if n != tt.expectedCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedCalls, n)
			}

			if n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, tt.args.t.HelperCalls)
			}
		})
	}
}

func TestRegexMatches(t *testing.T) {
	mockT := newMockTB()
	type args struct {
		t       *mockTB
		got     string
		pattern string
	}
	tests := []struct {
		name               string
		args               args
		expectedFatalCalls int
		expectedErrorCalls int
	}{
		{
			name: "Basic match",
			args: args{
				t:       mockT,
				got:     "abc123",
				pattern: `\w{3}\d{3}`,
			},
			expectedFatalCalls: 0,
			expectedErrorCalls: 0,
		},
		{
			name: "No match",
			args: args{
				t:       mockT,
				got:     "...",
				pattern: `\w{3}\d{3}`,
			},
			expectedFatalCalls: 0,
			expectedErrorCalls: 1,
		},
		{
			name: "Invalid regex",
			args: args{
				t:       mockT,
				got:     "abc123",
				pattern: `\1`,
			},
			expectedFatalCalls: 1,
			expectedErrorCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Reset()
			regexCache = make(map[string]*regexp.Regexp)

			RegexMatches(tt.args.t, tt.args.got, tt.args.pattern)
			n := len(tt.args.t.ErrorfCalls)
			m := len(tt.args.t.FatalfCalls)

			if n != tt.expectedErrorCalls {
				t.Errorf("expected %d calls to Errorf(), got %d", tt.expectedErrorCalls, n)
			}

			if m != tt.expectedFatalCalls {
				t.Errorf("expected %d calls to Fatalf(), got %d", tt.expectedFatalCalls, m)
			}

			if tt.args.t.HelperCalls != n+m {
				t.Errorf("expected %d calls to Helper(), got %d", n+m, tt.args.t.HelperCalls)
			}
		})
	}
}
