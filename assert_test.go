package assert

import (
	"errors"
	"testing"
)

type mockTB struct {
	*testing.T

	ErrorCalls  []errorParams
	ErrorfCalls []errorfParams
	HelperCalls int
}

type errorParams struct {
	args []any
}

type errorfParams struct {
	format string
	args   []any
}

func newMockTB(t *testing.T) *mockTB {
	return &mockTB{
		T:           t,
		ErrorCalls:  []errorParams{},
		ErrorfCalls: []errorfParams{},
		HelperCalls: 0,
	}
}

func (t *mockTB) Error(args ...any) {
	t.ErrorCalls = append(t.ErrorCalls, errorParams{args: args})
}

func (t *mockTB) Errorf(format string, args ...any) {
	t.ErrorfCalls = append(t.ErrorfCalls, errorfParams{format: format, args: args})
}

func (t *mockTB) Helper() {
	t.HelperCalls++
}

func TestError(t *testing.T) {
	mockT := newMockTB(t)
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
			Error(tt.args.t, tt.args.err)

			if n := len(tt.args.t.ErrorCalls); n != tt.expectedCalls {
				t.Errorf("expected %d calls to Error(), got %d", tt.expectedCalls, n)
			}

			if n := len(tt.args.t.ErrorCalls); n != tt.args.t.HelperCalls {
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, n)
			}
		})
	}
}
