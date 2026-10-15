# HashCracker

**Hash identification and password-candidate testing utility in Go**

[![CI](https://github.com/ispectr3/hashcracker/actions/workflows/ci.yml/badge.svg)](https://github.com/ispectr3/hashcracker/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go: 1.21+](https://img.shields.io/badge/Go-1.21%2B-00ADD8.svg?logo=go)](https://go.dev/)

HashCracker is a command-line tool written in Go. It identifies the format of an input hash and tests password candidates using either a wordlist or a generated brute-force space. The CLI accepts one hash or a file containing multiple hashes and uses configurable goroutines for the attacks.

Use this tool only with hashes and systems you own or are explicitly authorized to assess, such as an approved security assessment, forensic investigation, or CTF/lab. Do not use it against third-party credentials or systems without permission.

## Supported hash formats

| Format | Identification | Candidate testing |
| :--- | :--- | :--- |
| **MD5** | 32 hexadecimal characters | Dictionary and brute force |
| **SHA-1** | 40 hexadecimal characters | Dictionary and brute force |
| **SHA-256** | 64 hexadecimal characters | Dictionary and brute force |
| **SHA-512** | 128 hexadecimal characters | Dictionary and brute force |
| **bcrypt** | `$2a$`, `$2b$`, or `$2y$` format detection | Not implemented |

Identification is based on the input format and length; it does not prove which algorithm produced an ambiguous digest. bcrypt is recognized by format only and is not handled by the candidate-testing functions.

## Requirements and installation

- Go 1.21 or newer

```bash
git clone https://github.com/ispectr3/hashcracker.git
cd hashcracker
go build -o hashcracker .
```

The repository contains `wordlists/common.txt`, which the CLI uses as the default wordlist in dictionary mode when no `-wordlist` is supplied and that file is available from the current directory.

## Usage

Test a single hash against a wordlist:

```bash
./hashcracker -hash 5d41402abc4b2a76b9719d911017c592 -wordlist wordlists/common.txt
```

Process hashes from a file, one per line, with 16 workers:

```bash
./hashcracker -file target_hashes.txt -wordlist /path/to/wordlist.txt -threads 16
```

Use brute force for a short, explicitly bounded candidate space:

```bash
./hashcracker -hash 098f6bcd4621d373cade4e832627b4f6 \
  -mode brute -min-len 1 -max-len 4
```

### Command-line flags

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-hash` | empty | Single hash to process |
| `-file` | empty | File containing hashes; blank lines and lines beginning with `#` are skipped |
| `-wordlist` | empty | Wordlist path for dictionary mode |
| `-mode` | `dict` | `dict` or `brute` |
| `-charset` | `abcdefghijklmnopqrstuvwxyz0123456789` | Characters used by brute force |
| `-min-len` | `1` | Minimum brute-force candidate length |
| `-max-len` | `6` | Maximum brute-force candidate length |
| `-threads` | `8` | Number of worker goroutines |

Run `./hashcracker -h` to see the flags exposed by the CLI.

## Tests and build checks

The repository’s CI workflow runs the following commands on Go 1.21 and 1.22:

```bash
go test -v ./...
go build -v -o hashcracker main.go
```

The checked-in unit tests cover hash identification, including the supported hexadecimal lengths and bcrypt format detection, plus MD5 and SHA-256 generation.

## Repository structure

```text
.
├── main.go                       # CLI, flags, input loading, and output
├── cracker/
│   ├── bruteforce.go             # Generated-candidate attack
│   ├── cracker_test.go            # Unit tests
│   ├── dictionary.go              # Wordlist attack
│   ├── hasher.go                 # MD5, SHA-1, SHA-256, and SHA-512
│   └── identifier.go              # Format-based identification
├── wordlists/common.txt          # Included example wordlist
├── go.mod
├── .github/workflows/ci.yml      # CI test and build workflow
└── LICENSE
```

## Limitations

- Format detection cannot distinguish algorithms that share a digest length and hexadecimal representation.
- bcrypt is detected but cannot be cracked by the current implementation.
- Dictionary mode only tests candidates present in the supplied wordlist.
- Brute force is bounded by `-charset`, `-min-len`, and `-max-len`; the search space can grow rapidly.
- There is no claim of guaranteed recovery: a password is reported only when a tested candidate produces the target digest.

## License

HashCracker is open-source software licensed under the [MIT License](LICENSE).
