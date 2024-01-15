package miniblog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/miniblog/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"myMiniBlog/internal/miniblog/controller/v1/user"
	"myMiniBlog/internal/miniblog/store"
	"myMiniBlog/internal/pkg/known"
	"myMiniBlog/internal/pkg/log"
	"myMiniBlog/internal/pkg/middleware"
	proto "myMiniBlog/pkg/proto/miniblog/v1"
	"myMiniBlog/pkg/token"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfgFile string

// NewMiniBlogCommand 创建一个 *cobra.Command 对象. 之后，可以使用 Command 对象的 Execute 方法来启动应用程序.
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
		// 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果 `--version=true`，则打印版本并退出
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

	// 在这里您将定义标志和配置设置。

	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog configuration file. Empty string for no configuration file.")

	// Cobra 也支持本地标志，本地标志只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 添加 --version 标志
	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

// run 函数是实际的业务代码入口函数.
func run() error {
	if err := initStore(); err != nil {
		return err
	}
	// 打印所有的配置项及其值
	settings, _ := json.Marshal(viper.AllSettings())
	log.Infow(string(settings))
	token.Init(viper.GetString("jwt-secret"), known.XUsernameKey)
	// 打印 db -> username 配置项的值
	log.Infow(viper.GetString("db.username"))
	gin.SetMode(viper.GetString("runMode"))
	g := gin.New()
	// 全局中间件使用
	mws := []gin.HandlerFunc{gin.Recovery(), middleware.Cors, middleware.RequestID()}
	g.Use(mws...)
	installRouters(g)
	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	server := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}
	grpcsrv := startGRPCServer()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalw(err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("shutting down server")
	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	grpcsrv.GracefulStop()

	log.Infow("Server exiting")

	return nil
}

// 启动grpc服务
func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	grpcsrv := grpc.NewServer()
	proto.RegisterMiniBlogServer(grpcsrv, user.New(store.S, nil))
	log.Infow("Start to Listening the incoming requests on grpc address", "addr", viper.GetString("grpc.addr"))
	go func() {
		if err := grpcsrv.Serve(lis); err != nil {
			log.Fatalw(err.Error())
		}
	}()
	return grpcsrv
}
