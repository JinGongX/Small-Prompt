//go:build darwin
// +build darwin

package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -framework Vision  -framework Foundation
#include <stdlib.h>
#include <AppKit/AppKit.h>
#include "system.h"
*/
import "C"
import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sync"
	"unsafe"
)

// var mainWindow application.Window

//export hotKeyCallback
// func hotKeyCallback() {
// 	if mainWindow != nil {
// 		mainWindow.Show()
// 		mainWindow.Focus()
// 	}
// }

func GetPlainText() string {
	ptr := C.getPlainText()
	if ptr == nil {
		return ""
	}
	return C.GoString(ptr)
}

func GetHtmlText() string {
	ptr := C.getHtmlText()
	if ptr == nil {
		return ""
	}
	return C.GoString(ptr)
}

func GetImageBase64() string {
	ptr := C.getImageBase64()
	if ptr == nil {
		return ""
	}
	return C.GoString(ptr)
}

// func InitHotKey(window application.Window) {
// 	mainWindow = window
// 	C.RegisterHotKey()
// }

// func InitHotKeyDynamic(window application.Window, keyCode, modifiers uint32) {
// 	mainWindow = window
// 	//C.RegisterHotKeyDynamic()
// 	C.RegisterHotKeyDynamic(C.uint(keyCode), C.uint(modifiers))
// }

// 剪贴板写入
type ClipboardBridge struct{}

func NewClipboardBridge() *ClipboardBridge {
	return &ClipboardBridge{}
}

func CopyText(text string) bool {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return bool(C.setPlainText(ctext))
}

func CopyHTML(html string) bool {
	chtml := C.CString(html)
	defer C.free(unsafe.Pointer(chtml))
	return bool(C.setHtmlText(chtml))
}
func CopyImage(base64 string) bool {
	cimage := C.CString(base64)
	defer C.free(unsafe.Pointer(cimage))
	return bool(C.setImageBase64(cimage))
}

type HotKeyCallback func()

var (
	callbacks = make(map[string]HotKeyCallback)
	mutex     sync.RWMutex
)

//export hotKeyCallback
func hotKeyCallback(keyCode, modifiers C.uint) {
	id := fmt.Sprintf("%d_%d", uint32(keyCode), uint32(modifiers))

	mutex.RLock()
	cb, exists := callbacks[id]
	mutex.RUnlock()

	if exists {
		fmt.Printf("✅ 触发热键: %s\n", id)
		cb()
	} else {
		fmt.Printf("❌ 找不到回调: %s\n", id)
	}
}

// RegisterHotKeyWithCallback 注册热键并绑定 Go 回调
func RegisterHotKeyWithCallback(keyCode, modifiers uint32, cb HotKeyCallback) {
	id := fmt.Sprintf("%d_%d", keyCode, modifiers)

	mutex.Lock()
	callbacks[id] = cb
	mutex.Unlock()
	fmt.Printf("✈️ Registering hotkey: %d %d\n", keyCode, modifiers)
	C.RegisterHotKeyDynamic(C.uint(keyCode), C.uint(modifiers))
}

// 激活当前应用程序
func ActivateApp() {
	C.NSAppActivateIgnoringOtherApps() // 激活当前应用程序
}

// UnregisterHotKey 注销热键（Go -> C）
func UnregisterHotKey(keyCode, modifiers uint32) {
	id := fmt.Sprintf("%d_%d", keyCode, modifiers)

	mutex.Lock()
	delete(callbacks, id)
	mutex.Unlock()

	fmt.Printf("🧹 Unregistering hotkey: %d %d\n", keyCode, modifiers)
	C.UnregisterHotKey(C.uint(keyCode), C.uint(modifiers))
}

// SimulateCmdC triggers a Command+C key press on macOS.
func SimulateCmdC() {
	C.simulateCmdC()
}

func HideDock() {
	C.HideDockIcon()
}

func CheckAccessibilityPermission() bool {
	return bool(C.isAccessibilityEnabled())
}

func TriggerAccessibilityPrompt() bool {
	return bool(C.requestAccessibilityPermission())
}

func RecognizeTextFromImageMac(imagePath string) (string, error) {
	cpath := C.CString(imagePath)
	defer C.free(unsafe.Pointer(cpath))

	result := C.VisionOCR(cpath)
	defer C.free(unsafe.Pointer(result))

	return C.GoString(result), nil
}

func RecognizeImageBase64(base64str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}
	tmp := C.CBytes(data)
	defer C.free(tmp)
	result := C.VisionOCRFromMemory(tmp, C.int(len(data)))
	defer C.free(unsafe.Pointer(result))
	return C.GoString(result), nil
}

// RecognizeImageAndCopyToClipboard 解码 Base64 图片并识别，然后写入剪贴板
func RecognizeImageAndCopyToClipboard(base64str string) error {
	data, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return err
	}

	ptr := C.CBytes(data)
	defer C.free(ptr)

	res := C.VisionOCRFromMemory(ptr, C.int(len(data)))
	if res == nil {
		return errors.New("OCR failed: result is nil")
	}
	defer C.free(unsafe.Pointer(res))

	text := C.GoString(res)
	if text == "" {
		return errors.New("OCR returned empty string")
	}

	// 👇 写入系统剪贴板
	CopyText(text)
	return nil
}

// 开机自启
// plist 模板
const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
 "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
</dict>
</plist>
`

// EnableAutoStart 在 macOS 上启用开机自启
func EnableAutoStart(appName string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	// LaunchAgents 目录
	plistDir := filepath.Join(usr.HomeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(plistDir, 0755); err != nil {
		return err
	}

	// plist 文件路径
	plistPath := filepath.Join(plistDir, fmt.Sprintf("%s.plist", appName))

	// 获取可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// 写入 plist
	plistContent := fmt.Sprintf(plistTemplate, appName, exePath)
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return err
	}

	// 加载到 launchctl
	cmd := exec.Command("launchctl", "load", plistPath)
	return cmd.Run()
}

// DisableAutoStart 禁用开机自启
func DisableAutoStart(appName string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	plistPath := filepath.Join(usr.HomeDir, "Library", "LaunchAgents", fmt.Sprintf("%s.plist", appName))

	// 先卸载
	exec.Command("launchctl", "unload", plistPath).Run()

	// 删除 plist
	return os.Remove(plistPath)
}

// IsEnabled 检查是否已启用开机自启
func IsEnabled(appName string) bool {
	usr, err := user.Current()
	if err != nil {
		return false
	}

	plistPath := filepath.Join(usr.HomeDir, "Library", "LaunchAgents", fmt.Sprintf("%s.plist", appName))
	if _, err := os.Stat(plistPath); err == nil {
		return true
	}
	return false
}
