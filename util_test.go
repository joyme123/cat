package main

import "testing"

func TestIsInPac(t *testing.T) {
	list := parsePac("/home/jiang/projects/cat/pac.txt")
	if !isInPac("google.com", list) {
		t.Error("google.com is in pac")
	}

	if !isInPac("www.facebook.com", list) {
		t.Error("www.facebook.com is in pac")
	}
}
