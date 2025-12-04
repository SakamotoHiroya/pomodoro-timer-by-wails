package models

import "time"

type SessionMode string

const (
    ModeWork       SessionMode = "work"
    ModeShortBreak SessionMode = "short_break"
    ModeLongBreak  SessionMode = "long_break"
)

type SessionState struct {
    Mode                    SessionMode `json:"mode"`           // 現在のモード（作業／短休憩／長休憩）
	CurrentModeStartedAt time.Time   `json:"current_session_started_at"` // 現在のセッション開始時刻
    StartedAt               time.Time   `json:"started_at"`     // セッション開始時刻
    Paused                  bool        `json:"paused"`         // 一時停止中かどうか
    SessionCount            int         `json:"session_count"`  // 今日の通算セッション数（作業のみカウント）
}
