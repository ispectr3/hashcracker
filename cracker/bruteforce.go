// ============================================================================
// bruteforce.go - Módulo de Ataque por Brute Force
// Gera combinações de caracteres e compara hashes usando goroutines.
// Utiliza um padrão de geração recursiva com worker pool.
// ============================================================================

package cracker

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// BruteForceAttack realiza um ataque de brute force contra um hash
// Parâmetros:
//   - targetHash: hash que queremos crackear
//   - hashType: tipo do hash (MD5, SHA1, SHA256, SHA512)
//   - charset: conjunto de caracteres para gerar combinações
//   - minLen: comprimento mínimo das combinações
//   - maxLen: comprimento máximo das combinações
//   - numWorkers: número de goroutines paralelas
//
// Retorna: senha encontrada, número de tentativas, e erro (se houver)
func BruteForceAttack(targetHash, hashType, charset string, minLen, maxLen, numWorkers int) (string, int64, error) {
	// Normaliza o hash alvo
	targetHash = strings.ToLower(strings.TrimSpace(targetHash))

	// Converte o charset em slice de bytes para acesso indexado
	chars := []byte(charset)

	// Canal para enviar candidatos para os workers
	candidateChan := make(chan string, numWorkers*100)

	// Canal para receber o resultado
	resultChan := make(chan string, 1)

	// Contador atômico de tentativas
	var attempts int64

	// Flag de controle
	var found int32

	var wg sync.WaitGroup

	// Calcula o espaço total de busca para informar o usuário
	totalCombinations := int64(0)
	for l := minLen; l <= maxLen; l++ {
		combos := int64(1)
		for i := 0; i < l; i++ {
			combos *= int64(len(chars))
		}
		totalCombinations += combos
	}
	fmt.Printf("  [📊] Espaço de busca: %d combinações\n\n", totalCombinations)

	// Inicia os workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for candidate := range candidateChan {
				if atomic.LoadInt32(&found) == 1 {
					return
				}

				atomic.AddInt64(&attempts, 1)

				// Gera hash e compara
				computed := HashPassword(candidate, hashType)
				if computed == targetHash {
					atomic.StoreInt32(&found, 1)
					resultChan <- candidate
					return
				}

				// Progresso a cada 100.000 tentativas
				current := atomic.LoadInt64(&attempts)
				if current%100000 == 0 {
					progress := float64(current) / float64(totalCombinations) * 100
					fmt.Printf("\r\033[K  [⏳] %d/%d tentativas (%.1f%%)", current, totalCombinations, progress)
				}
			}
		}()
	}

	// Goroutine que gera todas as combinações
	go func() {
		for length := minLen; length <= maxLen; length++ {
			if atomic.LoadInt32(&found) == 1 {
				break
			}

			fmt.Printf("  [🔨] Testando comprimento: %d\n", length)

			// Gera combinações de forma iterativa usando "contagem" em base N
			// onde N = tamanho do charset
			generateCombinations(chars, length, candidateChan, &found)
		}
		close(candidateChan)
	}()

	// Aguarda resultado
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	result := <-resultChan
	finalAttempts := atomic.LoadInt64(&attempts)

	// Limpa a linha de progresso
	fmt.Print("\r\033[K")

	return result, finalAttempts, nil
}

// generateCombinations gera todas as combinações possíveis de um dado comprimento
// Usa um algoritmo iterativo baseado em "odômetro" (como um contador)
// Exemplo: charset "abc", length 2 -> aa, ab, ac, ba, bb, bc, ca, cb, cc
func generateCombinations(chars []byte, length int, out chan<- string, found *int32) {
	// indices[i] representa a posição atual no charset para cada posição da string
	indices := make([]int, length)
	charLen := len(chars)

	for {
		// Verifica se já encontrou
		if atomic.LoadInt32(found) == 1 {
			return
		}

		// Constrói a string candidata a partir dos índices atuais
		candidate := make([]byte, length)
		for i, idx := range indices {
			candidate[i] = chars[idx]
		}

		// Envia para o canal de workers
		out <- string(candidate)

		// Incrementa os índices (como um odômetro/conta-quilômetros)
		// Começa pela última posição e propaga o "carry"
		carry := true
		for i := length - 1; i >= 0 && carry; i-- {
			indices[i]++
			if indices[i] >= charLen {
				indices[i] = 0 // Reset e propaga carry
			} else {
				carry = false // Sem carry, para aqui
			}
		}

		// Se o carry chegou até o fim, todas as combinações foram geradas
		if carry {
			return
		}
	}
}
