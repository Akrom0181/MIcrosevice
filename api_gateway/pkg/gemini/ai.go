package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"user_api_gateway/api/models"
)

func AskFromGemini(content string, defaultTags []string) ([]string, error) {

	// Get the API key from the environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return []string{}, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	promt := fmt.Sprintf("Given the content: '%s', which of the following tags are most relevant? %v. Just write each tag only once.", content, defaultTags)

	requestBody := models.RequestBody{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: promt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// Send a POST request to the API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the response
	var response models.GeminiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	// Get the first candidate's content
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		// Get the first candidate's content
		aiResponse := response.Candidates[0].Content.Parts[0].Text

		// Split the response by commas
		matchedTags := []string{}
		for _, tag := range defaultTags {
			if strings.Contains(aiResponse, tag[1:]) {
				matchedTags = append(matchedTags, tag)
			}
		}
		return matchedTags, nil
	}

	return []string{}, fmt.Errorf("no valid response from AI")
}
