package common

type Config struct {
	WalletPath       string `yaml:"WalletPath"`
	HLConfigPath     string `yaml:"HLConfigPath"`
	HLWalletIdentity string `yaml:"HLWalletIdentity"`
	PublicKeyPath    string `yaml:"PublicKeyPath"`
	PrivateKeyPath   string `yaml:"PrivateKeyPath"`
	ChannelName      string `yaml:"ChannelName"`
	ContractName     string `yaml:"ContractName"`
	OrgMspId         string `yaml:"OrgMspId"`
	DbUser           string `yaml:"DbUser"`
	DbName           string `yaml:"DbName"`
	DbPassword       string `yaml:"DbPassword"`
	DbHost           string `yaml:"DbHost"`
	DbPort           string `yaml:"DbPort"`
	DbSslmode        string `yaml:"DbSslmode"`
	GraphIdPrefix    string `yaml:"GraphIdPrefix"`
	GrpcServerPort   int    `yaml:"GrpcServerPort"`
}
