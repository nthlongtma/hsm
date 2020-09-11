package main

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/gemalto/pkcs11"
)

const (
	modulePath     = "./module/libsofthsm2.so"
	slotID         = 866444829
	pin            = "12345678"
	labelSecretKey = "secret"
	labelRSAKey    = "rsa"

	plainText = "kbtg-tma team building"
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
	defer FinishContext(ctx)

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

	t.Run("Get HSM Information", func(t *testing.T) {
		info, err := ctx.GetInfo()
		if err != nil {
			t.Error(err)
		}
		t.Logf("HSM information: %s", ToJsonString(info))
	})
}

func TestGetMech(t *testing.T) {
	ctx, err := GetContext(modulePath)
	if err != nil {
		t.Error(err)
	}
	defer FinishContext(ctx)

	mechs, err := ctx.GetMechanismList(slotID)
	if err != nil {
		t.Error(err)
	}
	t.Logf("mechs: %+v", mechs)

	info, err := ctx.GetMechanismInfo(slotID, mechs)
	if err != nil {
		t.Error(err)
	}
	t.Logf("info: %+v", info)
}

// secret key
func TestCreateSecretKey(t *testing.T) {
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

	t.Run("Create-Secret-Key", func(t *testing.T) {
		k, err := CreateSecretKey(ctx, ss, labelSecretKey)
		if err != nil {
			t.Error(err)
		}
		t.Logf("create secret key successfully: %v", k)
	})
}

// find secret key that created from TestCreateSecretKey
func TestFindSecretKeys(t *testing.T) {
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

	t.Run("Find-Secret-Key", func(t *testing.T) {
		objs, err := FindKeys(ctx, ss, pkcs11.CKO_SECRET_KEY, labelSecretKey)
		if err != nil {
			t.Log("not found secret key")
		} else {
			t.Logf("found secret key: %v", objs[0])
		}
	})
}

// remove secret key that created from TestCreateSecretKey
func TestRemoveSecretKey(t *testing.T) {
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

	t.Run("Remove-Secret-Key", func(t *testing.T) {
		objs, err := FindKeys(ctx, ss, pkcs11.CKO_SECRET_KEY, labelSecretKey)
		if err != nil {
			t.Logf("not found secret key")
		} else {
			if err := RemoveKey(ctx, ss, objs[0]); err != nil {
				t.Error(err)
			}
			t.Logf("remove secret key successfully: %v", objs[0])
		}
	})
}

// use secret key for both encryption and decryption
func TestSymmetric(t *testing.T) {
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

	obj, err := CreateSecretKey(ctx, ss, "test-decrypt-symmetric")
	if err != nil {
		t.Error(err)
	}
	defer RemoveKey(ctx, ss, obj)

	padded, _ := pkcs7Pad([]byte(plainText), 32)

	t.Run("Symmetric", func(t *testing.T) {
		t.Run("ECB", func(t *testing.T) {
			cipher, decrypted := []byte{}, []byte{}
			iv := genIV(16)

			t.Run("Encrypt", func(t *testing.T) {
				cipher, err = Encrypt(ctx, ss, obj, pkcs11.CKM_AES_ECB, padded, iv)
				if err != nil {
					t.Error(err)
				}
				t.Logf("cipher: %s", base64.StdEncoding.EncodeToString(cipher))
			})

			t.Run("Decrypt", func(t *testing.T) {
				decrypted, err = Decrypt(ctx, ss, obj, pkcs11.CKM_AES_ECB, cipher, iv)
				if err != nil {
					t.Error(err)
				}

				unPadded, _ := pkcs7Unpad(decrypted, 32)
				t.Logf("decrypted: %s", unPadded)

				if bytes.Compare([]byte(plainText), unPadded) != 0 {
					t.Error("missmatch")
				}
			})
		})

		t.Run("CBC", func(t *testing.T) {
			cipher, decrypted := []byte{}, []byte{}
			iv := genIV(16)

			t.Run("Encrypt", func(t *testing.T) {
				cipher, err = Encrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC, padded, iv)
				if err != nil {
					t.Error(err)
				}
				t.Logf("cipher: %s", base64.StdEncoding.EncodeToString(cipher))
			})

			t.Run("Decrypt", func(t *testing.T) {
				decrypted, err = Decrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC, cipher, iv)
				if err != nil {
					t.Error(err)
				}

				unPadded, _ := pkcs7Unpad(decrypted, 32)
				t.Logf("decrypted: %s", unPadded)

				if bytes.Compare([]byte(plainText), unPadded) != 0 {
					t.Error("missmatch")
				}
			})
		})

		t.Run("CBC-PAD", func(t *testing.T) { // auto pading
			t.Run("CBC-PAD", func(t *testing.T) {
				cipher, decrypted := []byte{}, []byte{}
				iv := genIV(16)

				t.Run("Encrypt", func(t *testing.T) {
					cipher, err = Encrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, []byte(plainText), iv)
					if err != nil {
						t.Error(err)
					}
					t.Logf("cipher: %s", base64.StdEncoding.EncodeToString(cipher))
				})

				t.Run("Decrypt", func(t *testing.T) {
					decrypted, err = Decrypt(ctx, ss, obj, pkcs11.CKM_AES_CBC_PAD, cipher, iv)
					if err != nil {
						t.Error(err)
					}
					t.Logf("decrypted: %s", decrypted)

					if bytes.Compare([]byte(plainText), decrypted) != 0 {
						t.Error("missmatch")
					}
				})
			})

		})
	})
}

// public key and private key
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

	t.Run("Create-RSA-Keys", func(t *testing.T) {
		pbk, pvk, err := CreateRSAKeyPair(ctx, ss, labelRSAKey)
		if err != nil {
			t.Error(err)
		}
		t.Logf("create RSA key pair successfully {public - private}: %v - %v", pbk, pvk)
	})
}

// find keys from TestCreateRSAKeyPair
func TestFindRSAKeyPair(t *testing.T) {
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

	t.Run("Find-RSA-Keys", func(t *testing.T) {
		t.Run("Find-Public-Key", func(t *testing.T) {
			objs, err := FindKeys(ctx, ss, pkcs11.CKO_PUBLIC_KEY, labelRSAKey)
			if err != nil {
				t.Log("not found public key")
			} else {
				t.Logf("found public key: %v ", objs[0])
			}
		})
		t.Run("Find-Private-Key", func(t *testing.T) {
			objs, err := FindKeys(ctx, ss, pkcs11.CKO_PRIVATE_KEY, labelRSAKey)
			if err != nil {
				t.Log("not found private key")
			} else {
				t.Logf("found private key: %v", objs[0])
			}
		})
	})

}

// remove keys from TestCreateRSAKeyPair
func TestRemoveRSAKeyPair(t *testing.T) {
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

	t.Run("Remove-RSA-Keys", func(t *testing.T) {
		t.Run("Remove-Public-Key", func(t *testing.T) {
			objs, err := FindKeys(ctx, ss, pkcs11.CKO_PUBLIC_KEY, labelRSAKey)
			if err != nil {
				t.Log("not found public key")
			} else {
				if err := RemoveKey(ctx, ss, objs[0]); err != nil {
					t.Error(err)
				}
				t.Logf("remove public key successfully: %v", objs[0])
			}
		})
		t.Run("Remove-Private-Key", func(t *testing.T) {
			objs, err := FindKeys(ctx, ss, pkcs11.CKO_PRIVATE_KEY, labelRSAKey)
			if err != nil {
				t.Log("not found private key")
			} else {
				if err := RemoveKey(ctx, ss, objs[0]); err != nil {
					t.Error(err)
				}
				t.Logf("remove private key successfully: %v", objs[0])
			}
		})
	})

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

	pbk, pvk, err := CreateRSAKeyPair(ctx, ss, "test-sign-verify")
	if err != nil {
		t.Error(err)
	}
	defer RemoveKey(ctx, ss, pbk)
	defer RemoveKey(ctx, ss, pvk)

	t.Run("Sign-Verify", func(t *testing.T) {
		signature := []byte{}
		t.Run("Sign", func(t *testing.T) { //sign with private key
			if err := ctx.SignInit(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA1_RSA_PKCS, nil)}, pvk); err != nil {
				t.Error(err)
			}

			signature, err = ctx.Sign(ss, []byte(plainText))
			if err != nil {
				t.Error(err)
			}
			t.Logf("signature: %s", base64.StdEncoding.EncodeToString(signature))
		})

		t.Run("Verify", func(t *testing.T) { // verify with public key
			if err := ctx.VerifyInit(ss, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA1_RSA_PKCS, nil)}, pbk); err != nil {
				t.Error(err)
			}

			if err := ctx.Verify(ss, []byte(plainText), signature); err != nil {
				t.Error(err)
			}
			t.Log("verify successfully!")
		})

	})
}

func TestAsymmetric(t *testing.T) {
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

	pbk, pvk, err := CreateRSAKeyPair(ctx, ss, "test-asymmetric")
	if err != nil {
		t.Error(err)
	}
	defer RemoveKey(ctx, ss, pbk)
	defer RemoveKey(ctx, ss, pvk)

	cipher, decrypted := []byte{}, []byte{}

	t.Run("Asymmetric", func(t *testing.T) {
		iv := genIV(16)

		t.Run("Encrypt", func(t *testing.T) { // encrypt with public key
			cipher, err = Encrypt(ctx, ss, pbk, pkcs11.CKM_RSA_PKCS, []byte(plainText), iv)
			if err != nil {
				t.Error(err)
			}
			t.Logf("cipher: %s", base64.StdEncoding.EncodeToString(cipher))
		})
		t.Run("Decrypt", func(t *testing.T) { // decrypt with private key
			decrypted, err = Decrypt(ctx, ss, pvk, pkcs11.CKM_RSA_PKCS, cipher, iv)
			if err != nil {
				t.Error(err)
			}
			t.Logf("decrypted: %s", decrypted)
		})
	})
}
