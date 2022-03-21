package assert

import (
	"errors"
	"reflect"
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
