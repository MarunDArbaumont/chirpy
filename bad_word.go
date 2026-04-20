package main

import (
	"strings"
)

func replaceBadWord(body string, badWords map[string]struct{}) string {
	splitedBody := strings.Split(body, " ")
	for i, word := range splitedBody {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			splitedBody[i] = "****"
		}
	}
	cleanedBody := strings.Join(splitedBody, " ")
	return cleanedBody
}