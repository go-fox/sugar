package shost

import "testing"

func TestExtractHostPort(t *testing.T) {
	extract, err := Extract(":90", nil)
	if err != nil {
		t.Fatal(err)
	}
	println(extract)
}
