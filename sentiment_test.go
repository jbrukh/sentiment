package sentiment

import "testing"

func TestExclusion(t *testing.T) {
    h := NewHistogram("Jake")
    h.Absorb([]string{"Jake","Prat", "Prat"})
    t.Log(h.String())
        
    _, ok := h.Freq["Jake"] 
    if ok {
        t.Error("Jake should be exluded.")
    }

    freq, ok := h.Freq["Prat"]
    if !ok {
        t.Error("Prat should be included.")
    }

    if freq != 2 {
        t.Error("Prat should be counted twice.")
    }
}
