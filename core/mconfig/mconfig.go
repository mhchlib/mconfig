package mconfig

// MConfig ...
type MConfigConfig struct {
	Namspace        string
	RegistryAddress string
	RegistryType    string
	StoreType       string
	StoreAddress    string
	ServerIp        string
	ServerPort      int
}

// NewMConfig ...
func NewMConfig() *MConfigConfig {
	return &MConfigConfig{}
}
