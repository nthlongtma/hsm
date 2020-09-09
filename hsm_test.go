package main

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/gemalto/pkcs11"
)

const (
	modulePath = "./module/libsofthsm2.so"
	slotID     = 866444829
	pin        = "12345678"
	label      = "master-key"

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

func TestCreateKey(t *testing.T) {
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

	_, err = CreateKey(ctx, ss, label)
	if err != nil {
		t.Error(err)
	}
	t.Log("create key successfully")
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

	_, err = CreateKey(ctx, ss, label)
	if err != nil {
		t.Error(err)
	}

	objs, err := FindKeys(ctx, ss, pkcs11.CKO_SECRET_KEY, label)
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

	obj, err := CreateKey(ctx, ss, label)
	if err != nil {
		t.Error(err)
	}

	if err := RemoveKey(ctx, ss, obj); err != nil {
		t.Error(err)
	}
	t.Log("remove key successfully")
}

func TestEncrypt(t *testing.T) {
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

	obj, err := CreateKey(ctx, ss, label)
	if err != nil {
		t.Error(err)
	}

	t.Run("Encrypt", func(t *testing.T) {
		// t.Run("ECB", func(t *testing.T) {
		// 	cipher, err := Encryption(ctx, ss, obj, pkcs11.CKM_AES_ECB, []byte(plainText), []byte(iv))
		// 	if err != nil {
		// 		t.Error(err)
		// 	}
		// 	t.Logf("ECB result: %s", cipher)
		// })
		// t.Run("CBC", func(t *testing.T) {
		// 	cipher, err := Encryption(ctx, ss, obj, pkcs11.CKM_AES_CBC, []byte(plainText), []byte(iv))
		// 	if err != nil {
		// 		t.Error(err)
		// 	}
		// 	t.Logf("CBC result: %s", cipher)
		// })
		t.Run("CBC-PAD", func(t *testing.T) { // auto pading
			cipher, err := Encryption(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, []byte(plainText), []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("CBC-PAD result: %s", base64.StdEncoding.EncodeToString(cipher))
		})
	})
}

func TestDecrypt(t *testing.T) {
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

	obj, err := CreateKey(ctx, ss, label)
	if err != nil {
		t.Error(err)
	}

	t.Run("Encrypt", func(t *testing.T) {
		// t.Run("ECB", func(t *testing.T) {
		// 	cipher, err := Encryption(ctx, ss, obj, pkcs11.CKM_AES_ECB, []byte(plainText), []byte(iv))
		// 	if err != nil {
		// 		t.Error(err)
		// 	}
		// 	t.Logf("ECB result: %s", cipher)
		// })
		// t.Run("CBC", func(t *testing.T) {
		// 	cipher, err := Encryption(ctx, ss, obj, pkcs11.CKM_AES_CBC, []byte(plainText), []byte(iv))
		// 	if err != nil {
		// 		t.Error(err)
		// 	}
		// 	t.Logf("CBC result: %s", cipher)
		// })
		t.Run("CBC-PAD", func(t *testing.T) { // auto pading
			cipher, err := Encryption(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, []byte(plainText), []byte(iv))
			if err != nil {
				t.Error(err)
			}
			t.Logf("CBC-PAD result: %s", base64.StdEncoding.EncodeToString(cipher))

			decrypted, err := Decryption(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, cipher, []byte(iv))
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
