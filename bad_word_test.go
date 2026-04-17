package main

import "testing"

func TestReplaceBadWord(t *testing.T) {
	got := replaceBadWord("Hello everyone", []string{"hello", "bye"})
	if got != "**** everyone" {
		t.Errorf(`replaceBadWord("hello everyone", []string{"hello", "bye"}) = %v; want "**** everyone"`, got)
	}
}