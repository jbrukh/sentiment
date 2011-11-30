package sentiment

import "strings"

type Sanitizer struct {
    funcs []func(int,[]string)string
}

func NewSanitizer(funcs ...func(int,[]string)string) *Sanitizer {
    return &Sanitizer{
        funcs: funcs,
    }
}

func (s *Sanitizer) GetDocument(document []string) {
    result := make([]string, len(document)) //functions may add tokens
    for inx, _ := range document {
        for _, f := range s.funcs {
           clean := f(inx, document)
           if clean != "" {
              result = append(result, clean)
           }
        }
    }
}

func SanitizeToLower(inx int, document []string) (result string) {
    return strings.ToLower(document[inx])
}

func SanitizeNoMentions(inx int, document []string) (result string) {
    result = document[inx]
    if strings.HasPrefix(token, "@") {
        result = ""
    }
    return
}

func SanitizeNoLinks(inx int, document []string) (result string) {
    result = document[inx]
    if strings.HasPrefix(token, "http://") {
        result = ""
    }
    return
}

func SanitizePunctuation(inx int, document []string) (result string) {
    const punct = "?!~`@#$%^&*\\(\\)-_+={}\\[\\]:;|\\\\\"'/.,<>"
    leading := regexp.MustCompile("^["+punct+"]+|["+punct[2:]+"]+$")
    result = leading.ReplaceAll([]byte(document[inx]), []byte(""))
    return string(result)
}
