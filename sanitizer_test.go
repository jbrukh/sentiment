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

func compareOutput(t *testing.T, input, output []string, f SanitizerFunc) {
    if f == nil {
        return
    }
    result := f(input)
	Assert(t, len(result) == len(output), "length")
	for inx, word := range result {
		Assert(t, word == output[inx], inx, word)
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
	input := []string{"@jake", "jake"}
	output := []string{"jake"}
    compareOutput(t, input, output, SanitizeNoMentions)
}

func TestCombineNotes(t *testing.T) {
	input := []string{"not", "amazing", "not"}
	output := []string{"not", "amazing", "not", "not-amazing"}
    compareOutput(t, input, output, CombineNots)
}

func TestExclusions(t *testing.T) {
    excl := []string{"forbidden", "dontsay", "verboten"}
    input := []string{"dog", "forbidden", "cat", "dontsay", "verboten", "mouse"}
    output := []string{"dog", "cat", "mouse"}
    compareOutput(t, input, output, SanitizeExclusions(excl))
}

func TestExclusionsNil(t *testing.T) {
    excl := []string{}
    input := []string{"dog", "forbidden", "cat", "dontsay", "verboten", "mouse"}
    output := []string{"dog", "forbidden", "cat", "dontsay", "verboten", "mouse"}
    compareOutput(t, input, output, SanitizeExclusions(excl))
}
