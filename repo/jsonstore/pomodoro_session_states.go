package jsonstore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"pomodoro-timer-by-wails/internal/paths"
	"pomodoro-timer-by-wails/models"
)

const PomodoroSessionStatesFileName = "session_states.json"

func LoadPomodoroSessionStates(ctx context.Context) (models.SessionState, error) {
	dir, err := paths.AppDataPath("pomodoro-timer-by-wails")
	if err != nil {
		return models.SessionState{}, err
	}
	path := filepath.Join(dir, PomodoroSessionStatesFileName)

	fileData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// ファイルが存在しない場合は空の状態を返す
			return models.SessionState{}, nil
		}
		return models.SessionState{}, err
	}

	var data models.SessionState
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return models.SessionState{}, err
	}

	return data, nil
}

func SavePomodoroSessionStates(ctx context.Context, data models.SessionState) error {
	dir, err := paths.AppDataPath("pomodoro-timer-by-wails")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, PomodoroSessionStatesFileName)

	fileData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, fileData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func IsExistPomodoroSessionStates(ctx context.Context) (bool, error) {
	dir, err := paths.AppDataPath("pomodoro-timer-by-wails")
	if err != nil {
		return false, err
	}
	path := filepath.Join(dir, PomodoroSessionStatesFileName)

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
