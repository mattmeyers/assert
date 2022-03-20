package assert

import (
	"errors"
	"reflect"
	"testing"
)

// NoError asserts that the error is nil. This assertion will trigger a t.Fatal
// call if the error is not nil and is intended to be used when checking errors
// not directly related to the test. For a non-fatal check, use ErrorIs.
func NoError(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf(`expected no error, got "%+v"`, err)
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
