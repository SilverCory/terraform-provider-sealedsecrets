package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-provider-scaffolding/encryption"
)

var publicKey string
var input string

func main() {
	flag.StringVar(&publicKey, "public-key", os.Getenv("SECRET_PUBLIC_KEY"), "The public key to encrypt against")
	flag.StringVar(&input, "input", "", "The secret to encrypt")
	flag.Parse()

	if publicKey == "" {
		_, _ = fmt.Fprintf(os.Stderr, "The public key was not supplied, set the environment variable SECRET_PUBLIC_KEY or use the -public-key flag\n")
		os.Exit(1)
		return
	}

	enc, err := encryption.New(publicKey, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred setting up the encryptor: %v\n", err)
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

	// fmt.Printf("%q\n", input)

	if input == "" {
		_, _ = fmt.Fprintf(os.Stderr, "[WARNING] The input going to be encrypted is an empty string.\n", err)
	}

	out, err := enc.EncryptString(input)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occurred encrypting secret: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println(out)
}
