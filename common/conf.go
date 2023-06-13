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
	Debug bool `yaml:"debug"`
	Ec    Ec   `yaml:"ec"`
}

type Default struct {
}

type Ec struct {
	AuthOpenrc  string `yaml:"authOpenrc"`
	Flavor      string `yaml:"flavor"`
	Image       string `yaml:"image"`
	BootWithBdm bool   `yaml:"bootWithBdm"`
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()

}

func LoadConf() error {
	for _, file := range CONF_FILES {
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
		return fmt.Errorf("config file not found")
	}
	return nil
}

func InitConf() string {
	b, err := yaml.Marshal(CONF)
	if err != nil {
		os.Exit(1)
	}
	return string(b)
}

func LogConf(ecConf ConfGroup) {
	logging.Debug("*************** config ***************")
	groupTypes := reflect.TypeOf(ecConf)
	groupvalues := reflect.ValueOf(ecConf)
	for groupNum := 0; groupNum < groupTypes.NumField(); groupNum++ {
		if groupvalues.Field(groupNum).Kind() != reflect.Struct {
			logging.Debug("config: %s = %v", groupTypes.Field(groupNum).Name, groupvalues.Field(groupNum))
			continue
		}
		types := reflect.TypeOf(ecConf.Ec)
		values := reflect.ValueOf(ecConf.Ec)
		for num := 0; num < values.NumField(); num++ {
			logging.Debug("config: %s.%s = %v",
				groupTypes.Field(groupNum).Name,
				types.Field(num).Name,
				values.Field(num))
		}
	}
	logging.Debug("**************************************")
}
