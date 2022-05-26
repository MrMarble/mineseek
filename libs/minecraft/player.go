package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
	Data    responseData `json:"data"`
}

type responseData struct {
	Player player `json:"player"`
}
type player struct {
	Username string `json:"username"`
	UUID     string `json:"id"`
}

func GetUUID(name string) (string, error) {
	resp, err := http.Get("https://playerdb.co/api/player/minecraft/" + name)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var player response
	err = json.Unmarshal(body, &player)
	if err != nil {
		return "", err
	}

	if player.Success {
		return player.Data.Player.UUID, nil
	}

	return "", fmt.Errorf("Player not found")
}
