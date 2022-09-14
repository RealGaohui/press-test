package utils

import "time"

type Alert interface {
	SendWechat(message string) error
}

type Result struct {
	Backfill time.Duration
	FP Resource
	DB Resource
	WRK wrk
}

type wrk struct {
	Threads  int
	Connections int
}
