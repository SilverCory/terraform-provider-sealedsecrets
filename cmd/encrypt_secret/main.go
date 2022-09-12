package main

import (
	"flag"
	"fmt"
	"github.com/silvercory/terraform-provider-sealedsecrets/encryption"
	"io"
	"os"
)

var privateKey string
var publicKey string
var input string
var decrypt bool

func main() {
	flag.StringVar(&publicKey, "public-key", os.Getenv("SECRET_PUBLIC_KEY"), "The public key to encrypt against")
	flag.StringVar(&input, "input", "", "The secret to encrypt")

	flag.StringVar(&privateKey, "private-key", os.Getenv("SECRET_PRIVATE_KEY"), "The private key to decrypt with")
	flag.BoolVar(&decrypt, "decrypt", false, "Decrypt")

	flag.Parse()

	if !decrypt && publicKey == "" {
		privateKey = ""
		_, _ = fmt.Fprintf(os.Stderr, "The public key was not supplied, set the environment variable SECRET_PUBLIC_KEY or use the -public-key flag\n")
		os.Exit(1)
		return
	}

	if decrypt && privateKey == "" {
		publicKey = ""
		_, _ = fmt.Fprintf(os.Stderr, "The private key was not supplied, set the environment variable SECRET_PRIVATE_KEY or use the -private-key flag\n")
		os.Exit(1)
		return
	}

	enc, err := encryption.New(publicKey, privateKey)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occurred setting up the encryptor: %v\n", err)
		os.Exit(1)
		return
	}

	if input == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error occurred reading stdin (maybe set -input?): %v\n", err)
			os.Exit(1)
			return
		}

		input = string(data)
	}

	if input == "" {
		_, _ = fmt.Fprintf(os.Stderr, "[WARNING] The input going to be encrypted is an empty string.\n", err)
	}

	encFunc := enc.EncryptString
	if decrypt {
		encFunc = func(in string) (string, error) {
			out, _, err := enc.DecryptString(in)
			return out, err
		}
	}

	out, err := encFunc(input)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occurred encrypting secret: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println(out)
}
