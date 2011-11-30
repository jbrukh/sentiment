package sentiment

import "strings"
import "regexp"

type Sanitizer struct {
	funcs []func(int, []string) string
}

func NewSanitizer(funcs ...func(int, []string) string) *Sanitizer {
	return &Sanitizer{
		funcs: funcs,
	}
}

func (s *Sanitizer) GetDocument(document string) (result []string) {
    tokens := strings.Split(document, " ")
	result = make([]string, len(tokens)) //functions may add tokens
	for inx, _ := range tokens {
		for _, f := range s.funcs {
			clean := f(inx, tokens)
			if clean != "" {
				result = append(result, clean)
			}
		}
	}
    return
}

func SanitizeToLower(inx int, document []string) (result string) {
	return strings.ToLower(document[inx])
}

func SanitizeNoMentions(inx int, document []string) (result string) {
	result = document[inx]
	if strings.HasPrefix(result, "@") {
		result = ""
	}
	return
}

func SanitizeNoLinks(inx int, document []string) (result string) {
	result = document[inx]
	if strings.HasPrefix(result, "http://") {
		result = ""
	}
	return
}

func SanitizePunctuation(inx int, document []string) (result string) {
	const punct = "?!~`@#$%^&*\\(\\)-_+={}\\[\\]:;|\\\\\"'/.,<>"
	leading := regexp.MustCompile("^[" + punct + "]+|[" + punct[2:] + "]+$")
	return string(leading.ReplaceAll([]byte(document[inx]), []byte("")))
}
