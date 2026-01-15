package restapi

type Client struct {
	Address  string `yaml:"address"`
	AppToken string `yaml:"token" env-required:"true" env:"APP_TOKEN"`
}
