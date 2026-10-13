package models

import "testing"

func TestParseIntAmount(t *testing.T) {
	tests := []struct {
		name       string
		input      int64
		wantResult string
	}{
		{name: "whole unit with zero subunit", input: 100, wantResult: "1,00"},
		{name: "subunit with leading zero", input: 101, wantResult: "1,01"},
		{name: "multiply whole units with zero subunits", input: 10000, wantResult: "100,00"},
		{name: "two digit subunit", input: 123, wantResult: "1,23"},
		{name: "amount less than ten subunit", input: 1, wantResult: "0,01"},
		{name: "amount less than one unit", input: 99, wantResult: "0,99"},
		{name: "zero amount", input: 0, wantResult: "0,00"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := parseIntAmount(test.input)
			if res != test.wantResult {
				t.Errorf("Test: %s, want result %s, got result %s", test.name, test.wantResult, res)
			}
		})
	}
}
