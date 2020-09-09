package models

type (
	Config struct {
		ModulePath string `mapstructure:"module_path"`
		HSM        HSM    `mapstructure:"hsm"`
	}

	HSM struct {
		SlotID     uint   `mapstructure:"slot_id"`
		Pin        string `mapstructure:"pin"`
		KeyType    string `mapstructure:"key_type"`
		TokenLabel string `mapstructure:"token_label"`
	}
)
