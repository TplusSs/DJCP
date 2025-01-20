package cmd

import (
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run"
)

// linuxCmd represents the linux command
var pgsqlCmd = &cobra.Command{
	Use:   "pgsql",
	Short: "运行采集PostgreSQL安全配置核查功能",
	Long:  `基于远程登录功能,通过多线程的方法批量进行采集`,
	Run:   run.Pgsqlstart,
}

func init() {
	rootCmd.AddCommand(pgsqlCmd)
	pgsqlCmd.Flags().StringP("ip", "i", global.CmdPgsqlPath, "此参数是指定待远程采集的IP文件位置")
	pgsqlCmd.Flags().StringP("spript", "s", global.Split, "此参数是指定IP文件中的分隔字符")
	pgsqlCmd.Flags().StringP("value", "v", "", "此参数是单次执行")
	pgsqlCmd.Flags().BoolP("echo", "e", false, "此参数是控制控制台是否输出结果,默认不进行输出")
}
