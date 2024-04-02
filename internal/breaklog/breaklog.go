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
