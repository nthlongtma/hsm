package http_server

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
