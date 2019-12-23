package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Quote struct {
	Author    string `json:"author"`
	ID        int    `json:"id"`
	Quote     string `json:"quote"`
	Permalink string `json:"permalink"`
}

type RandomQuote struct {
	Success struct {
		Total int `json:"total"`
	} `json:"success"`
	Contents struct {
		Quotes    []Quote `json:"quotes"`
		Copyright string  `json:"copyright"`
	} `json:"contents"`
}

func (quote *Quote) printQuote() {
	fmt.Printf("❝ %s ❞\n", quote.Quote)
	fmt.Printf("\t ― %s\n", quote.Author)
	fmt.Printf("\nX: %s\n", quote.Permalink)
}

func printQuotes(quotes []Quote) {
	for _, quote := range quotes {
		quote.printQuote()
	}
}

type CachedQuotes []Quote

func (c CachedQuotes) printRandomQuotes() {
	size := len(c)
	if size > 0 {
		c[rand.Intn(size)].printQuote()
	}
}

func getRandomQuote() (*Quote, error) {
	url := "http://quotes.stormconsultancy.co.uk/random.json"
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "deepakshrma-random-quote")
	res, _ := spaceClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	quote := &Quote{}
	jsonErr := json.Unmarshal(body, quote)
	if jsonErr != nil {
		return nil, fmt.Errorf("Error while getting quote: %v", jsonErr)
	}
	return quote, nil
}

func printMore() {
	fmt.Println("##\n1. For More\n0. Exit \n##")
	fmt.Print("-> ")
}

func main() {
	cachedQuotes := CachedQuotes{}
	// var cachedQuotes CachedQuotes
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("## Random Quote ##")
	fmt.Println("---------------------")
	exit := false
	for !exit {
		printMore()
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("exit", text) == 0 || strings.Compare("0", text) == 0 {
			fmt.Println("Thanks, Have a good day!")
			exit = true
		}
		if !exit {
			ranQuote, err := getRandomQuote()
			if err != nil {
				log.Printf("Error while generating random quote: %v", err)
				cachedQuotes.printRandomQuotes()
			} else {
				ranQuote.printQuote()
				cachedQuotes = append(cachedQuotes, *ranQuote)
			}
		}
		/*
			//Reading Single UTF-8 Encoded Unicode Characters
			char, _, err := reader.ReadRune()

			if err != nil {
				fmt.Println(err)
			}
			// print out the unicode value i.e. A -> 65, a -> 97
			fmt.Println(char)

			switch char {
			case 'A':
				fmt.Println("A Key Pressed")
				break
			case 'a':
				fmt.Println("a Key Pressed")
				break
			}
		*/
	}

}
