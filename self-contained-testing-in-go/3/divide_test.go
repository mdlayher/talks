package main

import "testing"

// START TABLE OMIT
func TestDivide(t *testing.T) {
	tests := []struct {
		name    string // the name of the test being run
		a, b, c int    // inputs a and b, expected output c
		err     error  // expected error
	}{
		{
			name: "divide by zero",
			b:    0,
			err:  errDivideByZero,
		},
		{
			name: "OK",
			a:    4, b: 2, c: 2,
		},
	}
	// ...
	// END TABLE OMIT

	// START LOOP OMIT
	for _, tt := range tests {
		// Begin a subtest using the name of the test from table
		t.Run(tt.name, func(t *testing.T) {
			// Use test inputs to produce output result and error,
			// and check for expected error and expected output
			c, err := Divide(tt.a, tt.b)

			// START CHECK OMIT
			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n - got: %v",
					want, got)
			}
			// END CHECK OMIT

			if want, got := tt.c, c; want != got {
				t.Fatalf("unexpected result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
	// END LOOP OMIT
}
