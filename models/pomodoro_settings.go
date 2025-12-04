package models

type PomodoroSettings struct {
    WorkMinutes         int  `json:"work_minutes"`          // 作業時間（分）
    ShortBreakMinutes   int  `json:"short_break_minutes"`   // 短い休憩時間（分）
    LongBreakMinutes    int  `json:"long_break_minutes"`    // 長い休憩時間（分）
    LongBreakInterval   int  `json:"long_break_interval"`   // 何セッションごとに長い休憩を取るか
    AutoStartNext       bool `json:"auto_start_next"`       // セッション終了後に自動で次を開始するか
}

func DefaultPomodoroSettings() PomodoroSettings {
	return PomodoroSettings{
		WorkMinutes:         25,
		ShortBreakMinutes:   5,
		LongBreakMinutes:    15,
		LongBreakInterval:   4,
		AutoStartNext:       true,
	}
}
