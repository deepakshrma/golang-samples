package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/urfave/cli"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}
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

func (quote *Quote) printQuote() string {
	runes := []rune(quote.Quote)
	len := math.Min(float64(len(runes)), 230)
	quoteStr := string(runes[0:int(len)])
	quoteStr = fmt.Sprintf("❝ %s ❞\n  ― %s X Tweet From Bot::xdeepakv", quoteStr, quote.Author)
	return quoteStr
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
func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}
func main() {
	fmt.Println("Go-Twitter Bot v0.01")
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}
	app := cli.NewApp()
	//.Run(os.Args)
	app.Name = "Twitter Bot "
	app.Usage = "Let's Tweet together!"
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// We'll be using the same flag for all our commands
	// so we'll define it up here
	searchKey := cli.StringFlag{
		Name:  "key",
		Value: "golang",
	}
	searchKeys := append([]cli.Flag{}, searchKey)
	// we create our commands
	app.Commands = []cli.Command{
		{
			Name:  "post",
			Usage: "Post random quote on twitter",
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				// a simple lookup function
				// Post here
				ranQuote, err := getRandomQuote()
				if err == nil {
					quote := ranQuote.printQuote()
					fmt.Println(quote)
					client, err := getClient(&creds)
					if err != nil {
						log.Println("Error getting Twitter Client")
						log.Fatal(err)
					}
					tweet, resp, err := client.Statuses.Update(quote, nil)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("%+v\n", resp)
					log.Printf("%+v\n", tweet)
				}

				return nil
			},
		},
		{
			Name:  "search",
			Usage: "Search tweet by key",
			Flags: searchKeys,
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				// a simple lookup function
				// Post here
				client, err := getClient(&creds)
				if err != nil {
					log.Println("Error getting Twitter Client")
					log.Fatal(err)
				}
				key := c.String(searchKey.Name)

				search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
					Query: key,
				})

				if err != nil {
					log.Print(err)
				}

				log.Printf("%+v\n", resp)
				log.Printf("%+v\n", search)
				return nil
			},
		},
	}

	// start our application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
