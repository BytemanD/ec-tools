package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/common"
)

var InitConf = &cobra.Command{
	Use:   "init-conf [output]",
	Short: "生成配置文件",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		yamlData := common.InitConf()
		if len(args) == 0 {
			fmt.Println(yamlData)
		} else {
			fi, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, 0666)
			defer fi.Close()
			if err != nil {
				logging.Fatal("打开文件失败 %s", err)
			}
			fi.Write([]byte(yamlData))
		}
	},
}
