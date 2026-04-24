package auth

import "testing"

func TestHash(t *testing.T) {
	password := "12345"
	hashedPassword, _ := HashPassword(password)
	got, _ := CheckPasswordHash(password, hashedPassword)
	if got != true {
		t.Errorf(`HashPassword(password) = %v; want true`, got)
	}
}