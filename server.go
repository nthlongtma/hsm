package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"hsm/models"

	"github.com/gemalto/pkcs11"
)

type (
	Server struct {
		conf *models.Config
		hsm  HSM
	}
)

func NewServer(conf *models.Config, hsm HSM) Server {
	return Server{conf: conf, hsm: hsm}
}

func (s Server) HandleEncrypt(w http.ResponseWriter, r *http.Request) {
	req := models.EncryptRequest{}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Write(ToJson(models.EncryptResponse{
			ErrorCode:    "1111",
			ErrorMessage: err.Error(),
		}))
		return
	}
	log.Printf("encrypt request: %+v", req)

	// decode the plain text
	plainText, err := base64.StdEncoding.DecodeString(req.PlainText)
	if err != nil {
		w.Write(ToJson(models.EncryptResponse{
			ErrorCode:    "2222",
			ErrorMessage: err.Error(),
		}))
		return
	}
	log.Printf("plain text after decode: %s", string(plainText))

	// encrypt
	cipher, err := s.encrypt(pkcs11.CKO_SECRET_KEY, s.conf.HSM.N2kLabel, plainText)
	if err != nil {
		w.Write(ToJson(models.EncryptResponse{
			ErrorCode:    "3333",
			ErrorMessage: err.Error(),
		}))
		return
	}

	// encode cipher before response back
	eCipher := base64.StdEncoding.EncodeToString(cipher)

	w.Write(ToJson(models.EncryptResponse{
		ErrorCode:    "0000",
		ErrorMessage: "success",
		CipherText:   eCipher,
	}))
	return
}

func (s Server) HandleDecrypt(w http.ResponseWriter, r *http.Request) {
	req := models.DecryptRequest{}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Write(ToJson(models.DecryptResponse{
			ErrorCode:    "1111",
			ErrorMessage: err.Error(),
		}))
		return
	}
	log.Printf("request: %+v", req)

	// decode cipher text
	cipher, err := base64.StdEncoding.DecodeString(req.CipherText)
	if err != nil {
		w.Write(ToJson(models.DecryptResponse{
			ErrorCode:    "2222",
			ErrorMessage: err.Error(),
		}))
		return
	}

	// decrypt
	plainText, err := s.decrypt(pkcs11.CKO_SECRET_KEY, s.conf.HSM.N2kLabel, cipher)
	if err != nil {
		w.Write(ToJson(models.DecryptResponse{
			ErrorCode:    "3333",
			ErrorMessage: err.Error(),
		}))
		return
	}

	// encode plain text before response
	ePlainText := base64.StdEncoding.EncodeToString(plainText)

	w.Write(ToJson(models.DecryptResponse{
		ErrorCode:    "0000",
		ErrorMessage: "success",
		PlainText:    ePlainText,
	}))

	return
}

// encrypt with pre-generated iv then append the iv to the output.
func (s Server) encrypt(keyClass uint, keyLabel string, plainText []byte) ([]byte, error) {
	obj, err := FindKeys(s.hsm.ctx, s.hsm.ss, keyClass, keyLabel)
	if err != nil {
		return nil, fmt.Errorf("failed to find key: %v", err)
	}

	iv := genIV(s.conf.HSM.IVSize)
	cipher, err := Encrypt(s.hsm.ctx, s.hsm.ss, obj[0], pkcs11.CKM_AES_CBC_PAD, plainText, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %v", err)
	}

	// prepend iv to cipher
	ret := make([]byte, len(iv)+len(cipher))
	copy(ret[:len(iv)], iv)
	copy(ret[len(iv):], cipher)

	return ret, nil
}

// extract the iv from cipher and decrypt.
func (s Server) decrypt(keyClass uint, keyLabel string, cipher []byte) ([]byte, error) {
	obj, err := FindKeys(s.hsm.ctx, s.hsm.ss, keyClass, keyLabel)
	if err != nil {
		return nil, err
	}

	// extract iv and cipher
	iv := cipher[:s.conf.HSM.IVSize]
	c := cipher[s.conf.HSM.IVSize:]
	plain, err := Decrypt(s.hsm.ctx, s.hsm.ss, obj[0], pkcs11.CKM_AES_CBC_PAD, c, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plain, nil
}
