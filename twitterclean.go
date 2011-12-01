package sentiment

import "strings"
import "regexp"

// constants
const punct = "?!~`@#$%^&*\\(\\)\\-_+={}\\[\\]:;|\\\\\"'/.,<>"

var leading = regexp.MustCompile("^[" + punct + "]+|[" + punct[2:] + "]+$")
var StopWords = ReadFile("data/stopwords.txt")

func ToLower(words []string) (result []string) {
	return apply(words, func(input string) string {
		return strings.ToLower(input)
	})
}

func NoMentions(words []string) (result []string) {
	return filterIf(words, func(input string) bool {
		return strings.HasPrefix(input, "@")
	})
}

func NoLinks(words []string) (result []string) {
	return filterIf(words, func(input string) bool {
		return strings.HasPrefix(input, "http://")
	})
}

func NoNumbers(words []string) (result []string) {
	numbers := regexp.MustCompile("[0-9]+")
	return filterIf(words, func(input string) bool {
		return numbers.Match([]byte(input))
	})
}

func SmallWords(words []string) (result []string) {
	return filterIf(words, func(input string) bool {
		return len(input) < 3
	})
}

func Punctuation(words []string) (result []string) {
	result = apply(words, func(input string) string {
		return string(leading.ReplaceAll([]byte(input), []byte("")))
	})
	result = filterIf(words, func(input string) bool { return input == "" })
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
