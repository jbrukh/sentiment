//target:analyzer
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
	"twitterstream"
	"flag"
	"strings"
	. "sentiment"
	. "bayesian"
)

const DefaultThresh = .95

var username string  // twitter username
var password string  // twitter password
var track *string    // comma-delimited list of tracking keywords for twitter api
var c *Classifier    // the classifier
var san *Sanitizer   // the sanitizer
var exclList *string // list of excluded terms
var count [2]int     // the count of all classifications
var highCount [2]int // the count of all learned classifications
var thresh *float64  // threshold for learning
var printOnly *bool  // suppress classification?
var loadFile *string

func init() {
	// command-line flags
	track = flag.String("track", "", "comma-separated list of tracking terms")
	thresh = flag.Float64("thresh", DefaultThresh, "the confidence threshold required to learn new content")
	exclList = flag.String("exclude", "", "comma-separated list of keywords excluded from classification")
	printOnly = flag.Bool("print-only", false, "only print the Tweets, do not classify them")
	loadFile = flag.String("load-file", "", "specify classifier file")
	flag.Parse()

	// read the arguments
	args := flag.Args()
	if len(args) != 2 {
		println("Usage: [--help|<options>...] <username> <password>")
		os.Exit(1)
	}
	username = args[0]
	password = args[1]

	// load and train the classifier
	if *loadFile != "" {
		// from a file
		var err os.Error
		c, err = NewClassifierFromFile(*loadFile)
		if err != nil {
			println("error loading:", err.String())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "classifier is loaded: %v\n", c.WordCount())
	} else {
		// from scratch
		c = NewClassifier(Positive, Negative)
		LearnFile(c, "data/positive.txt", Positive)
		LearnFile(c, "data/negative.txt", Negative)
		fmt.Fprintf(os.Stderr, "classifier is trained: %v\n", c.WordCount())
	}

	// init the sanitizer
	excl := strings.Split(*exclList, ",")
	if *exclList != "" {
		fmt.Fprintf(os.Stderr, "excluding: %v\n", excl)
	}
	stopWords := ReadFile("data/stopwords.txt")
	fmt.Fprintf(os.Stderr, "stop words: %v\n", stopWords)
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

	// listen for Ctrl-C
	go signalHandler()
}

func signalHandler() {
	for {
		sig := <-signal.Incoming
		if strings.HasPrefix(sig.String(), "SIGTSTP") {
			// ctrl-Z
			t := time.LocalTime()
			name := t.Format("15-04-05") + ".data"
			println("\nsaving classifier to", name)
			err := c.WriteToFile(name)
			if err != nil {
				println("error", err)
			}
			os.Exit(0)
		} else if strings.HasPrefix(sig.String(), "SIGINT") {
			// Ctrol-C
			println("\nstopping without save")
			os.Exit(0)
		}
	}
}

func main() {
	// stream Twitter
	stream := make(chan *twitterstream.Tweet)
	client := twitterstream.NewClient(username, password)
	tracks := strings.Split(*track, ",")

	err := client.Track(tracks, stream)
	if err != nil {
		println(err.String())
	}
	fmt.Fprintf(os.Stderr, "track = %v\n", *track)

	// process the tweets
	for {
		tw := (<-stream).Text
		if !*printOnly {
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
	scores, inx, _ := c.ProbScores(doc)
	logScores, logInx, _ := c.LogScores(doc)
	class := c.Classes[inx]
	logClass := c.Classes[logInx]
	count[inx]++

	// the rate of positive sentiment
	posrate := float64(count[0]) / float64(count[0]+count[1])
	learned := ""

	// if above the threshold, then learn
	// this document
	if scores[inx] > *thresh {
		c.Learn(doc, class)
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

// pretty print all the classification information
func prettyPrintDoc(doc []string) {
	//fmt.Printf("\n%v\n", doc)
	fmt.Printf("\t")
	for _, word := range doc {
		fmt.Printf("%7s", abbrev(word, 5))
	}
	fmt.Println("")

	freqs := c.WordFrequencies(doc)
	for i := 0; i < 2; i++ {
		fmt.Printf("%6s", c.Classes[i])
		for j := 0; j < len(doc); j++ {
			fmt.Printf("%7.4f", freqs[i][j])
		}
		fmt.Println("")
	}
}

// abbrev will abreavate a word with ".." if it
// is too long. It is used for display purposes.
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
