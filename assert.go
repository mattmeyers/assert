package assert

import (
	"errors"
	"testing"
)

// NoError asserts that the error is nil.
func NoError(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Errorf(`expected no error, got "%+v"`, err)
	}
}

// Error asserts that the error is nil.
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
