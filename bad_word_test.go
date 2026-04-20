package main

import "testing"

func TestReplaceBadWord(t *testing.T) {
	badWords := map[string]struct{}{
		"hello": {},
		"bye":  {},
	}
	got := replaceBadWord("Hello everyone", badWords)
	if got != "**** everyone" {
		t.Errorf(`replaceBadWord("hello everyone", %v) = %v; want "**** everyone"`, badWords, got)
	}
}