package command

type Server struct {
	ServerName     string `json:"server"`
	Hostname       string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Password       string `json:"pass"`
	PrivateKeyPath string
}

type Config struct {
	Cmd     string
	Servers []Server
}
