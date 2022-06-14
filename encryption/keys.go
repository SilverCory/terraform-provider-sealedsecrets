package encryption

import (
	"bytes"
	"encoding/base64"
	"io"

	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/openpgp"
)

var encoder = base64.StdEncoding

func base64Decode(in string) (io.Reader, error) {
	decKey, err := encoder.DecodeString(in)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64 key: %w", err)
	}

	return bytes.NewBuffer(decKey), nil
}

func decodeKey(key string) (openpgp.EntityList, error) {
	keyBuf, err := base64Decode(key)
	if err != nil {
		return nil, err
	}

	el, err := openpgp.ReadArmoredKeyRing(keyBuf)
	if err != nil {
		return nil, fmt.Errorf("unable to read key ring: %w", err)
	}

	return el, nil
}

func hashString(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
