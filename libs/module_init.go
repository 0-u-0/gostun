package libs

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"os"
	"log"
)

func Init()  {
	LoadArgsModule()
	LoadConfigurationModule()
	LoadLoggerModule()


}

var (
	App               = kingpin.New("gostun", APP_NAME)
	config            = App.Flag("config", "Configuration file location").PlaceHolder(strings.Join(config_path_array,",")).Short('c').String()
)


func LoadArgsModule() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix(SERVER_TAG)

	App.Version(APP_VERSION)
	App.HelpFlag.Short('h')
	App.VersionFlag.Short('v')
	App.Parse(os.Args[1:])

}