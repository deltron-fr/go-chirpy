package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	secret := "Development"
	expiration := 1 * time.Second

	token, err := MakeJWT(userID, secret, expiration)
	if err != nil {
    	t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	haveUserID, err := ValidateJWT(token, secret)
	if err != nil {
    	t.Fatalf("ValidateJWT() error = %v, want nil", err)
	}


	if haveUserID != userID {
		t.Errorf("userID = %v, want match for %v", haveUserID, userID)
	}
}


func TestPassword(t *testing.T) {
    password := "testing"

    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword() error = %v", err)
    }

    match, err := CheckPasswordHash(password, hash)
    if err != nil {
        t.Fatalf("CheckPasswordHash() error = %v", err)
    }

    if !match {
        t.Errorf("expected passwords to match, got mismatch")
    }
}


func TestGetBearerToken(t *testing.T) {
	want := "TOKEN"
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("request error = %v", err)
	}

	req.Header.Set("Authorization", "Bearer TOKEN")
	
	have, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("GetBearerToken() error = %v", err)
	}

	if want != have {
		t.Errorf("token = %s, want match for %s", have, want)
	}

}