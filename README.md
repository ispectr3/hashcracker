# 🔓 Hash Cracker

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Security](https://img.shields.io/badge/CyberSecurity-Tool-red?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)

> Ferramenta de identificação e cracking de hashes escrita em Go. Utiliza goroutines para ataques paralelos de alta performance com suporte a dicionário e brute force.

---

## 🚀 Features

- 🏷️ **Identificação Automática** — Detecta MD5, SHA1, SHA256, SHA512 e bcrypt pelo formato e comprimento
- 📖 **Ataque por Dicionário** — Testa senhas a partir de wordlists usando worker pool
- 🔨 **Ataque Brute Force** — Gera combinações com charset e comprimento configuráveis
- ⚡ **Goroutines** — Processamento paralelo com número de threads ajustável
- 📊 **Indicador de Progresso** — Mostra tentativas, velocidade (hash/s) e porcentagem
- 📁 **Múltiplos Hashes** — Crackeia vários hashes de um arquivo de uma vez
- 🎨 **Output Colorido** — Interface bonita no terminal com cores ANSI

## 📋 Pré-requisitos

- [Go 1.21+](https://golang.org/dl/) instalado

## 📦 Instalação

```bash
git clone https://github.com/kaique/hashcracker-go.git
cd hashcracker-go
go build -o hashcracker main.go
```

## 🎯 Uso

### Identificar e crackear um hash com dicionário
```bash
./hashcracker -hash "5f4dcc3b5aa765d61d8327deb882cf99" -wordlist wordlists/common.txt
```

### Ataque brute force
```bash
./hashcracker -hash "e10adc3949ba59abbe56e057f20f883e" -mode brute -charset "abcdefghijklmnopqrstuvwxyz0123456789" -min-len 1 -max-len 6
```

### Crackear múltiplos hashes de um arquivo
```bash
./hashcracker -file hashes.txt -wordlist wordlists/common.txt -threads 16
```

### Flags

| Flag | Descrição | Padrão |
|------|-----------|--------|
| `-hash` | Hash único para crackear | — |
| `-file` | Arquivo com múltiplos hashes | — |
| `-wordlist` | Caminho para a wordlist | `wordlists/common.txt` |
| `-mode` | Modo: `dict` ou `brute` | `dict` |
| `-charset` | Caracteres para brute force | `a-z0-9` |
| `-min-len` | Comprimento mínimo (brute force) | `1` |
| `-max-len` | Comprimento máximo (brute force) | `6` |
| `-threads` | Número de goroutines | `8` |

## 📁 Estrutura do Projeto

```
hashcracker-go/
├── main.go                  # Entry point, flags e banner
├── cracker/
│   ├── identifier.go        # Identificação do tipo de hash
│   ├── hasher.go             # Funções de hashing (MD5, SHA1, SHA256, SHA512)
│   ├── dictionary.go         # Lógica do ataque por dicionário
│   └── bruteforce.go         # Lógica do ataque brute force
├── wordlists/
│   └── common.txt            # Top 10.000 senhas mais comuns
├── go.mod
├── README.md
├── LICENSE
└── .gitignore
```

## 🔧 Como Funciona

### Identificação de Hash
O módulo `identifier.go` analisa o comprimento e formato do hash:
| Tipo | Comprimento | Exemplo |
|------|-------------|---------|
| MD5 | 32 chars | `5f4dcc3b5aa765d61d8327deb882cf99` |
| SHA1 | 40 chars | `5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8` |
| SHA256 | 64 chars | `5e884898da2803...` |
| SHA512 | 128 chars | `b109f3bbbc244...` |
| bcrypt | Prefixo `$2a$` | `$2a$10$...` |

### Worker Pool (Goroutines)
O ataque usa um padrão **produtor-consumidor**:
1. **Produtor**: Lê a wordlist ou gera combinações
2. **Canal**: Distribui as palavras candidatas
3. **Workers (goroutines)**: Calculam hashes e comparam com o alvo

```
[Wordlist] → [Canal] → [Worker 1] → compara hash
                      → [Worker 2] → compara hash
                      → [Worker N] → compara hash
```

## 🧪 Hashes para Teste

```bash
# MD5 de "password"
echo "5f4dcc3b5aa765d61d8327deb882cf99" | ./hashcracker -hash -

# SHA256 de "123456"
./hashcracker -hash "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
```

## ⚠️ Disclaimer

> **Esta ferramenta é destinada apenas para fins educacionais e testes de segurança autorizados.** O uso desta ferramenta para crackear hashes sem permissão é ilegal. Utilize apenas em seus próprios hashes ou em ambientes de teste autorizados. O autor não se responsabiliza pelo uso indevido.

## 📄 Licença

Este projeto está sob a licença [MIT](LICENSE).

---

<p align="center">
  Feito com ❤️ por <a href="https://github.com/kaique">kaique</a>
</p>
