package sentiment

import (
	"testing"
	//    "strings"
)

func Assert(t *testing.T, condition bool, args ...interface{}) {
	if !condition {
		t.Fatal(args)
	}
}

func TestPunctuation(t *testing.T) {
	f := SanitizePunctuation
	input := []string{"poop!", "poop!?", "poop!!", "!poop", "!?poop", "hello", "hello&*^*&^$#^$%#",
		"hello[][][[\\||", "$#@hello>>><<,,..", ";\"./hello???"}
	output := []string{"poop!", "poop!?", "poop!!", "poop", "poop",
		"hello", "hello", "hello", "hello", "hello???"}

	result := f(input)
	Assert(t, len(result) == len(output), "length")
	for inx, word := range result {
		Assert(t, word == output[inx], inx, word)
	}
}

func TestNoMentions(t *testing.T) {
	f := SanitizeNoMentions
	input := []string{"@jake", "jake"}
	output := []string{"jake"}
	result := f(input)
	Assert(t, len(result) == len(output), "length")
	for inx, word := range result {
		Assert(t, word == output[inx], inx, word)
	}
}

func TestCombineNotes(t *testing.T) {
	f := CombineNots
	input := []string{"not", "amazing", "not"}
	output := []string{"not", "amazing", "not", "not-amazing"}
	result := f(input)
	Assert(t, len(result) == len(output), "length")
	for inx, word := range result {
		Assert(t, word == output[inx], inx, word)
	}
}
