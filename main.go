package main

import (
	"changeme/services"
	"embed"
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/locales/*.json
var I18nFS embed.FS

//go:embed assets/icon.png
var iconFS embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	appservice := &services.AppService{}
	suiService, errt := services.NewSuiStore()
	if errt != nil {
		// 处理错误，比如日志或退出
	}

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "Small Prompt",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(appservice),
			application.NewService(suiService),
			application.NewService(services.NewSystemInfo()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// 加上这个就搞定双开问题
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.yourcompany.yourapp", // 用你自己的唯一 ID，推荐反域名格式
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				// 如果主窗口存在，把它恢复 & 聚焦
				if appservice.SecondWindow != nil {
					appservice.SecondWindow.Restore()
					appservice.SecondWindow.Show()
					appservice.SecondWindow.Focus()
				}
			},
		},
		OnShutdown: func() {
			if appservice.Tray != nil {
				appservice.Tray.Destroy()
			}
		},
	})

	appservice.SetApp(app, I18nFS, suiService)

	scheduler := services.NewPromptScheduler(app.Context(), suiService)
	suiService.Start()
	// App 启动即恢复调度
	scheduler.Recalculate()
	services.Init(app, scheduler)

	services.LoadAndRegisterHotkeysFrom(suiService, 1)
	services.LoadAndRegisterHotkeysFrom(suiService, 2)

	langvalue, _ := suiService.GetAppConfig("language") // 获取语言配置
	iconBytes, _ := iconFS.ReadFile("assets/winicon.png")
	if runtime.GOOS == "darwin" {
		iconBytes, _ = iconFS.ReadFile("assets/icon.png")
	}
	services.InitTray(appservice, iconBytes, I18nFS, langvalue)
	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	// go func() {
	// 	// for {
	// 	// 	now := time.Now().Format(time.RFC1123)
	// 	// 	app.Event.Emit("time", now)
	// 	// 	time.Sleep(time.Second)
	// 	// }
	// }()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
