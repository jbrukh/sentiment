package sentiment

import "strings"
import "regexp"

// constants
const punct = "?!~`@#$%^&*\\(\\)\\-_+={}\\[\\]:;|\\\\\"'/.,<>"

var leading = regexp.MustCompile("^[" + punct + "]+|[" + punct[2:] + "]+$")
var notWords = []string{
	"not", "don't", "dont", "won't", "wont", "weren't", "werent",
	"can't", "cant", "isn't", "isnt", "aren't", "arent",
	"couldn't", "couldnt", "shouldn't", "shouldnt", "wouldn't", "wouldnt"}

// Make all the words lowercase.
func ToLower(words []string) (result []string) {
	return apply(words, func(input string) string {
		return strings.ToLower(input)
	})
}

// Filter out mentions
func NoMentions(words []string) (result []string) {
	return filterIf(words, func(input string) bool {
		return strings.HasPrefix(input, "@")
	})
}

// Filter out links
func NoLinks(words []string) (result []string) {
	return filterIf(words, func(input string) bool {
		return strings.HasPrefix(input, "http://")
	})
}

// Filter out pure numbers
func NoNumbers(words []string) (result []string) {
	numbers := regexp.MustCompile("[0-9]+")
	return filterIf(words, func(input string) bool {
		return numbers.Match([]byte(input))
	})
}

// Filter out single letter and empty words
func NoSmallWords(words []string) (result []string) {
	return filterIf(words, func(input string) bool {
		return len(input) <= 1
	})
}

// Remove excess punctuations and symbols
func Punctuation(words []string) (result []string) {
	result = apply(words, func(input string) string {
		return string(leading.ReplaceAll([]byte(input), []byte("")))
	})
	result = filterIf(words, func(input string) bool { return input == "" })
	return
}

// Results in combining negations with a dash.
func CombineNots(words []string) (result []string) {
	result = make([]string, 0, len(words))
	m := make(map[string]bool)
	for _, item := range notWords {
		m[item] = true
	}
	for inx := 0; inx < len(words); inx++ {
		word := words[inx]
		_, ok := m[word]
		var which string
		if ok && inx != len(words)-1 {
			which = word + "-" + words[inx+1]
			inx++ // skip the next
		} else {
			which = word
		}
		result = append(result, which)
	}
	return
}

func Exclusions(excl []string) SanitizerFunc {
	if len(excl) < 1 {
		return nil
	}
	m := make(map[string]bool, len(excl))
	for _, item := range excl {
		m[item] = true
	}
	return func(words []string) []string {
		return filterIf(words, func(input string) bool {
			_, ok := m[input]
			return ok
		})
	}
}
