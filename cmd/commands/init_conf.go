package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/BytemanD/easygo/pkg/global/logging"
	"github.com/BytemanD/ec-tools/common"
)

var output string

var DumpConf = &cobra.Command{
	Use:   "dump-conf",
	Short: "生成配置文件",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		yamlData := common.DumpConf()
		if output == "" {
			fmt.Println(yamlData)
		} else {
			fi, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0666)
			defer fi.Close()
			if err != nil {
				logging.Fatal("打开文件失败 %s", err)
			}
			fi.Write([]byte(yamlData))
		}
	},
}

func init() {
	DumpConf.Flags().StringVarP(&output, "output", "o", "", "保存文件路径")
}
