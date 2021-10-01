package config

import (
	"flag"
	"path"
)

const (
	ListenPortDefault           = ":40600"
	DataPathDefault             = "."
	HtmlServerDefault			= "html"
	CallingDataPathDefault      = "callings.csv"
	MembersDataPathDefault      = "members.csv"
	CallingModelDataPathDefault = "callings_model.csv"
	MembersModelDataPathDefault = "members_model.csv"

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
	HtmlServerPath string
}

func ParseConfig() Config {
	config := Config{}

	flag.StringVar(&config.DataPath, "data", DataPathDefault, "The path to the data files.")
	flag.StringVar(&config.ListenPort, "listen", ListenPortDefault, "Listen port. e.g. :8080")
	flag.Parse()

	config.CallingDataPath = path.Join(config.DataPath, CallingDataPathDefault)
	config.MembersDataPath = path.Join(config.DataPath, MembersDataPathDefault)
	config.CallingModelDataPath = path.Join(config.DataPath, CallingModelDataPathDefault)
	config.MembersModelDataPath = path.Join(config.DataPath, MembersModelDataPathDefault)
	config.HtmlServerPath = path.Join(config.DataPath, HtmlServerDefault)

	return config
}
