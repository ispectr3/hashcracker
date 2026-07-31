// ============================================================================
// dictionary.go - Módulo de Ataque por Dicionário
// Realiza ataque de dicionário usando wordlist e goroutines para
// paralelizar a comparação de hashes.
// ============================================================================

package cracker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// DictionaryAttack realiza um ataque de dicionário contra um hash
// Parâmetros:
//   - targetHash: hash que queremos crackear
//   - hashType: tipo do hash (MD5, SHA1, SHA256, SHA512)
//   - wordlistPath: caminho para o arquivo de wordlist
//   - numWorkers: número de goroutines paralelas
//
// Retorna: senha encontrada, número de tentativas, e erro (se houver)
func DictionaryAttack(targetHash, hashType, wordlistPath string, numWorkers int) (string, int64, error) {
	// Abre o arquivo de wordlist
	file, err := os.Open(wordlistPath)
	if err != nil {
		return "", 0, fmt.Errorf("erro ao abrir wordlist: %v", err)
	}
	defer file.Close()

	// Normaliza o hash alvo para lowercase
	targetHash = strings.ToLower(strings.TrimSpace(targetHash))

	// Canal para enviar palavras da wordlist para os workers
	wordChan := make(chan string, numWorkers*10)

	// Canal para receber o resultado (senha encontrada)
	resultChan := make(chan string, 1)

	// Contador atômico de tentativas (thread-safe)
	var attempts int64

	// Flag atômica para sinalizar que a senha foi encontrada
	var found int32

	// WaitGroup para aguardar todos os workers finalizarem
	var wg sync.WaitGroup

	// Inicia os workers (goroutines)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Cada worker consome palavras do canal
			for word := range wordChan {
				// Se já encontrou, para de processar
				if atomic.LoadInt32(&found) == 1 {
					return
				}

				// Incrementa o contador de tentativas
				atomic.AddInt64(&attempts, 1)

				// Gera o hash da palavra e compara com o alvo
				computed := HashPassword(word, hashType)
				if computed == targetHash {
					// Senha encontrada! Sinaliza para os outros workers
					atomic.StoreInt32(&found, 1)
					resultChan <- word
					return
				}

				// Exibe progresso a cada 100.000 tentativas
				current := atomic.LoadInt64(&attempts)
				if current%100000 == 0 {
					fmt.Printf("\r\033[K  [] %d senhas testadas...", current)
				}
			}
		}()
	}

	// Goroutine que lê a wordlist e envia palavras para o canal
	go func() {
		scanner := bufio.NewScanner(file)
		// Aumenta o buffer para suportar linhas longas
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		for scanner.Scan() {
			// Se já encontrou, para de ler
			if atomic.LoadInt32(&found) == 1 {
				break
			}

			word := strings.TrimSpace(scanner.Text())
			if word != "" {
				wordChan <- word
			}
		}
		// Fecha o canal quando terminar de ler
		close(wordChan)
	}()

	// Goroutine que aguarda os workers e fecha o canal de resultado
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Aguarda o resultado
	result := <-resultChan
	finalAttempts := atomic.LoadInt64(&attempts)

	// Limpa a linha de progresso
	fmt.Print("\r\033[K")

	return result, finalAttempts, nil
}
