package assert

import (
	"errors"
	"reflect"
	"regexp"
	"testing"
)

// FatalTB is a wrapper around a testing.TB that converts all calls to Error
// and Errorf to Fatal and Fatalf respectively. This can be used to convert
// any assertion to a fatal assertion.
type FatalTB struct {
	testing.TB
}

// Fatal builds a FatalTB for converting an assertion to a fatal assertion.
func Fatal(t testing.TB) *FatalTB {
	return &FatalTB{TB: t}
}

func (t *FatalTB) Error(args ...any) {
	t.TB.Fatal(args...)
}

func (t *FatalTB) Errorf(format string, args ...any) {
	t.TB.Fatalf(format, args...)
}

// NoError asserts that the error is nil.
func NoError(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Errorf(`expected no error, got "%+v"`, err)
	}
}

// Error asserts that the error is not nil.
func Error(t testing.TB, err error) {
	if err == nil {
		t.Helper()
		t.Error("expected error, got nil")
	}
}

// ErrorIs asserts that the error wraps the expected error according to the
// semantics of errors.Is.
func ErrorIs(t testing.TB, err, target error) {
	if !errors.Is(err, target) {
		t.Helper()
		t.Errorf(`expected error "%v", got "%v"`, target, err)
	}
}

// Equal asserts that two comparable values are equivalent.
func Equal[T comparable](t testing.TB, got, expected T) {
	if expected != got {
		t.Helper()
		t.Errorf(`expected "%v", got "%v"`, expected, got)
	}
}

// NotEqual asserts that two comparable values are not equivalent.
func NotEqual[T comparable](t testing.TB, got, expected T) {
	if expected == got {
		t.Helper()
		t.Errorf(`expect "%v" to not equal "%v"`, got, expected)
	}
}

// DeepEqual asserts that two comparable values are equivalent using
// reflect.DeepEqual.
func DeepEqual[T, R any](t testing.TB, got T, expected R) {
	if !reflect.DeepEqual(got, expected) {
		t.Helper()
		t.Errorf(`expected "%+v", got "%+v"`, expected, got)
	}
}

// DeepEqual asserts that two comparable values are not equivalent
// using reflect.DeepEqual.
func NotDeepEqual[T, R any](t testing.TB, got T, expected R) {
	if reflect.DeepEqual(got, expected) {
		t.Helper()
		t.Errorf(`expect "%v" to not equal "%v"`, got, expected)
	}
}

// Ordered represents all types that support the <, <=, >=, and > operators.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// GreaterThan asserts that a value is greater than an expected value.
func GreaterThan[T Ordered](t testing.TB, got, expected T) {
	if got <= expected {
		t.Helper()
		t.Errorf(`expected "%v" to be greater than "%v"`, got, expected)
	}
}

// GreaterThanOrEqual asserts a values is greater than or equal to an expected value.
func GreaterThanOrEqual[T Ordered](t testing.TB, got, expected T) {
	if got < expected {
		t.Helper()
		t.Errorf(`expected "%v" to be greater than or equal to "%v"`, got, expected)
	}
}

// LessThan than asserts a values is less than to an expected value.
func LessThan[T Ordered](t testing.TB, got, expected T) {
	if got >= expected {
		t.Helper()
		t.Errorf(`expected "%v" to be less than "%v"`, got, expected)
	}
}

// LessThanOrEqual asserts a values is less than or equal to an expected value.
func LessThanOrEqual[T Ordered](t testing.TB, got, expected T) {
	if got > expected {
		t.Helper()
		t.Errorf(`expected "%v" to be less than or equal to "%v"`, got, expected)
	}
}

// SliceContains asserts that a slice contains one or more values.
func SliceContains[T comparable](t testing.TB, slice []T, values ...T) {
	for _, v := range values {
		sliceContains(t, slice, v)
	}
}

// sliceContains is a private helper for checking the existence of a single
// element in a slice.
func sliceContains[T comparable](t testing.TB, slice []T, value T) {
	for _, s := range slice {
		if s == value {
			return
		}
	}

	t.Helper()
	t.Errorf(`slice does not contain value "%v"`, value)
}

// MapContains asserts that a map contains the provided key-value pair.
func MapContains[K, V comparable](t testing.TB, m map[K]V, key K, value V) {
	t.Helper()
	v, ok := m[key]
	if !ok {
		t.Errorf(`map does not contain key-value pair "%v: %v"`, key, value)
	} else if v != value {
		t.Errorf(`map contains key %v but not value %v`, key, value)
	}
}

// MapContainsKey asserts that the given map contains the provided keys.
func MapContainsKey[K comparable, V any](t testing.TB, m map[K]V, keys ...K) {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			t.Helper()
			t.Errorf(`map does not contain key "%v"`, k)
		}
	}
}

// regexCache holds all compiled regular expressions. The string pattern input
// is used as the key in the map.
var regexCache = make(map[string]*regexp.Regexp)

// RegexMatches asserts that a provided string is matched by the provided pattern.
// In order to avoid compiling regular expressions many times, they are compiled
// once and cached for future use.
func RegexMatches(t testing.TB, got string, pattern string) {
	r, ok := regexCache[pattern]
	if !ok {
		var err error
		r, err = regexp.Compile(pattern)
		if err != nil {
			t.Helper()
			t.Fatalf("failed to compile regular expression: %v", err)
			return
		}

		regexCache[pattern] = r
	}

	if !r.MatchString(got) {
		t.Helper()
		t.Errorf("received string %s not matched by pattern /%s/", got, pattern)
	}
}
