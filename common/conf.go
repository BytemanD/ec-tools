package common

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/BytemanD/easygo/pkg/global/logging"
)

var (
	CONF      ConfGroup
	CONF_FILE string
)
var CONF_FILES = []string{
	filepath.Join("etc", "ec-tools.yaml"),
	filepath.Join("/etc/ec-tools", "ec-tools.yaml"),
}

type ConfGroup struct {
	Debug      bool       `yaml:"debug"`
	Auth       Auth       `yaml:"auth"`
	Ec         Ec         `yaml:"ec"`
	Iperf      Iperf      `yaml:"iperf"`
	TestServer TestServer `yaml:"testServer"`
}

type Default struct {
}

type Ec struct {
	Flavor           string `yaml:"flavor"`
	Image            string `yaml:"image"`
	BootWithBdm      bool   `yaml:"bootWithBdm"`
	AvailabilityZone string `yaml:"availabilityZone"`
	Network          string `yaml:"network"`
}

type Auth struct {
	Url             string            `yaml:"url"`
	RegionName      string            `yaml:"regionName"`
	User            map[string]string `yaml:"user"`
	Project         map[string]string `yaml:"project"`
	TokenExpireTime int               `yaml:"tokenExpireTime"`
}

type Iperf struct {
	GuestPath     string `yaml:"guestPath"`
	LocalPath     string `yaml:"localPath"`
	ServerOptions string `yaml:"serverOptions"`
	ClientOptions string `yaml:"clientOptions"`
	// 输出QOS结果时，自动转化带宽单位
	ConvertBandwidthUnits bool `yaml:"convertBandwidthUnits"`
}

type TestServer struct {
	Times           int  `yaml:"times"`
	ContinueIfError bool `yaml:"continueIfError"`
	Workers         int  `yaml:"workers"`
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

func LoadConf(configFile string) error {
	viper.SetConfigType("yaml")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("ec-tools.yaml")
		viper.AddConfigPath("./etc")
		viper.AddConfigPath("/etc/ec-tools")
	}
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.Unmarshal(&CONF)
	return nil
}

func DumpConf() string {
	b, err := yaml.Marshal(CONF)
	if err != nil {
		logging.Fatal("dumpl conf failed, %s", err)
	}
	return string(b)
}

func LogLines() {
	logging.Debug("******************** config ********************")
	groupTypes, groupValues := reflect.TypeOf(CONF), reflect.ValueOf(CONF)
	for groupNum := 0; groupNum < groupTypes.NumField(); groupNum++ {
		optionTypes := groupTypes.Field(groupNum)
		options := groupValues.Field(groupNum)
		if options.Kind() == reflect.Struct {
			for num := 0; num < options.NumField(); num++ {
				logging.Debug("%-34s = %v",
					optionTypes.Name+"."+optionTypes.Type.Field(num).Name,
					options.Field(num))
			}
			continue
		}
		logging.Debug("%-34s = %v", optionTypes.Name, options)
	}
	logging.Debug("************************************************")
}
