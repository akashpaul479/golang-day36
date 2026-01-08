package dictionaryapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Dictionaries API response
type Dictionary []struct {
	Word     string `json:"word"`
	Phonetic string `json:"phonetic"`
	Origin   string `json:"origin,omitempty"`
	Meanings []struct {
		Partofspeech string `json:"partofspeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example,omitempty"`
		} `json:"definitions"`
	} `json:"meanings"`
}

// Handler to fetch word meaning
func DictionaryHandler(w http.ResponseWriter, r *http.Request) {
	Word := r.URL.Query().Get("word")

	if Word == "" {
		http.Error(w, "missing word", http.StatusBadRequest)
		return
	}
	// Build API url
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", Word)

	// Call dictionary API
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode json
	var dict Dictionary
	if err := json.NewDecoder(resp.Body).Decode(&dict); err != nil || len(dict) == 0 {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}
	// prepare response
	result := map[string]interface{}{
		"word":     dict[0].Word,
		"phonetic": dict[0].Phonetic,
		"origin":   dict[0].Origin,
		"Meanings": dict[0].Meanings,
	}
	// pretty print
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		http.Error(w, "Failed to format json", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// main func
func DictionaryAPI() {
	http.HandleFunc("/dictionary", DictionaryHandler)

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
