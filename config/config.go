package config

import (
	"flag"
	"path"
)

const (
	ListenPortClearDefault      = ":80"
	ListenPortDefault           = ":443"
	DataPathDefault             = "."
	HtmlServerDefault           = "html"
	CallingDataPathDefault      = "callings.json"
	MembersDataPathDefault      = "members.json"
	CallingModelDataPathDefault = "callings_model.json"
	MembersModelDataPathDefault = "members_model.json"

	MaxCallings = 300
	MaxMembers  = 500
)

type Config struct {
	ListenPort           string
	DataPath             string
	CallingDataPath      string
	MembersDataPath      string
	CallingModelDataPath string
	MembersModelDataPath string
	HtmlServerPath       string
}

func ParseConfig() Config {
	config := Config{}

	flag.StringVar(&config.DataPath, "data", DataPathDefault, "The path to the data and html files.")
	flag.StringVar(&config.ListenPort, "listen", ListenPortDefault, "Listen port for TLS. e.g. :443")
	flag.Parse()

	config.CallingDataPath = path.Join(config.DataPath, CallingDataPathDefault)
	config.MembersDataPath = path.Join(config.DataPath, MembersDataPathDefault)
	config.CallingModelDataPath = path.Join(config.DataPath, CallingModelDataPathDefault)
	config.MembersModelDataPath = path.Join(config.DataPath, MembersModelDataPathDefault)
	config.HtmlServerPath = path.Join(config.DataPath, HtmlServerDefault)

	return config
}
