package encryption

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
)

var ErrNoKeys = errors.New("no private or public key has been supplied")

type Encryptor struct {
	publicKey  openpgp.EntityList
	privateKey openpgp.EntityList
}

func New(publicKey string, privateKey string) (*Encryptor, error) {

	var encryptor = &Encryptor{}
	var err error

	if publicKey == "" && privateKey == "" {
		return nil, ErrNoKeys
	}

	if publicKey != "" {
		encryptor.publicKey, err = decodeKey(publicKey)
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt public key: %w", err)
		}
	}

	if privateKey != "" {
		encryptor.privateKey, err = decodeKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt private key: %w", err)
		}
	}

	return encryptor, nil
}

func (e *Encryptor) DecryptString(in string) (output string, hash string, err error) {
	inBuf, err := base64Decode(in)
	if err != nil {
		return "", "", err
	}

	md, err := openpgp.ReadMessage(inBuf, e.privateKey, nil, nil)
	if err != nil {
		return "", "", err
	}

	data, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", "", err
	}

	output = string(data)

	return output, hashString(output), nil
}

func (e *Encryptor) EncryptString(in string) (string, error) {
	var buf = new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, e.publicKey, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("creating encryption writer: %w", err)
	}

	_, err = w.Write([]byte(in))
	if err != nil {
		return "", fmt.Errorf("writing data to be encrypted: %w", err)
	}

	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("closing encryption writer: %w", err)
	}

	return encoder.EncodeToString(buf.Bytes()), nil
}
