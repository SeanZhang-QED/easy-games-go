package twitch

import (
	"fmt"
	"github.com/SeanZhang-QED/easy-games-go/config"
	"io"
	"net/http"
	"time"
)

func SearchByName(url string, gameName string) (string, error) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}

	req, _ := http.NewRequest("GET", url, nil)

	if gameName != "" {
		q := req.URL.Query()          // Get a copy of the query values.
		q.Add("name", gameName)       // Add a new value to the set.
		req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
	}

	req.Header.Add("Authorization", config.TOKEN)
	req.Header.Add("Client-Id", config.CLIENT_ID)

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get result from Twitch API.")
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Failed to read result from Twitch response body.")
		return "", err
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}
