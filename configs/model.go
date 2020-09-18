package configs

type (
	Config struct {
		ModulePath string  `mapstructure:"module_path"`
		HSM        HSM     `mapstructure:"hsm"`
		Servers    Servers `mapstructure:"servers"`
	}

	HSM struct {
		SlotID   uint   `mapstructure:"slot_id"`
		Pin      string `mapstructure:"pin"`
		KeyType  string `mapstructure:"key_type"`
		N2kLabel string `mapstructure:"n2k_label"`
		IVSize   int    `mapstructure:"iv_size"`
	}

	Servers struct {
		HTTP SeverInfo `mapstructure:"http"`
		GRPC SeverInfo `mapstructure:"grpc"`
	}

	SeverInfo struct {
		Port string `mapstructure:"port"`
		Path Path   `mapstructure:"path"`
	}
	Path struct {
		Encrypt string `mapstructure:"encrypt"`
		Decrypt string `mapstructure:"decrypt"`
	}
)
