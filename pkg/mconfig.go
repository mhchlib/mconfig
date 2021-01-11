package pkg

// MConfig ...
type MConfig struct {
	Namspace        string
	RegistryAddress string
	RegistryType    string
	EnableRegistry  bool
	StoreType       string
	StoreAddress    string
	ServerIp        string
	ServerPort      int
}

func NewMConfig() *MConfig {
	return &MConfig{}
}
