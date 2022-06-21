package main

import(
	"encoding/json"
	"fmt"
)

func getGameList(data string) ([]Game, error) {
	
	resp := TwitchGameResponse{}
	
	err := json.Unmarshal([]byte(data), &resp)

	if err != nil {
		fmt.Println("Failed to parse data from Twtich response.")
		return nil, err
	}
	
	return resp.Data, nil
}