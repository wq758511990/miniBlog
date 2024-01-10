package miniblog

import (
	"encoding/json"
	"fmt"
	"github.com/marmotedu/miniblog/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"myMiniBlog/internal/pkg/log"
)

var cfgFile string

func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 指定命令的名字，该名字会出现在帮助信息中
		Use: "miniblog",
		// 命令的简短描述
		Short: "A good Go practical project",
		// 命令的详细描述
		Long: `A good Go practical project, used to create user with basic information.

Find more miniblog information at:
	https://github.com/marmotedu/miniblog#readme`,
		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()
			// 初始化日志
			log.Init(logOptions())
			defer log.Sync() // Sync 将缓存中的日志刷新到磁盘文件中

			return run()
		},
		// 这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}
	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog configuration file. Empty string for no configuration file.")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

// run 函数是实际的业务代码入口函数.
func run() error {
	// 打印所有的配置项及其值
	settings, _ := json.Marshal(viper.AllSettings())
	log.Infow(string(settings))
	// 打印 db -> username 配置项的值
	log.Infow(viper.GetString("db.username"))
	return nil
}
