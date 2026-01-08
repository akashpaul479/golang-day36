package jokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JokeResponse struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}
type JokeOutput struct {
	Category  string `json:"category"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func JokeHandler(w http.ResponseWriter, r *http.Request) {

	url := "https://official-joke-api.appspot.com/random_joke"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "failed to fetch joke", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Joke API error", http.StatusBadGateway)
		return
	}

	var joke JokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&joke); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	result := JokeOutput{
		Category:  joke.Type,
		Setup:     joke.Setup,
		Punchline: joke.Punchline,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

// main func
func JokeAPI() {
	http.HandleFunc("/joke", JokeHandler)

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
