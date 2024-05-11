/*
 * Copyright (c) 2024 Samuel Dasilva
 *
 * This file is part of Boba Break.
 *
 * Boba Break is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Boba Break is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Boba Break. If not, see <https://www.gnu.org/licenses/>.
 */
package breaklog

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	NewLogEntry() BreakLogEntry
	BreakLogEntry
}

type FileBreakLogger struct {
	filepath string
	entries  []BreakLogEntry
}

func checkFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		_, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewFileBreakLogger(filepath string) (*FileBreakLogger, error) {
	entries := []BreakLogEntry{}
	err := checkFile(filepath)
	if err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() != 0 {
		file, err := os.ReadFile(filepath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(file, &entries)
		if err != nil {
			return nil, err
		}
	}

	return &FileBreakLogger{
		filepath: filepath,
		entries:  entries,
	}, nil
}

func NewBreakLogEntry(finding string, Reason string) *BreakLogEntry {
	return &BreakLogEntry{
		Timestamp: time.Now(),
		Reason:    Reason,
		Findings:  finding,
	}
}

func (f *FileBreakLogger) SetFilePath(fp string) {
	f.filepath = fp
}

func (f *FileBreakLogger) AddLogEntry(entry *BreakLogEntry) error {
	file, err := os.OpenFile(f.filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// Encode if necessary
	// _, err := file.WriteString()
	f.entries = append(f.entries, *entry)
	json_encoded, err := json.Marshal(f.entries)
	if err != nil {
		return err
	}
	_, err = file.Write(json_encoded)
	return err
}
