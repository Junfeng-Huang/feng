package command

import (
	"feng/framework/cobra"
)

// 初始化app命令及其子命令
// func initAppCommand() *cobra.Command {

// }

var AppCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}
