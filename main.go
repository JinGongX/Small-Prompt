package main

import (
	platform "changeme/platform/mac"
	"changeme/services"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/icon.png
var iconFS embed.FS

type AppService struct {
	App          *application.App
	secondWindow *application.WebviewWindow
	tipsWindow   *application.WebviewWindow
	mainwindow   *application.WebviewWindow
}

func (a *AppService) SetApp(app *application.App) {
	a.App = app
}

func (a *AppService) OpenSecondWindow() {
	if a.secondWindow != nil {
		fmt.Println("[DEBUG] secondWindow is not nil")
		screen, _ := a.secondWindow.GetScreen()                            // 获取屏幕信息
		a.secondWindow.SetPosition((screen.X+screen.Size.Width-340)*2, 10) //+10
		a.secondWindow.Show()                                              //.Hide()
	} else {
		fmt.Println("[ERROR] secondWindow is nil")
	}
}

func (a *AppService) OpenTipsWindow() {
	if a.tipsWindow != nil {
		fmt.Println("[DEBUG] tipsWindow is not nil")
		a.tipsWindow.Show() //.Hide()
		//a.tipsWindow.Restore() // 解除最小化状态（尤其是 Windows）
	} else {
		fmt.Println("[ERROR] tipsWindow is nil")
	}
}
func (a *AppService) HideTipsWindow() {
	if a.tipsWindow != nil {
		a.tipsWindow.Hide()
	}
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	appservice := &AppService{}
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
		Name:        "wetips",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(appservice),
			application.NewService(suiService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	mainwindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		//BackgroundColour: application.NewRGB(27, 38, 54),
		URL: "/",
	})

	secondWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "second",
		Name:      "second",
		Width:     340,
		Height:    224,
		MaxWidth:  340,
		MaxHeight: 224,
		MinWidth:  340,
		MinHeight: 224,
		Frameless: true, // ✅ 去除边框
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			//Backdrop:                application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHidden,
			Backdrop: application.MacBackdropTransparent, // 可选：背景透明
		},
		//BackgroundColour: application.NewRGB(27, 38, 54),
		URL:    "/#/second",
		Hidden: true, // ✅ 初始状态为隐藏
	})
	tipsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Tips Window",
		Name:      "tips",
		Width:     340,
		Height:    160,
		MaxWidth:  340,
		MaxHeight: 160,
		MinWidth:  340,
		MinHeight: 160,
		Frameless: true, // ✅ 去除边框
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			TitleBar:                application.MacTitleBarHidden,
			Backdrop:                application.MacBackdropTransparent, // 可选：背景透明
		},
		URL:    "/#/tips",
		Hidden: true, // ✅ 初始状态为隐藏
	})

	appservice.mainwindow = mainwindow
	appservice.secondWindow = secondWindow
	appservice.tipsWindow = tipsWindow

	appservice.SetApp(app)

	// db, erra := services.NewSuiStore()
	// defer db.Close()
	// if erra != nil {
	// }
	scheduler := services.NewPromptScheduler(app.Context(), suiService)
	suiService.Start()
	// App 启动即恢复调度
	scheduler.Recalculate()
	services.Init(app, scheduler)
	services.RegisterWindowAndCallback(2, tipsWindow, func() {
		fmt.Println("🔥 tipsWindow 被触发")
		platform.ActivateApp()
		tipsWindow.Show()
		tipsWindow.Focus()
	})

	services.RegisterWindowAndCallback(1, secondWindow, func() {
		fmt.Println("🔥 secondWindow 被触发")
		if secondWindow == nil {
			fmt.Println("🔥 secondWindow is nil")
			return
		}
		// screen, _ := secondWindow.GetScreen()                            // 获取屏幕信息
		// secondWindow.SetPosition((screen.X+screen.Size.Width-340)*2, 10) //+10
		platform.ActivateApp()
		secondWindow.Show()
		secondWindow.Focus() // 聚焦窗口
	})

	services.LoadAndRegisterHotkeysFrom(suiService, 1)
	services.LoadAndRegisterHotkeysFrom(suiService, 2)

	// mainwindow.RegisterHook(events.Common.WindowClosing, func(ctx *application.WindowEvent) {
	// 	mainwindow.Hide() // 隐藏窗口而不是关闭
	// 	fmt.Println("🔥 mainwindow is closing, but we hide it instead of closing.")
	// 	ctx.Cancel() // 取消关闭事件
	// })
	mainwindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		// Prevent the window from closing
		fmt.Println("🔥 mainwindow is closing, but we hide it instead of closing.")
		mainwindow.Hide()
		e.Cancel() // 取消关闭事件
	})
	secondWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		// Prevent the window from closing
		fmt.Println("🔥 secondWindow is closing, but we hide it instead of closing.")
		secondWindow.Hide()
		e.Cancel() // 取消关闭事件
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		// for {
		// 	now := time.Now().Format(time.RFC1123)
		// 	app.Event.Emit("time", now)
		// 	time.Sleep(time.Second)
		// }
	}()

	iconBytes, _ := iconFS.ReadFile("assets/icon.png")
	systray := app.SystemTray.New()
	//systray.SetLabel("My App")
	systray.SetIcon(iconBytes)
	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(iconBytes)
	}
	menu := application.NewMenu()
	menu.Add("打开面板").OnClick(func(ctx *application.Context) {
		screen, _ := secondWindow.GetScreen()                            // 获取屏幕信息
		secondWindow.SetPosition((screen.X+screen.Size.Width-340)*2, 10) //+10
		platform.ActivateApp()
		secondWindow.Show()
		secondWindow.Focus()
	})
	menu.Add("去写提示").OnClick(func(ctx *application.Context) {
		platform.ActivateApp()
		tipsWindow.Show()
		tipsWindow.Focus()
	})
	menu.Add("偏好设置").OnClick(func(ctx *application.Context) {
		// Handle click
		platform.ActivateApp()
		mainwindow.Show()
		mainwindow.Focus()
	})
	menu.AddSeparator()
	item := menu.Add("版本:v0.1")
	item.SetEnabled(false)
	menu.AddSeparator()
	// menu.Add("关于应用").OnClick(func(ctx *application.Context) {
	// 	app.Quit()
	// })
	menu.Add("退出应用").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systray.SetMenu(menu)

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
