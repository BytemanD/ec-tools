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
	AuthOpenrc     string `yaml:"authOpenrc"`
	Flavor         string `yaml:"flavor"`
	Image          string `yaml:"image"`
	BootWithBdm    bool   `yaml:"bootWithBdm"`
	IperfGuestPath string `yaml:"iperfGuestPath"`
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
		return fmt.Errorf("config file not found")
	}
	return nil
}

func InitConf() string {
	b, err := yaml.Marshal(CONF)
	if err != nil {
		logging.Fatal("加载配置失败 %s", err)
	}
	return string(b)
}

func LogConf() {
	logging.Debug("*************** config ***************")
	groupTypes := reflect.TypeOf(CONF)
	groupvalues := reflect.ValueOf(CONF)
	for groupNum := 0; groupNum < groupTypes.NumField(); groupNum++ {
		if groupvalues.Field(groupNum).Kind() != reflect.Struct {
			logging.Debug("config: %s = %v", groupTypes.Field(groupNum).Name, groupvalues.Field(groupNum))
			continue
		}
		types := reflect.TypeOf(CONF.Ec)
		values := reflect.ValueOf(CONF.Ec)
		for num := 0; num < values.NumField(); num++ {
			logging.Debug("config: %s.%s = %v",
				groupTypes.Field(groupNum).Name,
				types.Field(num).Name,
				values.Field(num))
		}
	}
	logging.Debug("**************************************")
}
