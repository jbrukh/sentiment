package sentiment

import (
	. "bayesian"
	"bufio"
	"os"
	"strings"
	"fmt"
)

const (
	Positive Class = "Positive"
	Negative Class = "Negative"
)

// ClassifierPrompt provides a command line prompt to
// query the classifier for scores.
func ClassifierPrompt(classifier *Classifier) {
	for {
		input := bufio.NewReader(os.Stdin)
		line, _, err := input.ReadLine()
		if err != nil || line == "quit" {
			println("exiting")
			break
		}
		scores, _, _ := classifier.Score(strings.Split(string(line), " "))
		fmt.Printf("%v\n", scores)
	}
}

func LearnFile(classifier *Classifier, name string, class Class) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0666)
	if err != nil {
		panic("could not open file")
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		words := strings.Split(string(line), " ")
		fmt.Printf("%v\n", words)
		classifier.Learn(words, class)
	}
}
