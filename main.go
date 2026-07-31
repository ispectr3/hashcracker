// ============================================================================
// Hash Cracker - Ferramenta de Identificação e Cracking de Hashes
// Autor: kaique
// Licença: MIT
// Descrição: Identifica tipos de hash e realiza ataques de dicionário
//            e brute force usando goroutines para performance.
//   USE APENAS PARA FINS EDUCACIONAIS E TESTES AUTORIZADOS!
// ============================================================================

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kaique/hashcracker-go/cracker"
)

// Cores ANSI para output no terminal
const (
	Red     = "\033[0;31m"
	Green   = "\033[0;32m"
	Yellow  = "\033[1;33m"
	Blue    = "\033[0;34m"
	Cyan    = "\033[0;36m"
	Magenta = "\033[0;35m"
	Bold    = "\033[1m"
	Reset   = "\033[0m"
)

// banner exibe a arte ASCII do programa
func banner() {
	fmt.Println(Cyan)
	fmt.Println(`        `)
	fmt.Println(`      `)
	fmt.Println(`  `)
	fmt.Println(`  `)
	fmt.Println(`        `)
	fmt.Println(`        `)
	fmt.Println()
	fmt.Println(`         `)
	fmt.Println(`   `)
	fmt.Println(`               `)
	fmt.Println(`               `)
	fmt.Println(`          `)
	fmt.Println(`            `)
	fmt.Println(Reset)
	fmt.Printf("  %s%s Hash Cracker v1.0%s\n", Bold, Magenta, Reset)
	fmt.Printf("  %s Ferramenta de Identificação e Cracking de Hashes%s\n", Blue, Reset)
	fmt.Printf("  %s  Use apenas para fins educacionais e testes autorizados!%s\n", Yellow, Reset)
	fmt.Println()
	fmt.Printf("  %s%s\n\n", Cyan, Reset)
}

func main() {
	// Definição das flags de linha de comando
	hashStr := flag.String("hash", "", "Hash único para crackear")
	hashFile := flag.String("file", "", "Arquivo contendo múltiplos hashes (um por linha)")
	wordlist := flag.String("wordlist", "", "Caminho para a wordlist")
	mode := flag.String("mode", "dict", "Modo de ataque: dict (dicionário) ou brute (brute force)")
	charset := flag.String("charset", "abcdefghijklmnopqrstuvwxyz0123456789", "Charset para brute force")
	minLen := flag.Int("min-len", 1, "Comprimento mínimo para brute force")
	maxLen := flag.Int("max-len", 6, "Comprimento máximo para brute force")
	threads := flag.Int("threads", 8, "Número de goroutines (threads)")
	flag.Parse()

	banner()

	// Validações de entrada
	if *hashStr == "" && *hashFile == "" {
		fmt.Printf("%s[] Especifique um hash (-hash) ou arquivo de hashes (-file)%s\n", Red, Reset)
		fmt.Printf("%s[] Use -h para ver todas as opções%s\n", Yellow, Reset)
		os.Exit(1)
	}

	// Coleta os hashes para crackear
	var hashes []string
	if *hashStr != "" {
		hashes = append(hashes, *hashStr)
	}
	if *hashFile != "" {
		fileHashes, err := loadHashesFromFile(*hashFile)
		if err != nil {
			fmt.Printf("%s[] Erro ao ler arquivo de hashes: %v%s\n", Red, err, Reset)
			os.Exit(1)
		}
		hashes = append(hashes, fileHashes...)
	}

	fmt.Printf("%s[] Total de hashes para crackear: %d%s\n\n", Cyan, len(hashes), Reset)

	// Processa cada hash
	for i, h := range hashes {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}

		fmt.Printf("%s%s\n", Cyan, Reset)
		fmt.Printf("%s[] Hash %d/%d: %s%s\n", Bold, i+1, len(hashes), h, Reset)

		// Identifica o tipo do hash
		hashType := cracker.IdentifyHash(h)
		fmt.Printf("%s[ ] Tipo identificado: %s%s%s\n", Blue, Bold, hashType, Reset)

		if hashType == "Desconhecido" {
			fmt.Printf("%s[ ] Tipo de hash não reconhecido, pulando...%s\n\n", Yellow, Reset)
			continue
		}

		startTime := time.Now()

		// Executa o ataque baseado no modo selecionado
		switch *mode {
		case "dict":
			// Ataque de dicionário
			if *wordlist == "" {
				// Tenta usar wordlist padrão
				defaultWL := "wordlists/common.txt"
				if _, err := os.Stat(defaultWL); err == nil {
					*wordlist = defaultWL
				} else {
					fmt.Printf("%s[] Especifique uma wordlist com -wordlist%s\n", Red, Reset)
					continue
				}
			}

			fmt.Printf("%s[] Modo: Dicionário | Wordlist: %s | Threads: %d%s\n", Blue, *wordlist, *threads, Reset)
			fmt.Printf("%s[] Iniciando ataque...%s\n\n", Yellow, Reset)

			result, attempts, err := cracker.DictionaryAttack(h, hashType, *wordlist, *threads)
			elapsed := time.Since(startTime)

			if err != nil {
				fmt.Printf("%s[] Erro: %v%s\n", Red, err, Reset)
			} else if result != "" {
				fmt.Printf("\n%s%s[] HASH CRACKEADA!%s\n", Bold, Green, Reset)
				fmt.Printf("%s    Hash:     %s%s\n", Green, h, Reset)
				fmt.Printf("%s    Senha:    %s%s%s\n", Green, Bold, result, Reset)
				fmt.Printf("%s    Tipo:     %s%s\n", Green, hashType, Reset)
			} else {
				fmt.Printf("\n%s[] Senha não encontrada na wordlist%s\n", Red, Reset)
			}

			speed := float64(attempts) / elapsed.Seconds()
			fmt.Printf("%s[] %d tentativas em %v (%.0f hash/s)%s\n", Cyan, attempts, elapsed.Round(time.Millisecond), speed, Reset)

		case "brute":
			// Ataque de brute force
			fmt.Printf("%s[] Modo: Brute Force | Charset: %d chars | Len: %d-%d | Threads: %d%s\n",
				Blue, len(*charset), *minLen, *maxLen, *threads, Reset)
			fmt.Printf("%s[] Iniciando ataque...%s\n\n", Yellow, Reset)

			result, attempts, err := cracker.BruteForceAttack(h, hashType, *charset, *minLen, *maxLen, *threads)
			elapsed := time.Since(startTime)

			if err != nil {
				fmt.Printf("%s[] Erro: %v%s\n", Red, err, Reset)
			} else if result != "" {
				fmt.Printf("\n%s%s[] HASH CRACKEADA!%s\n", Bold, Green, Reset)
				fmt.Printf("%s    Hash:     %s%s\n", Green, h, Reset)
				fmt.Printf("%s    Senha:    %s%s%s\n", Green, Bold, result, Reset)
				fmt.Printf("%s    Tipo:     %s%s\n", Green, hashType, Reset)
			} else {
				fmt.Printf("\n%s[] Senha não encontrada no espaço de busca%s\n", Red, Reset)
			}

			speed := float64(attempts) / elapsed.Seconds()
			fmt.Printf("%s[] %d tentativas em %v (%.0f hash/s)%s\n", Cyan, attempts, elapsed.Round(time.Millisecond), speed, Reset)

		default:
			fmt.Printf("%s[] Modo inválido: %s (use: dict ou brute)%s\n", Red, *mode, Reset)
		}

		fmt.Println()
	}

	fmt.Printf("%s%s[] Processamento concluído!%s\n\n", Bold, Magenta, Reset)
}

// loadHashesFromFile carrega hashes de um arquivo (um por linha)
func loadHashesFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hashes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			hashes = append(hashes, line)
		}
	}
	return hashes, scanner.Err()
}
