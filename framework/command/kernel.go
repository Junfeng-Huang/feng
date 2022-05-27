package command

import "feng/framework/cobra"

func RunCommand() {
	var rootCmd = &cobra.Command{
		Use:   "feng",
		Short: "feng 命令",
		Long:  "feng框架提供的命令行工具，使用这个工具执行框架自带命令",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.InitDefaultHelpFlag()
			cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	rootCmd.AddCommand(initAppCommand())
	rootCmd.AddCommand(initDevCommand())
	rootCmd.Execute()
}
