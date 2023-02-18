package model

import "github.com/isucon/isucandar/failure"

var (
	// 予期しないクリティカルなエラー
	ErrCritical failure.StringCode = "critical"

	// リクエスト失敗
	ErrRequestFailed failure.StringCode = "request_failed"
)
