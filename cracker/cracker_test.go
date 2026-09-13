package cracker

import (
	"testing"
)

func TestIdentifyHash(t *testing.T) {
	tests := []struct {
		hash     string
		expected string
	}{
		{"5d41402abc4b2a76b9719d911017c592", "MD5"},
		{"2fd4e1c67a2d28fced849ee1bb76e7391b93eb12", "SHA1"},
		{"2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", "SHA256"},
		{"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e", "SHA512"},
		{"$2a$12$e8N4ndr5k/W5W9W5W5W5W.5W5W5W5W5W5W5W5W5W5W5W5W5W5W5W5", "bcrypt"},
		{"not-a-valid-hash", "Desconhecido"},
	}

	for _, tt := range tests {
		result := IdentifyHash(tt.hash)
		if result != tt.expected {
			t.Errorf("IdentifyHash(%s) = %s, expected %s", tt.hash, result, tt.expected)
		}
	}
}

func TestHashPassword(t *testing.T) {
	word := "hello"
	md5Expected := "5d41402abc4b2a76b9719d911017c592"
	sha256Expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	if HashPassword(word, "MD5") != md5Expected {
		t.Errorf("MD5 mismatch for %s", word)
	}
	if HashPassword(word, "SHA256") != sha256Expected {
		t.Errorf("SHA256 mismatch for %s", word)
	}
}
