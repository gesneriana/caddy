package model

import "time"

// ShellConfig web应用的shell配置
type ShellConfig struct {
	Domain           string    `json:"domain"`
	InitShell        string    `json:"init_shell"`
	IsInit           bool      `json:"is_init"`
	SyncShell        string    `json:"sync_shell"`
	Interval         int       `json:"interval"`
	LastSyncTime     time.Time `json:"last_sync_time"`
	StartShell       string    `json:"start_shell"`
	VerificationCode string    `json:"verification_code"`
}
