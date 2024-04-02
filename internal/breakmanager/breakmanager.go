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
package breakmanager

import (
	"fmt"
	"time"
)

type breakEntry struct {
	timer    *time.Timer
	Duration time.Duration
}

type BreakManager interface {
	StartBreak()
	StopBreak()
}

type breaks struct {
	Breaks   []breakEntry
	Sessions uint
}

func (bs breaks) New(durStr string) (*breakEntry, error) {
	duration, err := time.ParseDuration(durStr)
	if err != nil {
		fmt.Println("ERROR in startCMD:", err)
		return nil, err
	}
	entry := breakEntry{
		Duration: duration,
	}
	return &entry, nil
}

func (b *breakEntry) StartBreak() {
	after := func() {
		fmt.Println("Break is done")
	}
	b.timer = time.AfterFunc(b.Duration, after)
	fmt.Println("Break Started")
	<-b.timer.C
}

func (b *breakEntry) StopBreak() {
	if !b.timer.Stop() {
		<-b.timer.C
	}
}

func (b *breakEntry) ResetBreak(d time.Duration) {
	b.timer.Reset(d)
}
