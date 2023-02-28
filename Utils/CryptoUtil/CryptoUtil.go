package CryptoUtil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type Crypto struct {
	key      string
	FilePath string
	Text     string
}

func NewKey() *Crypto {
	c := new(Crypto)
	c.key = "AnKoloft@~delNazok!12345" // key parameter must be 16, 24 or 32,
	return c
}

func (k *Crypto) Sha256() string {
	h := sha256.New()
	h.Write([]byte(k.Text))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (k *Crypto) Md5Sum() (string, error) {
	file, err := os.Open(k.FilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (k *Crypto) Encrypt() (string, error) {
	c, err := aes.NewCipher([]byte(k.key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err

	}
	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(k.Text), nil)), nil
}

func (k *Crypto) Decrypt() (string, error) {

	ciphertext, _ := base64.StdEncoding.DecodeString(k.Text)
	c, err := aes.NewCipher([]byte(k.key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	t, e := gcm.Open(nil, nonce, ciphertext, nil)
	if e != nil {
		return "", err
	}

	return bytes.NewBuffer(t).String(), nil

}
