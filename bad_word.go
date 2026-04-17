package main

import (
	"strings"
)

func replaceBadWord(body string, badWords []string) string {
	splitedBody := strings.Split(body, " ")
	for i, word := range splitedBody {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord  {
				splitedBody[i] = "****"
				continue
			}
		}
	}
	cleanedBody := strings.Join(splitedBody, " ")
	return cleanedBody
}