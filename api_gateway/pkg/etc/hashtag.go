package etc

import (
	"context"
	"strings"
	"unicode"
)

func FormatHashtag(hashtag string) string {
	if len(hashtag) == 0 {
		return ""
	}

	firstLetter := string(unicode.ToUpper(rune(hashtag[0])))
	rest := strings.ToLower(hashtag[1:])

	return firstLetter + rest
}

func GetTags(ctx context.Context, content string) ([]string, error) {
	var (
		response []string
	)

	words := strings.Fields(content)

	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			parts := strings.Split(word, "#")
			for _, part := range parts {
				if part != "" {
					cleanedWord := strings.TrimFunc(part, func(r rune) bool {
						return !unicode.IsLetter(r) && !unicode.IsDigit(r)
					})
					if cleanedWord != "" {
						formattedWord := FormatHashtag(cleanedWord)
						response = append(response, "#"+formattedWord)
					}
				}
			}
		}
	}

	return response, nil
}
