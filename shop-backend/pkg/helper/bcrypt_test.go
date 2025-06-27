package helper

import (
	"fmt"
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHashAndCompare(t *testing.T) {
	password := "securepass"

	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	log.Println("Gnerated Hash:", string(hash))

	// Compare correct password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("securepass"))
	if err != nil {
		t.Error("❌ Password does not match (should match)")
	} else {
		fmt.Println("✅ Password matched (Expected)")
	}

	// Compare wrong password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("wrongpass"))
	if err != nil {
		fmt.Println("✅ Correctly failed wrong password")
	} else {
		t.Error("❌ Incorrectly matched wrong password")
	}
}
