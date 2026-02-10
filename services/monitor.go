package services

import (
	"changeme/platform"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TipInfoEmit struct {
	Type    string
	Content string
}

var appInstance *application.App
var scheduler *PromptScheduler

func Init(app *application.App, schedulers *PromptScheduler) {
	appInstance = app
	scheduler = schedulers
	//scheduler.Recalculate()
}

// 新增提示后触发提示事件给展示页面
func EmitTipsEvent(contentType string) {
	if contentType == "immediate" { //立即类型，插入后直接触发提示
		latestTip := scheduler.PromGetLatest()
		appInstance.Event.Emit("tipEvent", TipInfo{
			ID:       latestTip.ID,
			Category: latestTip.Category,
			Type:     latestTip.Type,
			Title:    latestTip.Title,
			State:    latestTip.State,
			Pinned:   latestTip.Pinned,
			PinnedAt: latestTip.PinnedAt,
			SnoozeAt: latestTip.SnoozeAt,
			ExpireAt: latestTip.ExpireAt,
		})
	} else { //定时类型，重新计算调度
		scheduler.Recalculate()
	}
}

// 提示调度器
func (s *PromptScheduler) trigger(tip *TipInfo) {
	fmt.Printf("⏰ 触发提示 ID=%d, 内容=%s\n", tip.ID, tip.Title)
	// 1️⃣ 通知前端展示
	//appInstance.Event.Emit("tip:trigger", tip)
	appInstance.Event.Emit("tipEvent", TipInfo{
		ID:       tip.ID,
		Category: tip.Category,
		Type:     tip.Type,
		Title:    tip.Title,
		State:    2, // 已展示
		Pinned:   tip.Pinned,
		PinnedAt: tip.PinnedAt,
		SnoozeAt: tip.SnoozeAt,
		ExpireAt: tip.ExpireAt,
	})

	// 2️⃣ 更新状态
	_ = s.repo.UpTipsState(tip.ID, 2) // 设置为已展示
	fmt.Println("✅ 提示状态已更新为已展示")
	// 3️⃣ 继续调度下一条
	s.Recalculate()
}

type hotkeyEntry struct {
	ID       int
	Hotkey   Hotkey
	Window   application.Window
	Callback func()
}

var (
	hotkeys = make(map[int]*hotkeyEntry)
	mutex   sync.RWMutex
)

func RegisterHotkeyEntry(id int, keycode, modifiers uint32, win application.Window, cb func()) {
	mutex.Lock()
	defer mutex.Unlock()

	hotkeys[id] = &hotkeyEntry{
		ID: id,
		Hotkey: Hotkey{
			KeyCode:   keycode,
			Modifiers: modifiers,
		},
		Window:   win,
		Callback: cb,
	}
	fmt.Printf("✅ 注册热键 Entry: id=%d (%d, %d)\n", id, keycode, modifiers)
	platform.RegisterHotKeyWithCallback(keycode, modifiers, cb)
}

func RegistertargetHotkeyEntry(id int, keycode, modifiers uint32, cb func()) {
	mutex.Lock()
	defer mutex.Unlock()

	hotkeys[id] = &hotkeyEntry{
		ID: id,
		Hotkey: Hotkey{
			KeyCode:   keycode,
			Modifiers: modifiers,
		},
		Window:   nil, // 没有窗口绑定
		Callback: cb,
	}
	fmt.Printf("✅ 注册热键 Entry: id=%d (%d, %d)\n", id, keycode, modifiers)
	platform.RegisterHotKeyWithCallback(keycode, modifiers, cb)
}

func LoadAndRegisterHotkeysFrom(m *SuiStore, id int) {
	var hk Hotkey
	err := m.db.QueryRow("SELECT keycode, modifiers FROM hotkeys  where id = ?", id).Scan(&hk.KeyCode, &hk.Modifiers)
	if err != nil {
		fmt.Println("Error querying hotkeys:", err)
		return
	}
	mutex.RLock()
	entry, ok := hotkeys[id]
	mutex.RUnlock()

	if !ok {
		fmt.Println("❌ 未找到之前注册的窗口或回调，无法重新注册热键")
		return
	}
	platform.UnregisterHotKey(hk.KeyCode, hk.Modifiers) // 先注销旧的热键
	RegisterHotkeyEntry(id, hk.KeyCode, hk.Modifiers, entry.Window, entry.Callback)
}

func RegisterWindowAndCallback(id int, win application.Window, cb platform.HotKeyCallback) {
	mutex.Lock()
	hotkeys[id] = &hotkeyEntry{
		ID: id,
		Hotkey: Hotkey{
			KeyCode:   0, // 默认值，实际使用时不需要
			Modifiers: 0, // 默认值，实际使用时不需要
		},
		Window:   win,
		Callback: cb,
	}
	mutex.Unlock()
}
