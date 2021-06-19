package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous struct containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		}, {
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*24)),
			want: "17 Dec 2020 at 09:00",
		},
	}
	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case.
		// The first parameter to this is the name of the test and the second
		// an anonymous function containing the actual test for each case
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {

				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}
