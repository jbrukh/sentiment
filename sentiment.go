package sentiment

import (
    "bytes"
    "fmt"
    "strings"
)

type Histogram struct {
    Freq map[string] int
    Exclusions map[string] bool
}

func NewHistogram(excl ...string) *Histogram {
    exclSet := make(map[string] bool)
    if excl != nil {
        for _, item := range excl {
            exclSet[item] = true    
        }
    }
    return &Histogram {
        make(map[string] int),
        exclSet,
    }
}

func (h *Histogram) Absorb(items []string) {
    for _, item := range items {
        if !h.Exclusions[item] {
            if _, ok := h.Freq[item]; !ok {
                h.Freq[item] = 0
            }
            h.Freq[item] += 1
        }
    }
}

func (h *Histogram) AbsorbText(text, sep string) {
    h.Absorb(strings.Split(text, sep))
}

func (h *Histogram) String() string {
    buffer := bytes.NewBufferString("")
    for key, value := range h.Freq {
        fmt.Fprintf(buffer, "%v:%-10d\n", key, value)
    }
    return string(buffer.Bytes())
}
