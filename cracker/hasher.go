// ============================================================================
// hasher.go - Módulo de Funções de Hash
// Contém funções para gerar hashes MD5, SHA1, SHA256 e SHA512.
// ============================================================================

package cracker

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// HashPassword gera o hash de uma senha usando o algoritmo especificado
// Suporta: MD5, SHA1, SHA256, SHA512
func HashPassword(password string, hashType string) string {
	switch hashType {
	case "MD5":
		return hashMD5(password)
	case "SHA1":
		return hashSHA1(password)
	case "SHA256":
		return hashSHA256(password)
	case "SHA512":
		return hashSHA512(password)
	default:
		return ""
	}
}

// hashMD5 gera o hash MD5 de uma string
// MD5 produz um digest de 128 bits (16 bytes)
func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// hashSHA1 gera o hash SHA1 de uma string
// SHA1 produz um digest de 160 bits (20 bytes)
func hashSHA1(text string) string {
	hash := sha1.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// hashSHA256 gera o hash SHA256 de uma string
// SHA256 produz um digest de 256 bits (32 bytes)
func hashSHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// hashSHA512 gera o hash SHA512 de uma string
// SHA512 produz um digest de 512 bits (64 bytes)
func hashSHA512(text string) string {
	hash := sha512.Sum512([]byte(text))
	return hex.EncodeToString(hash[:])
}
