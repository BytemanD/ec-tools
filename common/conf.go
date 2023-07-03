package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"gopkg.in/yaml.v3"
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
	ServerOptions string `yaml:"serverOptions"`
	ClientOptions string `yaml:"clientOptions"`
	// 输出QOS结果时，自动转化带宽单位
	ConvertBandwidthUnits bool `yaml:"convertBandwidthUnits"`
}

type TestServer struct {
	Times int `yaml:"times"`
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()

}

func LoadConf(confFiles []string) error {
	if len(confFiles) == 0 {
		confFiles = CONF_FILES
	}
	for _, file := range confFiles {
		if !fileExists(file) {
			continue
		}
		CONF_FILE = file
		logging.Info("load conf from %s", file)
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("load config %s failed: %s", file, err)
		}
		yaml.Unmarshal(bytes, &CONF)
		break
	}
	if CONF_FILE == "" {
		return fmt.Errorf("config file not found, find paths: %v", confFiles)
	}
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
		if options.Kind() != reflect.Struct {
			logging.Debug("%-34s = %v", optionTypes.Name, options)
			continue
		}
		for num := 0; num < options.NumField(); num++ {
			logging.Debug("%-34s = %v",
				optionTypes.Name+"."+optionTypes.Type.Field(num).Name,
				options.Field(num))
		}
	}
	logging.Debug("************************************************")
}
