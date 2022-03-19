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

func newMockTB() *mockTB {
	return &mockTB{
		T:           &testing.T{},
		ErrorCalls:  []errorParams{},
		ErrorfCalls: []errorfParams{},
		HelperCalls: 0,
	}
}

func (t *mockTB) Reset() {
	t.ErrorCalls = []errorParams{}
	t.ErrorfCalls = []errorfParams{}
	t.HelperCalls = 0
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
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, n)
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
				t.Errorf("expected %d calls to Helper(), got %d", tt.expectedCalls, n)
			}
		})
	}
}
