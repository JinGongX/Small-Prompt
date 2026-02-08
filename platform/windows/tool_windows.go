//go:build windows
// +build windows

package platform

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.design/x/hotkey"
	"golang.org/x/sys/windows/registry"
)

const (
	CF_UNICODETEXT = 13
	CF_DIB         = 8
	GMEM_MOVEABLE  = 0x0002
	//WM_HOTKEY      = 0x0312
)

// var (
// 	// ID offset base for hotkeys
// 	hotkeyID uint32 = 1
// )

var (
	user32                     = syscall.NewLazyDLL("user32.dll")
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	openClipboard              = user32.NewProc("OpenClipboard")
	closeClipboard             = user32.NewProc("CloseClipboard")
	getClipboardData           = user32.NewProc("GetClipboardData")
	isClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	registerFormat             = user32.NewProc("RegisterClipboardFormatW")
	globalLock                 = kernel32.NewProc("GlobalLock")
	globalUnlock               = kernel32.NewProc("GlobalUnlock")
	globalSize                 = kernel32.NewProc("GlobalSize")
	emptyClipboard             = user32.NewProc("EmptyClipboard")
	setClipboardData           = user32.NewProc("SetClipboardData")
	globalAlloc                = kernel32.NewProc("GlobalAlloc")

	// modUser32            = syscall.NewLazyDLL("user32.dll")
	// procRegisterHotKey   = modUser32.NewProc("RegisterHotKey")
	// procUnregisterHotKey = modUser32.NewProc("UnregisterHotKey")
	// procGetMessageW      = modUser32.NewProc("GetMessageW")
	// procDispatchMessageW = modUser32.NewProc("DispatchMessageW")
	// procTranslateMessage = modUser32.NewProc("TranslateMessage")

	//callbacks = make(map[uint32]HotKeyCallback)
	//mutex sync.RWMutex

	procSetForeground = user32.NewProc("SetForegroundWindow")
	procShowWindow    = user32.NewProc("ShowWindow")
	procGetForeground = user32.NewProc("GetForegroundWindow")
)

func open() error {
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("failed to open clipboard")
	}
	return nil
}

func close() {
	closeClipboard.Call()
}

func registerClipboardFormat(name string) uint32 {
	ptr, _ := syscall.UTF16PtrFromString(name)
	ret, _, _ := registerFormat.Call(uintptr(unsafe.Pointer(ptr)))
	return uint32(ret)
}

// ✅ 读取纯文本
func GetPlainText() string {
	if err := open(); err != nil {
		return ""
	}
	defer close()

	ret, _, _ := getClipboardData.Call(uintptr(CF_UNICODETEXT))
	if ret == 0 {
		return ""
	}

	lock, _, _ := globalLock.Call(ret)
	if lock == 0 {
		return ""
	}
	defer globalUnlock.Call(ret)

	u16 := (*[1 << 20]uint16)(unsafe.Pointer(lock))
	length := 0
	for u16[length] != 0 {
		length++
	}

	return string(utf16.Decode(u16[:length]))
}

// ✅ 读取 HTML 片段
func GetHtmlText() string {
	if err := open(); err != nil {
		return ""
	}
	defer close()

	cfHTML := registerClipboardFormat("HTML Format")
	available, _, _ := isClipboardFormatAvailable.Call(uintptr(cfHTML))
	if available == 0 {
		return ""
	}

	ret, _, _ := getClipboardData.Call(uintptr(cfHTML))
	if ret == 0 {
		return ""
	}

	size, _, _ := globalSize.Call(ret)
	if size == 0 {
		return ""
	}

	lock, _, _ := globalLock.Call(ret)
	if lock == 0 {
		return ""
	}
	defer globalUnlock.Call(ret)

	buf := (*[1 << 30]byte)(unsafe.Pointer(lock))
	data := make([]byte, size)
	copy(data, buf[:size])

	full := string(data)

	startMarker := "<!--StartFragment-->"
	endMarker := "<!--EndFragment-->"
	startIdx := strings.Index(full, startMarker)
	endIdx := strings.Index(full, endMarker)
	if startIdx == -1 || endIdx == -1 {
		return full // fallback
	}
	startIdx += len(startMarker)
	return full[startIdx:endIdx]
}

// ✅ 读取图像并转为 base64 编码的 PNG（从 CF_DIB）
func GetImageBase64() string {
	if err := open(); err != nil {
		return ""
	}
	defer close()

	available, _, _ := isClipboardFormatAvailable.Call(CF_DIB)
	if available == 0 {
		return ""
	}

	ret, _, _ := getClipboardData.Call(CF_DIB)
	if ret == 0 {
		return ""
	}

	size, _, _ := globalSize.Call(ret)
	if size == 0 {
		return ""
	}

	lock, _, _ := globalLock.Call(ret)
	if lock == 0 {
		return ""
	}
	defer globalUnlock.Call(ret)

	// Windows BITMAPINFOHEADER 结构
	type BITMAPINFOHEADER struct {
		Size          uint32
		Width         int32
		Height        int32
		Planes        uint16
		BitCount      uint16
		Compression   uint32
		SizeImage     uint32
		XPelsPerMeter int32
		YPelsPerMeter int32
		ClrUsed       uint32
		ClrImportant  uint32
	}

	header := (*BITMAPINFOHEADER)(unsafe.Pointer(lock))
	if header.BitCount != 32 {
		return ""
	}

	width := int(header.Width)
	height := int(header.Height)

	// pointer to pixel data (after header)
	pixels := (*[1 << 30]byte)(unsafe.Pointer(uintptr(lock) + uintptr(header.Size)))

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// DIB is bottom-up, need to invert rows
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			src := ((height-1-y)*width + x) * 4
			r := pixels[src+2]
			g := pixels[src+1]
			b := pixels[src+0]
			a := pixels[src+3]
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// 写入剪贴板

// ✅ 写入纯文本
func CopyText(text string) bool {
	if err := open(); err != nil {
		return false
	}
	defer close()
	emptyClipboard.Call()

	data := utf16.Encode([]rune(text + "\x00"))
	size := len(data) * 2
	hMem, _, _ := globalAlloc.Call(GMEM_MOVEABLE, uintptr(size))
	if hMem == 0 {
		return false
	}

	ptr, _, _ := globalLock.Call(hMem)
	if ptr == 0 {
		return false
	}
	copy((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:len(data)], data)
	globalUnlock.Call(hMem)

	_, _, _ = setClipboardData.Call(CF_UNICODETEXT, hMem)
	return true
}

// ✅ 写入 HTML
func CopyHTML(html string) bool {
	if err := open(); err != nil {
		return false
	}
	defer close()
	emptyClipboard.Call()

	cfHTML := registerClipboardFormat("HTML Format")

	header := `Version:0.9
StartHTML:00000097
EndHTML:00000197
StartFragment:00000133
EndFragment:00000161
<html><body><!--StartFragment-->`
	footer := `<!--EndFragment--></body></html>`
	full := header + html + footer

	data := []byte(full)
	hMem, _, _ := globalAlloc.Call(GMEM_MOVEABLE, uintptr(len(data)))
	if hMem == 0 {
		return false
	}

	ptr, _, _ := globalLock.Call(hMem)
	if ptr == 0 {
		return false
	}
	copy((*[1 << 30]byte)(unsafe.Pointer(ptr))[:len(data)], data)
	globalUnlock.Call(hMem)

	_, _, _ = setClipboardData.Call(uintptr(cfHTML), hMem)
	return true
}

// ✅ 写入图片（base64 PNG → CF_DIB）
func CopyImage(base64PNG string) bool {
	if !strings.HasPrefix(base64PNG, "data:image/png;base64,") {
		return false
	}

	pngData, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(base64PNG, "data:image/png;base64,"))
	if err != nil {
		return false
	}

	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return false
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	type BITMAPINFOHEADER struct {
		Size          uint32
		Width         int32
		Height        int32
		Planes        uint16
		BitCount      uint16
		Compression   uint32
		SizeImage     uint32
		XPelsPerMeter int32
		YPelsPerMeter int32
		ClrUsed       uint32
		ClrImportant  uint32
	}

	header := BITMAPINFOHEADER{
		Size:        40,
		Width:       int32(width),
		Height:      int32(height),
		Planes:      1,
		BitCount:    32,
		Compression: 0,
		SizeImage:   uint32(width * height * 4),
	}

	headerSize := int(unsafe.Sizeof(header))
	pixelSize := width * height * 4
	totalSize := headerSize + pixelSize

	hMem, _, _ := globalAlloc.Call(GMEM_MOVEABLE, uintptr(totalSize))
	if hMem == 0 {
		return false
	}
	ptr, _, _ := globalLock.Call(hMem)
	if ptr == 0 {
		return false
	}
	defer globalUnlock.Call(hMem)

	buf := (*[1 << 30]byte)(unsafe.Pointer(ptr))
	headerBytes := (*[40]byte)(unsafe.Pointer(&header))
	copy(buf[:headerSize], headerBytes[:])

	offset := headerSize
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buf[offset+0] = byte(b >> 8)
			buf[offset+1] = byte(g >> 8)
			buf[offset+2] = byte(r >> 8)
			buf[offset+3] = byte(a >> 8)
			offset += 4
		}
	}

	if err := open(); err != nil {
		return false
	}
	defer close()

	emptyClipboard.Call()
	_, _, _ = setClipboardData.Call(CF_DIB, hMem)

	return true
}

// 热键
//type HotKeyCallback func()

// type MSG struct {
// 	HWND    uintptr
// 	Message uint32
// 	WParam  uintptr
// 	LParam  uintptr
// 	Time    uint32
// 	Pt      struct {
// 		X int32
// 		Y int32
// 	}
// }

// 启动热键监听线程
// func init() {
// 	//go hotkeyListener()
// 	go startListener()
// }

// func hotkeyListener() {
// 	var msg MSG
// 	for {
// 		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
// 		if int32(ret) == -1 {
// 			break
// 		}
// 		if msg.Message == WM_HOTKEY {
// 			id := uint32(msg.WParam)
// 			mutex.RLock()
// 			cb, ok := callbacks[id]
// 			mutex.RUnlock()
// 			if ok {
// 				fmt.Printf("✅ 热键触发 (id=%d)\n", id)
// 				cb()
// 			} else {
// 				fmt.Printf("❌ 无回调 (id=%d)\n", id)
// 			}
// 		}
// 		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
// 		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
// 	}
// }

// RegisterHotKeyWithCallback 注册热键并绑定回调
// func RegisterHotKeyWithCallback(keyCode, modifiers uint32, cb HotKeyCallback) {
// 	id := hotkeyID
// 	hotkeyID++

// 	mutex.Lock()
// 	callbacks[id] = cb
// 	mutex.Unlock()

// 	r, _, err := procRegisterHotKey.Call(0, uintptr(id), uintptr(modifiers), uintptr(keyCode))
// 	if r == 0 {
// 		fmt.Printf("❌ 热键注册失败: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("✈️ 注册热键 (id=%d, key=%d, mod=%d)\n", id, keyCode, modifiers)
// }

// // UnregisterHotKey 注销热键
// func UnregisterHotKey(keyCode, modifiers uint32) {
// 	for id := range callbacks {
// 		_, _, _ = procUnregisterHotKey.Call(0, uintptr(id))
// 	}
// 	mutex.Lock()
// 	callbacks = map[uint32]HotKeyCallback{}
// 	mutex.Unlock()
// 	fmt.Println("🧹 所有热键已注销")
// }

const (
	SW_RESTORE = 9
)

func SetForegroundWindow(hwnd uintptr) {
	procSetForeground.Call(hwnd)
}

func ShowWindow(hwnd uintptr, nCmdShow int) {
	procShowWindow.Call(hwnd, uintptr(nCmdShow))
}

func GetForegroundWindow() uintptr {
	hwnd, _, _ := procGetForeground.Call()
	return hwnd
}

func ActivateApp() {
	hwnd := GetForegroundWindow() // 或你自己保存的主窗口句柄
	if hwnd != 0 {
		ShowWindow(hwnd, SW_RESTORE)
		SetForegroundWindow(hwnd)
	}
}

// 改用"golang.design/x/hotkey" 调用热键，直接调用windows API会无法监听热键消息
type HotKeyCallback func()

type hotkeyEntry struct {
	hk  *hotkey.Hotkey
	key uint32
	mod uint32
}

var (
	callbacks  = make(map[int]hotkeyEntry) //make(map[int]*hotkey.Hotkey)
	cbHandlers = make(map[int]HotKeyCallback)
	idCounter  = 1
	mutex      sync.RWMutex
	startOnce  sync.Once
)

// RegisterHotKeyWithCallback 注册热键并绑定回调（使用 golang.design/x/hotkey 实现）
func RegisterHotKeyWithCallback(keyCode, modifiers uint32, cb HotKeyCallback) {
	mod := toHotkeyModifiers(modifiers)
	key := toHotkeyKey(keyCode)

	hk := hotkey.New(mod, key)
	err := hk.Register()
	if err != nil {
		fmt.Printf("❌ 热键注册失败: %v\n", err)
		return
	}

	localID := idCounter
	idCounter++

	mutex.Lock()
	callbacks[localID] = hotkeyEntry{hk: hk, key: keyCode, mod: modifiers} // hk
	cbHandlers[localID] = cb
	mutex.Unlock()

	// 启动监听线程
	go func(id int, h *hotkey.Hotkey, cb HotKeyCallback) {
		for {
			<-h.Keydown()
			fmt.Printf("✅ 热键触发 (id=%d)\n", id)
			cb()
		}
	}(localID, hk, cb)
	// 仅启动一次监听线程
	// 启动统一监听线程（只启动一次）
	//startOnce.Do(startListener)

	fmt.Printf("✈️ 注册热键 (id=%d, key=%v, mod=%v)\n", localID, key, mod)

}

// 启动统一监听线程
// func startListener() {
// 	fmt.Println("🚀 启动热键监听线程")
// 	//itemp:=0
// 	for id, hk := range callbacks {
// 		//itemp++
// 		cb := cbHandlers[id]
// 		if cb == nil {
// 			go func(id int, h *hotkey.Hotkey) {
// 				for {
// 					select {
// 					case <-h.Keydown():
// 						mutex.RLock()
// 						if cb, ok := cbHandlers[id]; ok {
// 							fmt.Printf("✅ 热键触发 (id=%d)\n", id)
// 							cb() // 调用回调
// 						} else {
// 							fmt.Printf("❌ 无回调 (id=%d)\n", id)
// 						}
// 						mutex.RUnlock()
// 					}
// 				}
// 			}(id, hk)
// 		}
// 	}
// 	// go func() {
// 	// 	fmt.Println("🚀 热键监听线程启动")
// 	// 	for {
// 	// 		mutex.RLock()
// 	// 		for id, hk := range callbacks {
// 	// 			select {
// 	// 			case <-hk.Keydown():
// 	// 				if cb, ok := cbHandlers[id]; ok {
// 	// 					go cb() // 非阻塞调用
// 	// 				}
// 	// 			default:
// 	// 			}
// 	// 		}
// 	// 		mutex.RUnlock()
// 	// 	}
// 	// }()
// }

// UnregisterHotKey 注销所有已注册热键
func UnregisterHotKey(_keyCode, _modifiers uint32) {
	mutex.Lock()
	defer mutex.Unlock()

	for id, hk := range callbacks {
		if _keyCode == hk.key && _modifiers == hk.mod {
			_ = hk.hk.Unregister() // 注销热键
			delete(callbacks, id)  // 删除回调
			//delete(cbHandlers, id) // 删除回调处理器
			fmt.Printf("🧹 注销热键 (id=%d, key=%v ", id, _keyCode)
		}
		//_ = hk.Unregister()
		//fmt.Printf(" 注销热键 (id=%d)\n", id)
	}

	// callbacks = make(map[int]*hotkey.Hotkey)
	// cbHandlers = make(map[int]HotKeyCallback)

	// idCounter = 1
	//startOnce = sync.Once{} // 重置 once，以便下次可以重新启动监听
}

func toHotkeyModifiers(modifiers uint32) []hotkey.Modifier {
	var mods []hotkey.Modifier
	if modifiers&1 != 0 {
		mods = append(mods, hotkey.ModAlt)
	}
	if modifiers&2 != 0 {
		mods = append(mods, hotkey.ModCtrl)
	}
	if modifiers&4 != 0 {
		mods = append(mods, hotkey.ModShift)
	}
	if modifiers&8 != 0 {
		mods = append(mods, hotkey.ModWin)
	}
	return mods
}

func toHotkeyKey(keyCode uint32) hotkey.Key {
	return hotkey.Key(keyCode)
}

func SimulateCmdC() {

}

//

func HideDock() {
	// Windows 不需要隐藏 Dock 图标
	// 但可以在这里添加其他逻辑
	fmt.Println("Windows 平台不支持隐藏 Dock 图标")
}
func CheckAccessibilityPermission() bool {
	fmt.Println("Windows 不需要检查辅助功能权限")
	return true // Windows 不需要检查辅助功能权限
}
func TriggerAccessibilityPrompt() bool {
	return true // Windows 不需要检查辅助功能权限
}

func RecognizeTextFromImageMac(imagePath string) (string, error) {
	return imagePath, nil
}
func RecognizeImageBase64(base64str string) (string, error) {
	return base64str, nil
}

func RecognizeImageAndCopyToClipboard(base64str string) error {
	return nil
}

// ----开机自启
// 注册表路径
const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

// EnableAutoStart 启用开机自启 2
func EnableAutoStart(appName string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	return k.SetStringValue(appName, exePath)
}

// DisableAutoStart 禁用开机自启
func DisableAutoStart(appName string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	// 删除注册表键值
	return k.DeleteValue(appName)
}

// IsEnabled 检查是否启用开机自启
func IsEnabled(appName string) bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(appName)
	return err == nil
}
