package main

import "fmt"
import "os"

var username string
var password string

func init() {
        if len(os.Args) != 3 {
                    fmt.Print("Usage: <username> <password>\n")
                            os.Exit(1)
                                }
                                    username = os.Args[1]
                                        password = os.Args[2]
                                            fmt.Printf("username: %v password: %v\n", username, password)
}

func main() {

        stream := make(chan *twitterstream.Tweet)
            client := twitterstream.NewClient("username", "password")

                err := client.Track([]string{"#OWS"}, stream)
                    if err != nil {
                                println(err.String())
                                    }
                                        for {
                                                    tw := <-stream
                                                            println(tw.User.Screen_name, ": ", tw.Text)
                                                                }
}
