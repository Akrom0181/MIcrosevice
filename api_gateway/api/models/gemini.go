package models

type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type RequestBody struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}
