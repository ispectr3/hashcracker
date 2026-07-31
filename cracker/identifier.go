// ============================================================================
// identifier.go - Módulo de Identificação de Hash
// Identifica o tipo de hash baseado no comprimento e formato.
// ============================================================================

package cracker

import (
	"regexp"
	"strings"
)

// IdentifyHash identifica o tipo de hash baseado no comprimento e padrão
// Retorna uma string com o tipo identificado (MD5, SHA1, SHA256, SHA512, bcrypt)
func IdentifyHash(hash string) string {
	hash = strings.TrimSpace(hash)

	// Verifica se é bcrypt pelo prefixo característico ($2a$, $2b$, $2y$)
	if isBcrypt(hash) {
		return "bcrypt"
	}

	// Verifica se é um hash hexadecimal válido
	if !isHex(hash) {
		return "Desconhecido"
	}

	// Identifica pelo comprimento da string hexadecimal
	switch len(hash) {
	case 32:
		// MD5 produz hashes de 128 bits = 32 caracteres hex
		return "MD5"
	case 40:
		// SHA1 produz hashes de 160 bits = 40 caracteres hex
		return "SHA1"
	case 64:
		// SHA256 produz hashes de 256 bits = 64 caracteres hex
		return "SHA256"
	case 128:
		// SHA512 produz hashes de 512 bits = 128 caracteres hex
		return "SHA512"
	default:
		return "Desconhecido"
	}
}

// isHex verifica se uma string contém apenas caracteres hexadecimais válidos
func isHex(s string) bool {
	match, _ := regexp.MatchString("^[a-fA-F0-9]+$", s)
	return match
}

// isBcrypt verifica se o hash tem formato bcrypt
// Formato bcrypt: $2a$XX$... ou $2b$XX$... ou $2y$XX$...
func isBcrypt(hash string) bool {
	match, _ := regexp.MatchString(`^\$2[aby]\$\d{2}\$.{53}$`, hash)
	return match
}
