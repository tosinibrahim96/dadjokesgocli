/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a randolm dad joke",
	Long:  `This command will fetch a random dad joke from the API and display it to the user`,
	Run: func(cmd *cobra.Command, args []string) {
		term, _ := cmd.Flags().GetString("term")

		if term!= "" {
            getRandomJokeWithTerm(term)
        }else{
			getRandomJoke()
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.PersistentFlags().String("term", "", "A search term for a dad joke")
}

type Joke struct {
	ID     int    `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResult struct {
	Results json.RawMessage `json:"results"`
	Searchterm string `json:"search_term"`
	Status int `json:"status"`
	TotalJokes int `json:"total_jokes"`
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	res := getJokeData(url)
	joke := Joke{}

	if err := json.Unmarshal(res, &joke); err != nil {
		log.Printf("Error unmarshalling joke data - %v", err)
	}

	fmt.Println(string(joke.Joke))
}

func getRandomJokeWithTerm(term string){
	total, results := getJokeDataWithTerm(term)
	randomiseJokeList(total, results)
}

func randomiseJokeList(length int, jokeList []Joke){

	rand.Seed(time.Now().Unix())

	min :=0
	max :=length -1

	if length < min {
		err := fmt.Errorf("No jokes found for search term")
		fmt.Println(err.Error())
	}else{
		randomNum := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNum].Joke)
	}
}

func getJokeDataWithTerm(term string) (totalJokes int, jokeList []Joke) {

	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", url.QueryEscape(term))
	res := getJokeData(url)
	searchResult := SearchResult{}

	if err := json.Unmarshal(res, &searchResult); err != nil {
		log.Printf("Error unmarshalling joke data - %v", err)
	}

	jokes := []Joke{}
	if err := json.Unmarshal(searchResult.Results, &jokes); err != nil {
		log.Printf("Error unmarshalling jokes results data - %v", err)
	}

	return searchResult.TotalJokes, jokes

}

func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(http.MethodGet, baseAPI, nil)

	if err != nil {
		log.Printf("Error creating request - %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "DadJokesGoCLI (github.com/tosinibrahim96/dadjokesgocli)")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
        log.Printf("Error processing request - %v", err)
    }

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Printf("Error reading response - %v", err)
	}

	return body
}
