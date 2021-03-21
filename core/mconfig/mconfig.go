package mconfig

// MConfig ...
type MConfigConfig struct {
	Namspace        string
	RegistryAddress string
	RegistryType    string
	EnableRegistry  bool
	StoreType       string
	StoreAddress    string
	ServerIp        string
	ServerPort      int
}

func NewMConfig() *MConfigConfig {
	return &MConfigConfig{}
}
