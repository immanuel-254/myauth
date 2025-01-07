package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
)

func GenerateAESKey() []byte {
	key := make([]byte, 32) // AES-256 requires a 32-byte key
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}
	return key
}

func GenerateRSAKeys() (*rsa.PrivateKey, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SaveKeyToPEM(key interface{}, filedir, filePath string) error {
	// create the directorate
	err := os.MkdirAll(filedir, os.ModePerm)
	if err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create PEM block
	var pemBlock *pem.Block
	switch key := key.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		}
	case *rsa.PublicKey:
		pemBytes, err := x509.MarshalPKIXPublicKey(key)
		if err != nil {
			return err
		}
		pemBlock = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pemBytes,
		}
	default:
		return fmt.Errorf("unsupported key type")
	}

	// Write PEM block to file
	err = pem.Encode(file, pemBlock)
	if err != nil {
		return err
	}
	return nil
}

func CreateRSAKeyFile(path, file string) error {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		privateKey, err := GenerateRSAKeys()
		if err != nil {
			return err
		}

		// Save private key to .pem file
		err = SaveKeyToPEM(privateKey, path, fmt.Sprintf("%s/%s", path, file))
		if err != nil {
			fmt.Println("Error saving private key:", err)
			return err
		}

	}
	return nil

}

func CreateAESKeyFile() {
	aesFilePath := "aes_key.pem"

	_, err := os.Stat(aesFilePath)

	if os.IsNotExist(err) {
		aesKey := GenerateAESKey()

		// Encode key to PEM format
		block := &pem.Block{
			Type:  "AES KEY",
			Bytes: aesKey,
		}

		// Write PEM encoded key to a file
		file, err := os.Create(aesFilePath)
		if err != nil {
			fmt.Println("Error generating AES key file:", err)
			return
		}
		defer file.Close()

		err = pem.Encode(file, block)
		if err != nil {
			fmt.Println("Error Writing AES key to file:", err)
			return
		}

	} else {
		return
	}
}

func LoadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	// Read the PEM file
	pemData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Decode PEM data
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, err
	}

	// Parse the RSA private key
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func LoadPublicKey(pemString string) (*rsa.PublicKey, error) {
	// Decode the PEM block
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// Parse the public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Type assert to *rsa.PublicKey
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	return rsaPublicKey, nil
}

func LoadAESKey(filename string) ([]byte, error) {
	// Read the PEM file
	pemData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// Check if the PEM block type is "AES KEY"
	if block.Type != "AES KEY" {
		return nil, errors.New("unexpected PEM block type")
	}

	// Extract the AES key
	key := block.Bytes

	return key, nil
}

func PublicKeyPem(pubkey *rsa.PublicKey) (string, error) {
	// Convert the public key to DER format
	derBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}

	// Create a PEM block for the public key
	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derBytes,
	}

	// Encode the PEM block to a string
	pemString := string(pem.EncodeToMemory(pemBlock))

	return pemString, nil
}

func EncryptGCMAES(key []byte, plaintext string) ([]byte, []byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	return ciphertext, nonce, nil
}

func DecryptGCMAES(key []byte, nonce []byte, ct []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}

	pt := make([]byte, len(ct))
	c.Decrypt(pt, ct)

	return plaintext, nil
}

func HybridEncrypt(data, pubKey string) (string, string, string, error) {
	key := GenerateAESKey()
	encdata, nonce, err := EncryptGCMAES(key, data)

	if err != nil {
		return "", "", "", err
	}

	pubkey, err := LoadPublicKey(pubKey)
	if err != nil {
		return "", "", "", err
	}

	// Encrypt the message using RSA-OAEP
	enckey, err := rsa.EncryptOAEP(
		sha256.New(), // Random source
		rand.Reader,
		pubkey, // Public key
		key,    // Message to encrypt
		nil,    // Label (use nil for no label)
	)

	if err != nil {
		return "", "", "", err
	}

	return hex.EncodeToString(encdata), hex.EncodeToString(enckey), hex.EncodeToString(nonce), nil
}

func HybridDecrypt(data, nonce, key string) (string, error) {
	encdata, err := hex.DecodeString(data)

	if err != nil {
		return "", err
	}

	encnonce, err := hex.DecodeString(nonce)

	if err != nil {
		return "", err
	}

	enckey, err := hex.DecodeString(key)

	if err != nil {
		return "", err
	}

	privkey, err := LoadPrivateKey("private_key.pem")
	if err != nil {
		return "", err
	}

	deckey, err := rsa.DecryptOAEP(
		sha256.New(), // Random source
		rand.Reader,
		privkey,
		enckey,
		nil,
	)

	if err != nil {
		return "", err
	}

	decdata, err := DecryptGCMAES(deckey, encnonce, encdata)

	if err != nil {
		return "", err
	}

	return string(decdata), nil
}

func EncryptAESFromPEM(plaintext string) ([]byte, []byte, error) {
	// Load AESKey
	aesKey, err := LoadAESKey("aes_key.pem")
	if err != nil {
		return nil, nil, err
	}

	// Encrypt message
	aesencData, nonce, err := EncryptGCMAES(aesKey, plaintext)
	if err != nil {
		return nil, nil, err
	}

	return aesencData, nonce, err

}

func DecryptAESFromPEM(ct, nonce []byte) ([]byte, error) {
	// Load AESKey
	aesKey, err := LoadAESKey("aes_key.pem")
	if err != nil {
		return nil, err
	}

	// Encrypt message
	aesdecData, err := DecryptGCMAES(aesKey, nonce, ct)
	if err != nil {
		return nil, err
	}

	return aesdecData, err
}
