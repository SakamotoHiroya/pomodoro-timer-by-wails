package jsonstore

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"pomodoro-timer-by-wails/internal/paths"
	"pomodoro-timer-by-wails/models"
)

const PomodoroSettingsFileName = "settings.json"

func LoadPomodoroSettings(ctx context.Context) (models.PomodoroSettings, error) {
	dir, err := paths.AppDataPath("pomodoro-timer-by-wails")
	if err != nil {
		return models.PomodoroSettings{}, err
	}
	path := filepath.Join(dir, PomodoroSettingsFileName)

	fileData, err := os.ReadFile(path)
	if err != nil {
		SavePomodoroSettings(ctx, models.DefaultPomodoroSettings())
		fileData, err = os.ReadFile(path)
		if err != nil {
			return models.PomodoroSettings{}, err
		}
	}

	var data models.PomodoroSettings
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return models.PomodoroSettings{}, err
	}

	return data, nil
}

func SavePomodoroSettings(ctx context.Context, data models.PomodoroSettings) error {
	dir, err := paths.AppDataPath("pomodoro-timer-by-wails")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, PomodoroSettingsFileName)

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

func IsExistPomodoroSettings(ctx context.Context) (bool, error) {
	dir, err := paths.AppDataPath("pomodoro-timer-by-wails")
	if err != nil {
		return false, err
	}
	path := filepath.Join(dir, PomodoroSettingsFileName)

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
