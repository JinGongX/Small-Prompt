//go:build !darwin && !windows
// +build !darwin,!windows

package platform

// clipboard.go
// 统一接口

func GetPlainText() string
func GetHtmlText() string
func GetImageBase64() string
func CopyText(text string) bool
func CopyHTML(html string) bool
func CopyImage(base64Image string) bool

// func InitHotKey(window application.Window)
// func InitHotKeyDynamic(window application.Window, keyCode, modifiers uint32)
// func RegisterHotKeyDynamic(keyCode, modifiers uint32) bool
// func RegisterHotKey() bool
// func UnregisterHotKey() bool
func NewClipboardBridge() *ClipboardBridge

type ClipboardBridge struct{}

func (cb *ClipboardBridge) CopyText(text string) bool

func HideDock()

func CheckAccessibilityPermission() bool {
	// 在非 macOS 和 Windows 平台上，通常不需要检查辅助功能权限
	// 但可以在这里添加其他逻辑
	return true // 默认返回 true
}
func TriggerAccessibilityPrompt() bool {
	return true //   在非 macOS 和 Windows 平台上，通常不需要检查辅助功能权限
}
