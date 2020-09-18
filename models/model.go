package models

type (
	Config struct {
		ModulePath string `mapstructure:"module_path"`
		HSM        HSM    `mapstructure:"hsm"`
	}

	HSM struct {
		SlotID   uint   `mapstructure:"slot_id"`
		Pin      string `mapstructure:"pin"`
		KeyType  string `mapstructure:"key_type"`
		N2kLabel string `mapstructure:"n2k_label"`
		IVSize   int    `mapstructure:"iv_size"`
	}
)

const (
	TypeUnknown     Type = "UNKNOWN"
	TypeN2KKGK      Type = "N2K_KGK"
	TypeN2KInternal Type = "N2K_INTERNAL"
)

type (
	EncryptRequest struct {
		Type      Type   `json:"type"`
		PlainText string `json:"plainText"`
	}
	EncryptResponse struct {
		ErrorCode    string `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
		CipherText   string `json:"cipherText"`
	}

	DecryptRequest struct {
		Type       Type   `json:"type"`
		CipherText string `json:"cipherText"`
	}
	DecryptResponse struct {
		ErrorCode    string `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
		PlainText    string `json:"plainText"`
	}

	Type string
)
