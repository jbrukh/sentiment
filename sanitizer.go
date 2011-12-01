package sentiment

import "strings"
import "regexp"

var whitespace = regexp.MustCompile("[\\r\\n\\t ]+")

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
	document = string(whitespace.ReplaceAll([]byte(document), []byte(" ")))
	result = strings.Split(document, " ")
	for _, f := range s.funcs {
		if f != nil {
			result = f(result)
		}
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

func filterIf(words []string, f func(string) bool) (result []string) {
	result = make([]string, 0, len(words))
	for _, word := range words {
		if !f(word) {
			result = append(result, word)
		}
	}
	return
}
