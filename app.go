package main

import (
	"context"
	"errors"
	"pomodoro-timer-by-wails/models"
	"pomodoro-timer-by-wails/repo/jsonstore"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) GetPomodoroSettings() (models.PomodoroSettings, error) {
	settings, err := jsonstore.LoadPomodoroSettings(a.ctx)
	if err != nil {
		return models.PomodoroSettings{}, err
	}
	return settings, nil
}

func (a *App) GetPomodoroSessionStates() (models.SessionState, error) {
	sessionStates, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return models.SessionState{}, err
	}
	return sessionStates, nil
}

func (a *App) StartPomodoro() error {
	now := time.Now()
	sessionStates := models.SessionState{
		Mode:                 models.ModeWork,
		CurrentModeStartedAt: now,
		StartedAt:            now,
		Paused:               false,
		SessionCount:         0,
	}

	err := jsonstore.SavePomodoroSessionStates(a.ctx, sessionStates)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetCurrentSessionState() (models.SessionState, error) {
	sessionStates, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return models.SessionState{}, err
	}
	return sessionStates, nil
}

func (a *App) GetCurrentTime() (time.Duration, error) {
	err := a.UpdateSessionState()
	if err != nil {
		return 0, err
	}
	sessionStates, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return 0, err
	}
	return time.Since(sessionStates.StartedAt), nil
}

func (a *App) PausePomodoro() error {
	sessionStates, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return err
	}
	sessionStates.Paused = true
	err = jsonstore.SavePomodoroSessionStates(a.ctx, sessionStates)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) UpdateSessionState() error {
	currentState, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return err
	}

	// セッションが開始されていない場合は何もしない
	if currentState.CurrentModeStartedAt.IsZero() {
		return nil
	}

	settings, err := jsonstore.LoadPomodoroSettings(a.ctx)
	if err != nil {
		return err
	}

	// 現在のModeの持続時間を取得
	var currentModeDuration time.Duration
	switch currentState.Mode {
	case models.ModeWork:
		currentModeDuration = time.Duration(settings.WorkMinutes) * time.Minute
	case models.ModeShortBreak:
		currentModeDuration = time.Duration(settings.ShortBreakMinutes) * time.Minute
	case models.ModeLongBreak:
		currentModeDuration = time.Duration(settings.LongBreakMinutes) * time.Minute
	default:
		return nil
	}

	// 現在のModeの終了時間を計算
	currentModeEndTime := currentState.CurrentModeStartedAt.Add(currentModeDuration)
	now := time.Now()

	// 現在時刻が終了時間を過ぎている場合、次のステートに更新
	if now.After(currentModeEndTime) || now.Equal(currentModeEndTime) {
		nextState, err := a.GetNextSessionState()
		if err != nil {
			return err
		}
		err = jsonstore.SavePomodoroSessionStates(a.ctx, nextState)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) ResumePomodoro() error {
	sessionStates, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return err
	}
	sessionStates.Paused = false
	// 再開時は現在時刻をCurrentModeStartedAtに設定（一時停止時間を考慮）
	sessionStates.CurrentModeStartedAt = time.Now()
	err = jsonstore.SavePomodoroSessionStates(a.ctx, sessionStates)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetRemainingTime() (time.Duration, error) {
	err := a.UpdateSessionState()
	if err != nil {
		return 0, err
	}

	currentState, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return 0, err
	}

	// セッションが開始されていない場合は0を返す
	if currentState.CurrentModeStartedAt.IsZero() {
		return 0, nil
	}

	settings, err := jsonstore.LoadPomodoroSettings(a.ctx)
	if err != nil {
		return 0, err
	}

	// 現在のModeの持続時間を取得
	var currentModeDuration time.Duration
	switch currentState.Mode {
	case models.ModeWork:
		currentModeDuration = time.Duration(settings.WorkMinutes) * time.Minute
	case models.ModeShortBreak:
		currentModeDuration = time.Duration(settings.ShortBreakMinutes) * time.Minute
	case models.ModeLongBreak:
		currentModeDuration = time.Duration(settings.LongBreakMinutes) * time.Minute
	default:
		return 0, errors.New("invalid mode")
	}

	// 現在のModeの終了時間を計算
	currentModeEndTime := currentState.CurrentModeStartedAt.Add(currentModeDuration)
	now := time.Now()

	// 残り時間を計算
	remainingTime := currentModeEndTime.Sub(now)
	if remainingTime < 0 {
		return 0, nil
	}

	return remainingTime, nil
}

func (a *App) ResetPomodoro() {
	jsonstore.SavePomodoroSessionStates(a.ctx, models.SessionState{})
}

func (a *App) GetNextSessionState() (models.SessionState, error) {
	currentState, err := jsonstore.LoadPomodoroSessionStates(a.ctx)
	if err != nil {
		return models.SessionState{}, err
	}

	settings, err := jsonstore.LoadPomodoroSettings(a.ctx)
	if err != nil {
		return models.SessionState{}, err
	}

	// 現在のModeの持続時間を取得
	var currentModeDuration time.Duration
	switch currentState.Mode {
	case models.ModeWork:
		currentModeDuration = time.Duration(settings.WorkMinutes) * time.Minute
	case models.ModeShortBreak:
		currentModeDuration = time.Duration(settings.ShortBreakMinutes) * time.Minute
	case models.ModeLongBreak:
		currentModeDuration = time.Duration(settings.LongBreakMinutes) * time.Minute
	default:
		return models.SessionState{}, errors.New("invalid mode")
	}

	// 現在のModeの終了時間を計算
	currentModeEndTime := currentState.CurrentModeStartedAt.Add(currentModeDuration)

	// 次のModeを決定
	var nextMode models.SessionMode
	nextSessionCount := currentState.SessionCount

	if currentState.Mode == models.ModeWork {
		// 作業セッションが終了した場合
		nextSessionCount++
		// LongBreakIntervalの倍数かどうかで判定
		if nextSessionCount%settings.LongBreakInterval == 0 {
			nextMode = models.ModeLongBreak
		} else {
			nextMode = models.ModeShortBreak
		}
	} else {
		// 休憩セッションが終了した場合、次の作業セッションへ
		nextMode = models.ModeWork
	}

	// 次のステートを作成
	nextState := models.SessionState{
		Mode:                 nextMode,
		CurrentModeStartedAt: currentModeEndTime,
		StartedAt:            currentState.StartedAt,
		Paused:               false,
		SessionCount:         nextSessionCount,
	}

	return nextState, nil
}
