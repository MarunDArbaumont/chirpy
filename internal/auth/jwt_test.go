package auth

import (
	"testing"
	"time"
	"net/http"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	parsedID, _ := uuid.Parse("a8e9de63-a35c-4b40-9606-2f89854cc3e3")
	time, _ := time.ParseDuration("45s")
	token, _ := MakeJWT(parsedID ,"my-test-secret", time)
	_, err := ValidateJWT(token, "my-test-secret-but-not-the-right")
	if err == nil {
		t.Error("there should be an error here")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer THIS_IS_THE_TOKEN")
	got, _ := GetBearerToken(header)
	if got != "THIS_IS_THE_TOKEN" {
		t.Errorf(`GetBearerToken(header) = %v should be "THIS_IS_THE_TOKEN"`, got)
	}
}
