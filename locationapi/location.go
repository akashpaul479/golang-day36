package locationapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	City                 string  `json:"city"`
	Locality             string  `json:"locality"`
	PrincipalSubdivision string  `json:"principalSubdivision"`
	CountryName          string  `json:"countryName"`
	CountryCode          string  `json:"countryCode"`
	Continent            string  `json:"continent"`
	LookupSource         string  `json:"lookupSource"`
}

const (
	bigdatacloud = "https://api.bigdatacloud.net/data/reverse-geocode-client"
)

func LocationHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if lat == "" || lon == "" {
		http.Error(w, "lat and lon query parameters is required", http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&localityLanguage=en", bigdatacloud, lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to call", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "bigdatacloud returned error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	var location LocationResponse
	if err := json.Unmarshal(body, &location); err != nil {
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)

}

// main func
func LocationAPI() {
	http.HandleFunc("/api/location", LocationHandler)

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
