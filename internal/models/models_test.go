package models

import "testing"

func TestParseIntAmount(t *testing.T) {
	res := parseIntAmount(123)
	if res != "1,23" {
		t.Errorf("Want result \"1,23\", got result %s", res)
	}
}
