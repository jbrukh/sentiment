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

var username string        // twitter username
var password string        // twitter password
var track *string          // comma-delimited list of tracking keywords for twitter api
var classifier *Classifier // the classifier
var san *Sanitizer         // the sanitizer
var exclList *string       // list of excluded terms
var count [2]int           // the count of all classifications
var highCount [2]int       // the count of all learned classifications
var thresh *float64        // threshold for learning
var printOnly *bool        // suppress classification

func init() {
	track = flag.String("track", "", "comma-separated list of tracking terms")
	thresh = flag.Float64("thresh", DefaultThresh, "the confidence threshold required to learn new content")
	exclList = flag.String("exclude", "", "comma-separated list of keywords excluded from classification")
    printOnly = flat.Bool("print-only", false, "only print the Tweets, do not classify them")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		println("Usage: [--track|--thresh|--exclude|--help] <username> <password>")
		os.Exit(1)
	}
	username = args[0]
	password = args[1]

	// train the classifier
	classifier = NewClassifier(Positive, Negative)
	LearnFile(classifier, "data/positive.txt", Positive)
	LearnFile(classifier, "data/negative.txt", Negative)
	fmt.Println("classifier is trained!")

	// init the sanitizer
	excl := strings.Split(*exclList, ",")
	if *exclList != "" {
		fmt.Printf("excluding: %v\n", excl)
	}

	stopWords := ReadFile("data/stopwords.txt")
	fmt.Printf("stop words: %v\n", stopWords)
	san = NewSanitizer(
		ToLower,
		NoMentions,
		NoLinks,
		NoNumbers,
		Punctuation,
        NoSmallWords,
		CombineNots,
		Exclusions(excl),
		Exclusions(stopWords),
	)
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
		tw := (<-stream).Text
        if (!*printOnly) {
		    process(tw)
        } else {
            fmt.Println(tw)
        }
	}
}

// Obtain a the tweet, classify it, learn it
// if necessary, calculate the positive tweet
// rate and print information.
func process(document string) {
	fmt.Printf("\n> %v\n\n", document)
	// the sanitized document
	doc := san.GetDocument(document)
	if len(doc) < 1 {
		return
	}

	// classification of this document
	scores, inx, _ := classifier.Probabilities(doc)
	logScores, logInx, _ := classifier.Scores(doc)
	class := classifier.Classes[inx]
	logClass := classifier.Classes[logInx]
	count[inx]++

	// the rate of positive sentiment
	posrate := float64(count[0]) / float64(count[0]+count[1])
	learned := ""

	// if above the threshold, then learn
	// this document
	if scores[inx] > *thresh {
		classifier.Learn(doc, class)
		learned = "***"
	}

	// print info
	prettyPrintDoc(doc)
	fmt.Printf("%7.5f %v %v\n", scores[inx], class, learned)
	fmt.Printf("%7.2f %v\n", logScores[logInx], logClass)
	if logClass != class {
		// incorrect classification due to underflow
		fmt.Println("CLASSIFICATION ERROR!")
	}
	fmt.Printf("%7.5f (posrate)\n", posrate)
	//fmt.Printf("%5.5f (high-probability posrate)\n", highrate)
}

func prettyPrintDoc(doc []string) {
	//fmt.Printf("\n%v\n", doc)
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
