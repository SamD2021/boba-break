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
