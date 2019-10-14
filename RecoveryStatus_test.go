package main

import "testing"

func TestBuildFromStringOk(t *testing.T) {
	builder := NewPostgresXlogLocationBuilder()

	stringSrc := "1B18/3F9E56B8"

	result, err := builder.BuildFromString(stringSrc)

	if err != nil {
		t.Errorf("Got unexepected error: %s", err)
	} else {
		if result.Upper != 6936 {
			t.Errorf("Expected upper value of 6936, got %d", result.Upper)
		}
		if result.Lower != 1067341496 {
			t.Errorf("Expected lower value of 1067341496, got %d", result.Lower)
		}
	}
}
