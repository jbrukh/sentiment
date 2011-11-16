package main

import "fmt"
import "os"
import "twitterstream"
import "flag"
import "strings"
import "sentiment"

var username string
var password string
var track *string

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

	hist := sentiment.NewHistogram()
	hist.Exclude(tracks)
    hist.Exclude(sentiment.CommonEnglish())
    hist.Exclude(sentiment.TwitterTrash())

	for {
		tw := <-stream
		text := sanitize(tw.Text)
		println(tw.User.Screen_name, ": ", text)
		hist.AbsorbText(text, " ")
		printPops(hist.MostPopular())
	}
}

func sanitize(text string) string {
    return strings.ToLower(text)
}

func printPops(pops sentiment.TokenPops) {
    fmt.Println("")
    for _, value := range pops[:5] {
        fmt.Printf("%v:%5d\n", value.Token, value.Pop)   
    }
    fmt.Println("")
}
