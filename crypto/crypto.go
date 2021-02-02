package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
)

const RsaLen = 2048

func GetKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	keyFile := "cert/key.pem"
	mysalt := "Serpico78"
	priv, err := privateKeyFromFile(keyFile, mysalt)
	if err != nil {
		log.Println("Key pem file error, create a new one", err)
		priv, _ = rsa.GenerateKey(rand.Reader, RsaLen)
		err = savePrivateKeyInFile(keyFile, priv, mysalt)
		if err != nil {
			return nil, nil, err
		}
	}

	return priv, &priv.PublicKey, err
}

func Encrypt(plain []byte, pubkey *rsa.PublicKey) []byte {

	//è interessante notare la procedura ibrida della criptazione.
	// Viene generata una nuova chiave random la quale viene poi criptata con la chiave pubblica
	// e messa in testa al file. La chiave della sessione viene criptata con rsa.
	// Mentre il file viene creiptato con aes che è una procedura di cifrazione simmetrica.
	key := make([]byte, 256/8) // AES-256
	io.ReadFull(rand.Reader, key)

	encKey, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, key, nil)
	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, aesgcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciph := aesgcm.Seal(nil, nonce, plain, nil)
	s := [][]byte{encKey, nonce, ciph}
	return bytes.Join(s, []byte{})
}

func Decrypt(ciph []byte, priv *rsa.PrivateKey) ([]byte, error) {
	//Per primo viene estratta la chiave per la decriptazione via aes.
	// La chiave è in testa al file ed è codificata in rsa. La decriptazione della chiave per
	// la sessione aes è possibile solo via rsa utilizzando la chiave privata in formato pem.
	encKey := ciph[:RsaLen/8]
	ciph = ciph[RsaLen/8:]
	key, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, nil)

	block, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(block)
	nonce := ciph[:aesgcm.NonceSize()]
	ciph = ciph[aesgcm.NonceSize():]

	return aesgcm.Open(nil, nonce, ciph, nil)
}

func savePrivateKeyInFile(file string, priv *rsa.PrivateKey, pwd string) error {
	der := x509.MarshalPKCS1PrivateKey(priv)
	pp := []byte(pwd)
	block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", der, pp, x509.PEMCipherAES256)
	if err != nil {
		return err
	}
	log.Println("Save the key in ", file)
	return ioutil.WriteFile(file, pem.EncodeToMemory(block), 0644)
}

func privateKeyFromFile(file string, pwd string) (*rsa.PrivateKey, error) {
	der, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(der)

	der, err = x509.DecryptPEMBlock(block, []byte(pwd))
	if err != nil {
		return nil, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(der)
	return priv, nil
}
