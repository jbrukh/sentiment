//target:analyzer
package main

import (
	"fmt"
	"os"
	"twitterstream"
	"flag"
	"strings"
	. "sentiment"
	. "bayesian"
)

const DefaultThresh = .95

var username string
var password string
var track *string
var top *int
var classifier *Classifier
var san *Sanitizer
var count [2]int
var highCount [2]int
var thresh *float64

func init() {
	track = flag.String("track", "", "comma-separated list of tracking terms")
	thresh = flag.Float64("thresh", DefaultThresh, "the confidence threshold required to learn new content")

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		println("Usage: [flags...] <username> <password>")
		os.Exit(1)
	}
	username = args[0]
	password = args[1]

	// train the classifier
	classifier = NewClassifier(Positive, Negative)
	LearnFile(classifier, "data/positive.txt", Positive)
	LearnFile(classifier, "data/negative.txt", Negative)
	fmt.Println("classifier is trained...")

	// init the sanitizer
	san = NewSanitizer(SanitizeToLower,
		SanitizeNoMentions,
		SanitizeNoLinks,
		SanitizeNoNumbers,
		SanitizePunctuation,
		CombineNots)
}

func main() {
	stream := make(chan *twitterstream.Tweet)
	client := twitterstream.NewClient(username, password)

	fmt.Printf("track = %v\n", *track)
	tracks := strings.Split(*track, ",")

	err := client.Track(tracks, stream)
	if err != nil {
		println(err.String())
	}

	for {
		tw := <-stream
		process(tw.Text)
	}
}

// Obtain a the tweet, classify it, learn it
// if necessary, calculate the positive tweet
// rate and print information.
func process(document string) {
	// the sanitized document
	doc := san.GetDocument(document)

	// classification of this document
	scores, inx, _ := classifier.Probabilities(doc)
	class := classifier.Classes[inx]
	count[inx]++
	highCount[inx]++

	// the rate of positive sentiment
	posrate := float64(count[0]) / float64(count[0]+count[1])
	highrate := float64(highCount[0]) / float64(highCount[0]+highCount[1])
	learned := ""

	// if above the threshold, then learn
	// this document
	if scores[inx] > *thresh {
		classifier.Learn(doc, class)
		learned = "***"
		highCount[inx]++
	}

	// print info
	prettyPrintDoc(doc)
	fmt.Printf("%2.5f %v %v\n", scores[inx], class, learned)
	fmt.Printf("%2.5f (all posrate)\n", posrate)
	fmt.Printf("%2.5f (high-probability posrate)\n", highrate)
}

func prettyPrintDoc(doc []string) {
	fmt.Printf("\n%v\n", doc)
	fmt.Printf("\t")
	for _, word := range doc {
		fmt.Printf("%7s", abbrev(word, 5))
	}
	fmt.Println("")

	freqs := classifier.WordFrequencies(doc)
	for i := 0; i < 2; i++ {
		fmt.Printf("%6s", classifier.Classes[i])
		for j := 0; j < len(doc); j++ {
			fmt.Printf("%7.4f", freqs[i][j])
		}
		fmt.Println("")
	}
}

func abbrev(word string, max int) (result string) {
	result = word
	if max < 5 {
		panic("max must be at least 5")
	}
	if len(word) > max {
		result = word[0:max-2] + ".."
	}
	return
}
