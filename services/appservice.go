package services

import (
	"changeme/platform"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type AppService struct {
	App              *application.App
	repo             *SuiStore
	SecondWindow     *application.WebviewWindow
	TipsWindow       *application.WebviewWindow
	Mainwindow       *application.WebviewWindow
	TipsLoadPosition bool
	Tray             *application.SystemTray
	I18nFS           embed.FS
	Menu             *application.Menu
	MenuQuit         *application.MenuItem
	MenuSetting      *application.MenuItem
	MenuWriteprompts *application.MenuItem
	MenuWetips       *application.MenuItem
	MenuVersion      *application.MenuItem
}

func (a *AppService) SetApp(app *application.App, i18nFS embed.FS, repo *SuiStore) {
	a.App = app
	a.I18nFS = i18nFS
	a.repo = repo

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	secondWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "second",
		Name:      "second",
		Width:     340,
		Height:    234,
		MaxWidth:  340,
		MaxHeight: 234,
		MinWidth:  340,
		MinHeight: 234,
		Frameless: true, // ✅ 去除边框
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			//Backdrop:                application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHidden,
			Backdrop: application.MacBackdropTransparent, // 可选：背景透明
		},
		BackgroundType: application.BackgroundTypeTransparent, //windows系统下需要设置背景类型为透明
		//BackgroundColour: application.NewRGB(27, 38, 54),
		URL: "/#/second",
		//Hidden: true, // ✅ 初始状态为隐藏
	})

	mainwindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		//BackgroundColour: application.NewRGB(27, 38, 54),
		URL:    "/",
		Hidden: true, // ✅ 初始状态为隐藏
	})
	tipsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Tips Window",
		Name:      "tips",
		Width:     360,
		Height:    160,
		MaxWidth:  360,
		MaxHeight: 160,
		MinWidth:  360,
		MinHeight: 160,
		Frameless: true, // ✅ 去除边框
		Mac: application.MacWindow{
			//InvisibleTitleBarHeight: 50,
			TitleBar: application.MacTitleBarHidden,
			Backdrop: application.MacBackdropTransparent, // 可选：背景透明
		},
		BackgroundType: application.BackgroundTypeTransparent, //windows系统下需要设置背景类型为透明
		URL:            "/#/tips",
		Hidden:         true, // ✅ 初始状态为隐藏
		AlwaysOnTop:    true, // ✅ 窗口始终在最前
	})

	a.Mainwindow = mainwindow
	a.SecondWindow = secondWindow
	a.TipsWindow = tipsWindow
	a.TipsLoadPosition = false

	//
	secondWindow.OnWindowEvent(events.Common.WindowShow, func(e *application.WindowEvent) {
		screen, _ := secondWindow.GetScreen() // 获取屏幕信息
		x := screen.X + screen.Size.Width - 340
		if runtime.GOOS == "darwin" {
			x = int(float64(x) * float64(screen.ScaleFactor)) // macOS 需要考虑屏幕缩放
		}
		secondWindow.SetPosition(x, 10)
	})

	RegisterWindowAndCallback(2, tipsWindow, func() {
		if !a.TipsLoadPosition {
			a.TipsLoadPosition = true
			TipsWindowPosition(a)
		}
		fmt.Println("🔥 tipsWindow 被触发")
		platform.ActivateApp()
		tipsWindow.Show()
		tipsWindow.Focus()
	})

	RegisterWindowAndCallback(1, secondWindow, func() {
		fmt.Println("🔥 secondWindow 被触发")
		if secondWindow == nil {
			fmt.Println("🔥 secondWindow is nil")
			return
		}
		SecondWindowPosition(a)
		platform.ActivateApp()
		secondWindow.Show()
		secondWindow.Focus() // 聚焦窗口
	})

	a.Mainwindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		// Prevent the window from closing
		fmt.Println("🔥 mainwindow is closing, but we hide it instead of closing.")
		a.Mainwindow.Hide()
		e.Cancel() // 取消关闭事件
	})
	a.SecondWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		// Prevent the window from closing
		fmt.Println("🔥 secondWindow is closing, but we hide it instead of closing.")
		a.SecondWindow.Hide()
		e.Cancel() // 取消关闭事件
	})

}

func (a *AppService) OpenSecondWindow() {
	if a.SecondWindow != nil {
		fmt.Println("[DEBUG] secondWindow is not nil")
		SecondWindowPosition(a)
		platform.ActivateApp() //  windows系统下需要调用 ActivateApp 来激活应用，否则无法聚焦
		a.SecondWindow.Show()  //.Hide()
		a.SecondWindow.Focus() // 聚焦窗口
	} else {
		fmt.Println("[ERROR] secondWindow is nil")
	}
}

func (a *AppService) OpenTipsWindow() {
	if a.TipsWindow != nil {
		fmt.Println("[DEBUG] tipsWindow is not nil")
		if !a.TipsLoadPosition {
			a.TipsLoadPosition = true
			TipsWindowPosition(a)
		}
		platform.ActivateApp() //  windows系统下需要调用 ActivateApp 来激活应用，否则无法聚焦
		a.TipsWindow.Show()    //.Hide()
		a.TipsWindow.Focus()   // 聚焦窗口 windows系统下需要调用 ActivateApp 来激活应用，否则无法聚焦
		//a.tipsWindow.Restore() // 解除最小化状态（尤其是 Windows）
	} else {
		fmt.Println("[ERROR] tipsWindow is nil")
	}
}
func (a *AppService) HideTipsWindow() {
	if a.TipsWindow != nil {
		a.TipsWindow.Hide()
	}
}

func LoadLang(i18nemb embed.FS, lang string) (*Lang, error) {
	file := fmt.Sprintf("frontend/src/locales/%s.json", lang)

	data, err := i18nemb.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var i18nva Lang
	err = json.Unmarshal(data, &i18nva)

	return &i18nva, err
}

// 配置菜单栏菜单
func LoadNewAppMenu(a *AppService, langvalue string) {
	fmt.Println("🔄 正在加载新菜单语言:", langvalue)
	lang, langerr := LoadLang(a.I18nFS, langvalue)
	if langerr != nil {
		log.Fatal(langerr)
	}
	a.MenuWetips.SetLabel(lang.App.Menu.Wetips)
	a.MenuWriteprompts.SetLabel(lang.App.Menu.Writeprompts)
	a.MenuSetting.SetLabel(lang.App.Menu.Preferences)
	a.MenuVersion.SetLabel(lang.App.Menu.Version)
	a.MenuQuit.SetLabel(lang.App.Menu.Quit)
	//a.Tray.SetMenu(menu)
	if runtime.GOOS == "darwin" {
		a.Menu.Update() // 刷新系统托盘菜单显示
	} else {
		a.Tray.SetMenu(a.Menu) // 重新设置菜单以刷新显示
	}
}

func InitTray(a *AppService, iconBytes []byte, i18nemb embed.FS, langvalue string) {
	a.Tray = a.App.SystemTray.New()
	a.Tray.SetIcon(iconBytes)
	if runtime.GOOS == "darwin" {
		a.Tray.SetTemplateIcon(iconBytes)
	}
	LoadAppMenu(a, langvalue)
}

func LoadAppMenu(a *AppService, langvalue string) {
	lang, langerr := LoadLang(a.I18nFS, langvalue)
	if langerr != nil {
		log.Fatal(langerr)
	}

	a.Menu = application.NewMenu()
	a.MenuWetips = a.Menu.Add(lang.App.Menu.Wetips).OnClick(func(ctx *application.Context) {
		SecondWindowPosition(a)

		platform.ActivateApp()
		a.SecondWindow.Show()
		a.SecondWindow.Focus()
	})
	a.MenuWriteprompts = a.Menu.Add(lang.App.Menu.Writeprompts).OnClick(func(ctx *application.Context) {
		if !a.TipsLoadPosition {
			a.TipsLoadPosition = true
			TipsWindowPosition(a)
		}

		platform.ActivateApp()
		a.TipsWindow.Show()
		a.TipsWindow.Focus()
	})
	a.MenuSetting = a.Menu.Add(lang.App.Menu.Preferences).OnClick(func(ctx *application.Context) {
		// Handle click
		platform.ActivateApp()
		a.Mainwindow.Show()
		a.Mainwindow.Focus()
	})
	a.Menu.AddSeparator()
	item := a.Menu.Add(lang.App.Menu.Version)
	item.SetEnabled(false)
	a.MenuVersion = item
	a.Menu.AddSeparator()
	// menu.Add("关于应用").OnClick(func(ctx *application.Context) {
	// 	app.Quit()
	// })
	a.MenuQuit = a.Menu.Add(lang.App.Menu.Quit).OnClick(func(ctx *application.Context) {
		a.App.Quit()
	})

	a.Tray.SetMenu(a.Menu)
}

// 设定第二窗口位置
func SecondWindowPosition(a *AppService) {
	screen, _ := a.SecondWindow.GetScreen() // 获取屏幕信息
	x := screen.X + screen.Size.Width - 340
	if runtime.GOOS == "darwin" {
		x = int(float64(x) * float64(screen.ScaleFactor)) // macOS 需要考虑屏幕缩放
	}
	a.SecondWindow.SetPosition(x, 10)
}

// 设定写提示窗口位置
func TipsWindowPosition(a *AppService) {
	screen, _ := a.TipsWindow.GetScreen()          // 获取屏幕信息
	x := screen.X + (screen.Size.Width-360)/2      // 居中显示
	y := screen.Y + screen.Size.Height - 160 - 160 // 屏幕高度的1/4位置
	if runtime.GOOS == "darwin" {
		x = int(float64(x) * float64(screen.ScaleFactor)) // macOS 需要考虑屏幕缩放
		y = int(float64(y) * float64(screen.ScaleFactor))
	}
	a.TipsWindow.SetPosition(x, y)
}

// 设置语言并更新菜单显示
func (a *AppService) SetLanguage(lang string) error {
	current, _ := a.repo.GetAppConfig("language")
	if current == lang {
		return nil
	}
	LoadNewAppMenu(a, lang) // 更新菜单语言
	return a.repo.SetAppConfig("language", lang)
}
