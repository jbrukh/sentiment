package sentiment

import (
	"bytes"
	"fmt"
	"strings"
    "sort"
)

type Histogram struct {
	Freq       map[string]int
	Exclusions map[string]bool
}

// NewHistogram returns a new, empty histogram
func NewHistogram() *Histogram {
	return &Histogram{
		make(map[string]int),
		make(map[string]bool),
	}
}

// Exclude provides a list of strings that will
// be excluded from being processed by the histogram.
func (h *Histogram) Exclude(excl []string) {
	if excl != nil {
		for _, item := range excl {
			if item != "" {
				h.Exclusions[item] = true
			}
		}
	}
}

// Absorb will add the specified list of tokens to
// the histogram. A token will be added unless it is
// exluded or the empty string.
func (h *Histogram) Absorb(items []string) {
	for _, item := range items {
		if !h.Exclusions[item] && item != "" {
			if _, ok := h.Freq[item]; !ok {
				h.Freq[item] = 0
			}
			h.Freq[item] += 1
		}
	}
}

// AbsorbText will absorb the specified block of text,
// doing a rudimentary tokenization of it using the
// strings.Split() function.
func (h *Histogram) AbsorbText(text, sep string) {
	h.Absorb(strings.Split(text, sep))
}


func (h *Histogram) String() string {
	buffer := bytes.NewBufferString("")
	for key, value := range h.Freq {
		fmt.Fprintf(buffer, "%v: %d, ", key, value)
	}
	return string(buffer.Bytes())
}

// Token popularity
type TokenPop struct {
    Token string
    Pop int
}
type TokenPops []TokenPop

func (p TokenPops) Swap(i, j int) { p[i], p[j] = p[j], p[i] } 
func (p TokenPops) Len() int { return len(p) } 
func (p TokenPops) Less(i, j int) bool { return p[i].Pop > p[j].Pop } 

func (h *Histogram) MostPopular() TokenPops {
    pops := make(TokenPops, len(h.Freq))
    i := 0
    for k, v := range h.Freq {
        pops[i] = TokenPop{k, v}; i++
    }
    sort.Sort(pops)
    return pops
}
