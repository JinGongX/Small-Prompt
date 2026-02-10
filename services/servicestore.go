package services

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type SuiStore struct {
	db *sql.DB
}

func getSafeDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, "SuiTips")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, "suitips.db"), nil
}
func NewSuiStore() (*SuiStore, error) {
	dbPath, err := getSafeDBPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建表
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS hotkeys (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    keycode INTEGER NOT NULL,
	    modifiers INTEGER NOT NULL,
	    description TEXT,
		target TEXT 
	);

	CREATE TABLE IF NOT EXISTS tipinfo (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    type TEXT NOT NULL,
		category TEXT NOT NULL,
	    title TEXT NOT NULL,
		desc TEXT,
		createdat INTEGER DEFAULT (strftime('%s','now')),
	    expireat INTEGER NOT NULL,
		pinned INTEGER NOT NULL DEFAULT 0,
  		pinnedat INTEGER,
  		snoozeat INTEGER,
		state INTEGER DEFAULT 1,
		snoozecount INTEGER DEFAULT 0
	);
	`)
	if err != nil {
		return nil, err
	}
	// 检查是否已存在数据
	row := db.QueryRow(`SELECT COUNT(*) FROM hotkeys`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		var sqlhotkeys = OSinithotkeys()
		_, err = db.Exec(sqlhotkeys)
		if err != nil {
			return nil, err
		}
	}

	return &SuiStore{db: db}, nil
}

func (cs *SuiStore) Close() {
	cs.db.Close()
}

func (cs *SuiStore) Start() error {
	// 这里可以初始化数据库或其它启动逻辑
	// 同步过期提示
	_ = cs.SyncExpiredTips()

	return nil
}

func (cs *SuiStore) Stop() error {
	cs.Close()
	return nil
}

type Hotkey struct {
	ID        int    `json:"id"`        // 热键ID
	KeyCode   uint32 `json:"keycode"`   // 键码
	Modifiers uint32 `json:"modifiers"` // 修饰键
}

// 快捷键修改
func (cs *SuiStore) UpHotkey(id int, key int, modifier int) error {
	_, err := cs.db.Exec(`
        UPDATE hotkeys 
        SET keycode = ?, modifiers = ? 
        WHERE id = ?
    `, key, modifier, id)
	return err
}

func (cs *SuiStore) GetHotkeys() ([]Hotkey, error) {
	rows, err := cs.db.Query("SELECT id, keycode, modifiers FROM hotkeys")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotkeys []Hotkey
	for rows.Next() {
		var hk Hotkey
		if err := rows.Scan(&hk.ID, &hk.KeyCode, &hk.Modifiers); err != nil {
			return nil, err
		}
		hotkeys = append(hotkeys, hk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hotkeys, nil
}

func (cs *SuiStore) InsTips(contentType, category, content string, expireAt int64, snoozeAt int64) error {
	_, err := cs.db.Exec(
		"INSERT INTO tipinfo (type,category, title,desc, expireat,snoozeat) VALUES (?, ?, ?,'', ?,?)",
		contentType, category, content, expireAt, snoozeAt,
	)
	if err != nil {
		return err
	}
	EmitTipsEvent(contentType) //触发提示事件
	return nil
}

type TipInfo struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	CreatedAt   int64  `json:"createdAt"`
	ExpireAt    int64  `json:"expireat"`
	State       int    `json:"state"`
	Pinned      int    `json:"pinned"`
	PinnedAt    *int64 `json:"pinnedat"`
	SnoozeAt    *int64 `json:"snoozeat"`
	Snoozecount int    `json:"snoozecount"`
}

func (cs *SuiStore) GetLatestTip() (*TipInfo, error) {
	row := cs.db.QueryRow(`
		SELECT id, type, category, title, desc, createdat, expireat, state, pinned, pinnedat, snoozeat, snoozecount
		FROM tipinfo
		ORDER BY createdat DESC
		LIMIT 1
	`)
	var e TipInfo
	if err := row.Scan(
		&e.ID,
		&e.Type,
		&e.Category,
		&e.Title,
		&e.Desc,
		&e.CreatedAt,
		&e.ExpireAt,
		&e.State,
		&e.Pinned,
		&e.PinnedAt,
		&e.SnoozeAt,
		&e.Snoozecount,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有任何记录
		}
		return nil, err
	}
	return &e, nil
}

func (cs *SuiStore) GetTips() ([]TipInfo, error) {
	rows, err := cs.db.Query("SELECT id, type,category, title, desc, createdat, expireat, state,pinned,pinnedat,snoozeat,snoozecount FROM tipinfo WHERE state IN (1,2,8) ORDER BY createdAt DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tips []TipInfo
	for rows.Next() {
		var e TipInfo
		if err := rows.Scan(&e.ID, &e.Type, &e.Category, &e.Title, &e.Desc, &e.CreatedAt, &e.ExpireAt, &e.State, &e.Pinned, &e.PinnedAt, &e.SnoozeAt, &e.Snoozecount); err != nil {
			return nil, err
		}

		tips = append(tips, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tips, nil
}

// 更新提示状态
func (cs *SuiStore) UpTipsState(id int, state int) error {
	_, err := cs.db.Exec(`
		UPDATE tipinfo
		SET state = ?
		WHERE id = ?
	`, state, id)
	return err
}

func (cs *SuiStore) UpTipsPinned(id int, pinned int) error {
	_, err := cs.db.Exec(`
		UPDATE tipinfo
		SET pinned = ?, pinnedat = ?
		WHERE id = ?
	`, pinned, time.Now().Unix(), id)
	return err
}

// 延迟提示
func (cs *SuiStore) UpTipsDelayed(id int, expireAt int64) error {
	_, err := cs.db.Exec(`
		UPDATE tipinfo
		SET expireat = ?, snoozecount = snoozecount + 1
		WHERE id = ? AND snoozecount < 3 AND state != 8
	`, expireAt, id)
	return err
}

// 同步过期的提示
func (s *SuiStore) SyncExpiredTips() error {
	now := time.Now().Unix()
	_, err := s.db.Exec(`
        UPDATE tipinfo
        SET state = 8 -- 设置为完成
        WHERE type='scheduled' AND state IN (1, 2)
          AND  expireat <= ? 
		  ;
		
		UPDATE tipinfo
        SET state = 2 -- 设置为显示
        WHERE type='scheduled' AND state =1
          AND ( 
		   snoozeat <= ? AND expireat >= ?
          );
    `, now, now, now, now)

	return err
}

func (r *SuiStore) FindNextScheduled(now int64) (*TipInfo, error) {
	row := r.db.QueryRow(`
		SELECT id, title, type, category, state,
		       pinned, pinnedat, snoozeat, expireat
		FROM tipinfo
		WHERE state = 1
		  AND type = 'scheduled'
		  AND COALESCE(snoozeat, expireat) >= ?
		ORDER BY COALESCE(snoozeat, expireat) ASC
		LIMIT 1
	`, now)

	var tip TipInfo
	err := row.Scan(
		&tip.ID,
		&tip.Title,
		&tip.Type,
		&tip.Category,
		&tip.State,
		&tip.Pinned,
		&tip.PinnedAt,
		&tip.SnoozeAt,
		&tip.ExpireAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tip, nil
}
