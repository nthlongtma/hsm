package main

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/gemalto/pkcs11"
)

const (
	modulePath  = "./module/libsofthsm2.so"
	slotID      = 866444829
	pin         = "12345678"
	keyAESLabel = "aes-key"
	keyRSALabel = "rsa-key"

	plainText = "kbtg-tma team building"
	iv        = "0123456789abcdef"
)

func TestGetContext(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	FinishContext(ctx)
}

func TestGetSession(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)
	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)
}

func TestGetSlotList(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	FinishContext(ctx)

	sl, err := ctx.GetSlotList(true)
	if err != nil {
		t.Error(err)
	}
	t.Logf("slot list: %+v", sl)
}

func TestGetInfo(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	info, err := ctx.GetInfo()
	if err != nil {
		t.Error(err)
	}
	t.Logf("HSM information: %+v", info)
}

func TestCreateAESKey(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	k, err := CreateAESKey(ctx, ss, keyAESLabel)
	if err != nil {
		t.Error(err)
	}
	t.Logf("create AES key successfully: %v", k)
}

func TestFindKeys(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	_, err = CreateAESKey(ctx, ss, keyAESLabel)
	if err != nil {
		t.Error(err)
	}

	objs, err := FindKeys(ctx, ss, pkcs11.CKO_SECRET_KEY, keyAESLabel)
	if err != nil {
		t.Error(err)
	}
	t.Logf("found %d key", len(objs))
}

func TestRemoveKey(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	obj, err := CreateAESKey(ctx, ss, keyAESLabel)
	if err != nil {
		t.Error(err)
	}

	if err := RemoveKey(ctx, ss, obj); err != nil {
		t.Error(err)
	}
	t.Log("remove key successfully")
}

func TestSymmetricEncrypt(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	obj, err := CreateAESKey(ctx, ss, keyAESLabel)
	if err != nil {
		t.Error(err)
	}

	t.Run("Encrypt", func(t *testing.T) {
		t.Run("ECB", func(t *testing.T) {
			padPlainText, _ := pkcs7Pad([]byte(plainText), 32)
			cipher, err := Encrypt(ctx, ss, obj, pkcs11.CKM_AES_ECB, padPlainText, []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("ECB result: %s", base64.StdEncoding.EncodeToString(cipher))
		})
		t.Run("CBC", func(t *testing.T) {
			padPlainText, _ := pkcs7Pad([]byte(plainText), 32)
			cipher, err := Encrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC, padPlainText, []byte(iv))
			if err != nil {
				t.Error(err)
			}

			t.Logf("CBC result: %s", base64.StdEncoding.EncodeToString(cipher))
		})
		t.Run("CBC-PAD", func(t *testing.T) { // auto pading
			cipher, err := Encrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, []byte(plainText), []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("CBC-PAD result: %s", base64.StdEncoding.EncodeToString(cipher))
		})
	})
}

func TestSymmetricDecrypt(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	obj, err := CreateAESKey(ctx, ss, keyAESLabel)
	if err != nil {
		t.Error(err)
	}

	t.Run("Decrypt", func(t *testing.T) {
		t.Run("ECB", func(t *testing.T) {
			padded, _ := pkcs7Pad([]byte(plainText), 32)
			cipher, err := Encrypt(ctx, ss, obj, pkcs11.CKM_AES_ECB, padded, []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("ECB result: %s", base64.StdEncoding.EncodeToString(cipher))

			decrypted, err := Decrypt(ctx, ss, obj, pkcs11.CKM_AES_ECB, cipher, []byte(iv))
			if err != nil {
				t.Error(err)
			}
			unPadded, _ := pkcs7Unpad(decrypted, 32)
			t.Logf("decrypted: %s", unPadded)

			if bytes.Compare([]byte(plainText), unPadded) != 0 {
				t.Error("missmatch")
			}
		})
		t.Run("CBC", func(t *testing.T) {
			padded, _ := pkcs7Pad([]byte(plainText), 32)
			cipher, err := Encrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC, padded, []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("CBC result: %s", base64.StdEncoding.EncodeToString(cipher))

			decrypted, err := Decrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC, cipher, []byte(iv))
			if err != nil {
				t.Error(err)
			}
			unPadded, _ := pkcs7Unpad(decrypted, 32)
			t.Logf("decrypted: %s", unPadded)

			if bytes.Compare([]byte(plainText), unPadded) != 0 {
				t.Error("missmatch")
			}
		})
		t.Run("CBC-PAD", func(t *testing.T) { // auto pading
			cipher, err := Encrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, []byte(plainText), []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("CBC-PAD result: %s", base64.StdEncoding.EncodeToString(cipher))

			decrypted, err := Decrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, cipher, []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("decrypted: %s", decrypted)

			if bytes.Compare([]byte(plainText), decrypted) != 0 {
				t.Error("missmatch")
			}
		})
	})
}

func TestCreateRSAKeyPair(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	pbk, pvk, err := CreateRSAKeyPair(ctx, ss, keyRSALabel)
	if err != nil {
		t.Error(err)
	}
	t.Logf("create RSA key pair successfully: %v - %v", pbk, pvk)
}

func TestSignAndVerify(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, slotID, pin)
	if err != nil {
		t.Error(err)
	}
	defer FinishSession(ctx, ss)

	pbk, pvk, err := CreateRSAKeyPair(ctx, ss, keyRSALabel)
	if err != nil {
		t.Error(err)
	}

	t.Run("Sign-Verify", func(t *testing.T) {
		if err := ctx.SignInit(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA1_RSA_PKCS, nil)}, pvk); err != nil {
			t.Error(err)
		}

		sign, err := ctx.Sign(ss, []byte(plainText))
		if err != nil {
			t.Error(err)
		}
		t.Logf("sign successfully, signature: %s", base64.StdEncoding.EncodeToString(sign))

		if err := ctx.VerifyInit(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA1_RSA_PKCS, nil)}, pbk); err != nil {
			t.Error(err)
		}

		if err := ctx.Verify(ss, []byte(plainText), sign); err != nil {
			t.Error(err)
		}
		t.Log("verify successfully")
	})
}
