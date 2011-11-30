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

var username string
var password string
var track *string
var top *int
var classifier *Classifier
var san *Sanitizer
var count [2]int

func init() {
	track = flag.String("track", "", "comma-separated list of tracking terms")
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
		document := sanitize(tw.Text)
		process(document)
	}
}

const thresh = .95

func process(document []string) {
	fmt.Printf("\n%v\n", document)
	scores, inx, _ := classifier.Probabilities(document)
	class := classifier.Classes[inx]
	count[inx]++
	posrate := fmt.Sprintf("%2.5f", float32(count[0])/float32(count[0]+count[1]))
	var learned string
	if scores[inx] > thresh {
		classifier.Learn(document, class)
		learned = "***"
	}
	fmt.Printf("%2.5f %v %v\n", scores[inx], class, learned)
	fmt.Printf("%s (Positive Rate)\n", posrate)
}

func sanitize(text string) (document []string) {
	return san.GetDocument(text)
}
