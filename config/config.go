package config

import (
	"flag"
	"path"
)

const (
	ListenPortClearDefault = ":80"
	ListenPortDefault      = ":443"
	DataPathDefault        = "."
	HtmlServerDefault      = "html"
	CallingDataFileDefault = "callings.json"
	MembersDataFileDefault = "members.json"

	CookieName = "id"

	MaxCallings = 300
	MaxMembers  = 500
)

type Config struct {
	ListenPort     string
	DataPath       string
	CallingFile    string
	MembersFile    string
	HtmlServerPath string
}

func ParseConfig() Config {
	config := Config{}

	flag.StringVar(&config.DataPath, "data", DataPathDefault, "The path to the data and html files.")
	flag.StringVar(&config.ListenPort, "listen", ListenPortDefault, "Listen port for TLS. e.g. :443")
	flag.Parse()

	config.HtmlServerPath = path.Join(config.DataPath, HtmlServerDefault)
	config.CallingFile = CallingDataFileDefault
	config.MembersFile = MembersDataFileDefault

	return config
}
