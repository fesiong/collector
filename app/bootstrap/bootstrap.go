package bootstrap

import (
	"collector/app/route"
	"collector/config"
	"collector/core"
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"os/exec"
	"runtime"
	"time"
)

type Bootstrap struct {
	Application *iris.Application
	Port        int
	LoggerLevel string
}

func New(port int, loggerLevel string) *Bootstrap {
	var bootstrap Bootstrap
	bootstrap.Application = iris.New()
	bootstrap.Port = port
	bootstrap.LoggerLevel = loggerLevel

	//crond
	core.Crond()

	return &bootstrap
}

func (bootstrap *Bootstrap) loadGlobalMiddleware() {
	bootstrap.Application.Use(recover.New())
}

func (bootstrap *Bootstrap) loadRoutes() {
	route.Register(bootstrap.Application)
}

func (bootstrap *Bootstrap) Serve() {
	bootstrap.Application.Logger().SetLevel(bootstrap.LoggerLevel)
	bootstrap.loadGlobalMiddleware()
	bootstrap.loadRoutes()
	pugEngine := iris.Django(fmt.Sprintf("%stemplate", config.ExecPath), ".html")

	if config.ServerConfig.Env == "development" {
		//测试环境下动态加载
		pugEngine.Reload(true)
	}

	pugEngine.AddFunc("stampToDate", TimestampToDate)
	bootstrap.Application.RegisterView(pugEngine)

	go Open(fmt.Sprintf("http://127.0.0.1:%d", config.ServerConfig.Port))

	bootstrap.Application.Run(
		iris.Addr(fmt.Sprintf(":%d", bootstrap.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithoutBodyConsumptionOnUnmarshal,
	)
}

func TimestampToDate(in uint, layout string) string {
	t := time.Unix(int64(in), 0)
	return t.Format(layout)
}

func (bootstrap *Bootstrap) Shutdown() error {
	bootstrap.Application.Shutdown(context.Background())

	return nil
}

func Open(uri string) {
	time.Sleep(1 * time.Second)
	var commands = map[string]string{
		"windows": "cmd /c start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}

	run, ok := commands[runtime.GOOS]
	if !ok {
		fmt.Println(fmt.Sprintf("请手动在浏览器中打开网址： %s", uri))
		return
	}

	cmd := exec.Command(run, uri)
	err := cmd.Start()
	if err != nil {
		fmt.Println(fmt.Sprintf("请手动在浏览器中打开网址： %s", uri))
	}
}
