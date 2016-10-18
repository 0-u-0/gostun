package libs

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"os"
	"log"
)

func Init()  {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix(SERVER_TAG)

	LoadArgsModule()
	LoadConfigurationModule()
	LoadLoggerModule()
	LoadEntryModule()

}

var (
	App               = kingpin.New("gostun", APP_NAME)
	config            = App.Flag("config", "Configuration file location").PlaceHolder(strings.Join(config_path_array,",")).Short('c').String()
	external_ip       = App.Flag("external_ip","TURN Server public/private address mapping, if the server is behind NAT.").Short('x').String()
	min_port 		  = App.Flag("min_port","Lower bound of the UDP port range for relay endpoints allocation.").Default("49152").Int()
	max_port 		  = App.Flag("max_port","Upper bound of the UDP port range for relay endpoints allocation.").Default("65535").Int()
)


func LoadArgsModule() {
	App.Version(APP_VERSION)
	App.HelpFlag.Short('h')
	App.VersionFlag.Short('v')
	App.Parse(os.Args[1:])
}