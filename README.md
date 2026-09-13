# HashCracker

**Multi-Hash Algorithm Identification & Verification Engine in Go**

[![CI](https://github.com/ispectr3/hashcracker/actions/workflows/ci.yml/badge.svg)](https://github.com/ispectr3/hashcracker/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go: 1.21+](https://img.shields.io/badge/Go-1.21%2B-00ADD8.svg?logo=go)](https://golang.org)

HashCracker is a lightweight cryptanalysis utility developed in Go for offline credential assessment, forensics, and CTF challenges. It provides automated hash type identification based on bit-length and structural signatures, combined with multi-threaded dictionary matching and brute-force verification across standard cryptographic digest algorithms.

---

## Supported Hash Types

| Algorithm | Digest Length | Bits | Format Signature |
| :--- | :--- | :--- | :--- |
| **MD5** | 32 hex chars | 128-bit | `^[a-fA-F0-9]{32}$` |
| **SHA-1** | 40 hex chars | 160-bit | `^[a-fA-F0-9]{40}$` |
| **SHA-256** | 64 hex chars | 256-bit | `^[a-fA-F0-9]{64}$` |
| **SHA-512** | 128 hex chars | 512-bit | `^[a-fA-F0-9]{128}$` |
| **bcrypt** | 60 chars | Variable | `^\$2[aby]\$\d{2}\$.{53}$` (Detection) |

---

## Installation

### Prerequisites

- Go 1.21+

```bash
git clone https://github.com/ispectr3/hashcracker.git
cd hashcracker
go build -o hashcracker main.go
```

---

## Usage

```bash
# Identify and crack a single hash via dictionary
./hashcracker -hash 5d41402abc4b2a76b9719d911017c592 -wordlist wordlists/common.txt

# Crack multiple hashes from a file with 16 worker goroutines
./hashcracker -file target_hashes.txt -wordlist /usr/share/wordlists/rockyou.txt -threads 16

# Brute force attack on short passwords
./hashcracker -hash 098f6bcd4621d373cade4e832627b4f6 -mode brute -min-len 1 -max-len 4
```

---

## Automated Testing

```bash
go test -v ./...
```

---

## License

HashCracker is open-source software licensed under the [MIT License](LICENSE).
