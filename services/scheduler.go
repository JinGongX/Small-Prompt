package services

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type PromptScheduler struct {
	repo  *SuiStore
	timer *time.Timer
	mu    sync.Mutex
	ctx   context.Context
}

func NewPromptScheduler(ctx context.Context, repo *SuiStore) *PromptScheduler {
	return &PromptScheduler{
		repo: repo,
		ctx:  ctx,
	}
}

func (s *PromptScheduler) PromGetLatest() *TipInfo {
	GetLatestTip, err := s.repo.GetLatestTip()
	if err != nil {
		fmt.Println("Error getting latest tip:", err)
		return nil
	}
	if GetLatestTip == nil {
		fmt.Println("No tips available to show.")
		return nil
	}
	return GetLatestTip
}

func (s *PromptScheduler) Recalculate() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}

	now := time.Now().Unix()
	fmt.Println("⏰ 重新计算提示调度..."+" 当前时间:", now)
	tip, err := s.repo.FindNextScheduled(now)
	if err != nil || tip == nil {
		fmt.Println("⏰ 没有找到下一个提示，调度结束。", err)
		return
	}
	fmt.Printf("⏰ 下一个提示 ID=%d, 内容=%s, 触发时间=%d\n", tip.ID, tip.Title, tip.ExpireAt)
	triggerAt := tip.ExpireAt
	if tip.SnoozeAt != nil {
		triggerAt = *tip.SnoozeAt
	}

	delay := time.Duration(triggerAt-now) * time.Second
	if delay < 0 {
		delay = 0
	}
	fmt.Printf("⏰ 提示将在 %s 后触发\n", delay.String())
	s.timer = time.AfterFunc(delay, func() {
		s.trigger(tip)
	})
}
