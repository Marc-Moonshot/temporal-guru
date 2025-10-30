package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Pass at least 1 argument.")
		return
	}

	var species = os.Args[1]

	url := `https://pokeapi.co/api/v2/pokemon/` + species
	fmt.Println(url)

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, getErr := httpClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	var mon map[string]any
	json.Unmarshal(body, &mon)

	pretty, _ := json.MarshalIndent(mon, "", "  ")
	fmt.Println(string(pretty))
}
