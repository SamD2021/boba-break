package breaklog

import (
	"encoding/json"
	"os"
	"time"
)

type BreakLogEntry struct {
	Timestamp      time.Time     `json:"timestamp"`
	Reason         string        `json:"reason,omitempty"` // Optional field
	WorkInProgress string        `json:"work_in_progress"`
	Findings       string        `json:"findings"`
	Duration       time.Duration `json:"duration"`
}

type BreakLogger interface {
	AddLogEntry(entry BreakLogEntry) error
}

type FileBreakLogger struct {
	filepath string
}

func (f *FileBreakLogger) SetFilePath(fp string) {
	f.filepath = fp
}

func (f *FileBreakLogger) AddLogEntry(entry BreakLogEntry) error {
	file, err := os.OpenFile(f.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// Encode if necessary
	// _, err := file.WriteString()
	json_encoded, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = file.Write(json_encoded)
	return err
}
