package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ApiResponse []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio,omitempty"`
	} `json:"phonetics"`
	Origin   string `json:"origin"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
			Synonyms   []any  `json:"synonyms"`
			Antonyms   []any  `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no word to define")
		os.Exit(1)
	}

	resp, err := http.Get(fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", os.Args[1]))

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("no definition found")
		} else {
			fmt.Println("failed to define the word")
		}

		os.Exit(0)
	}

	if err != nil {
		fmt.Println("failed to define the word")
		os.Exit(1)
	}

	var data ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		fmt.Println("failed to read api response")
	}

	prettyPrint(&data)
}

func prettyPrint(data *ApiResponse) {
	if data == nil || len(*data) == 0 {
		fmt.Println("No data to print.")
		return
	}

	fmt.Printf("Word: %s\n", (*data)[0].Word)

	if len((*data)[0].Phonetics) > 0 {
		fmt.Println("Phonetics:")
		for _, phonetic := range (*data)[0].Phonetics {
			fmt.Printf("  - %s\n", phonetic.Text)
		}
	}

	fmt.Printf("Origin: %s\n", (*data)[0].Origin)

	if len((*data)[0].Meanings) > 0 {
		fmt.Println("Meanings:")
		for _, meaning := range (*data)[0].Meanings {
			fmt.Printf("  %s:\n", meaning.PartOfSpeech)
			for _, definition := range meaning.Definitions {
				fmt.Printf("    - Definition: %s\n", definition.Definition)
				if definition.Example != "" {
					fmt.Printf("      Example: %s\n", definition.Example)
				}

				if len(definition.Synonyms) > 0 {
					fmt.Printf("      Synonyms: %s\n", strings.Join(toStringSlice(definition.Synonyms), ", "))
				}

				if len(definition.Antonyms) > 0 {
					fmt.Printf("      Antonyms: %s\n", strings.Join(toStringSlice(definition.Antonyms), ", "))
				}
			}
		}
	}
}

func toStringSlice(slice []interface{}) []string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}
	return strSlice
}
