package sentiment

import "strings"
import "regexp"

// SanitizerFunc will operate on an entire document and return
// the result. Note that the length of the processed array
// need not be the same as the input.
type SanitizerFunc func(words []string) (result []string)

// A Sanitizer will apply a bunch of Sanitizer functions 
// in sequence.
type Sanitizer struct {
	funcs []SanitizerFunc
}

// NewSanitizer returns a new Sanitizer.
func NewSanitizer(funcs ...SanitizerFunc) *Sanitizer {
	return &Sanitizer{
		funcs: funcs,
	}
}

func (s *Sanitizer) GetDocument(document string) (result []string) {
	whitespace := regexp.MustCompile("[\\r\\n\\t ]+")
	document = string(whitespace.ReplaceAll([]byte(document), []byte(" ")))
	result = strings.Split(document, " ")
	for _, f := range s.funcs {
		result = f(result)
	}
	return
}

func apply(words []string, f func(string) string) (result []string) {
	result = words
	for inx, word := range words {
		result[inx] = f(word)
	}
	return
}

func filter(words []string, f func(string) bool) (result []string) {
	result = make([]string, 0, len(words))
	for _, word := range words {
		if f(word) {
			result = append(result, word)
		}
	}
	return
}

func SanitizeToLower(words []string) (result []string) {
	return apply(words, func(input string) string {
		return strings.ToLower(input)
	})
}

func SanitizeNoMentions(words []string) (result []string) {
	return filter(words, func(input string) bool {
		return !strings.HasPrefix(input, "@")
	})
}

func SanitizeNoLinks(words []string) (result []string) {
	return filter(words, func(input string) bool {
		return !strings.HasPrefix(input, "http://")
	})
}

func SanitizeNoNumbers(words []string) (result []string) {
	numbers := regexp.MustCompile("[0-9]+")
	return filter(words, func(input string) bool {
		return !numbers.Match([]byte(input))
	})
}

func SanitizePunctuation(words []string) (result []string) {
	const punct = "?!~`@#$%^&*\\(\\)\\-_+={}\\[\\]:;|\\\\\"'/.,<>"
	leading := regexp.MustCompile("^[" + punct + "]+|[" + punct[2:] + "]+$")
	result = apply(words, func(input string) string {
		return string(leading.ReplaceAll([]byte(input), []byte("")))
	})
	result = filter(words, func(input string) bool { return input != "" })
	return
}

func CombineNots(words []string) (result []string) {
	result = words
	for inx, word := range words {
		if word == "not" && inx != len(words)-1 {
			result = append(result, "not-"+words[inx+1])
		}
	}
	return
}
