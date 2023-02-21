package model

import (
	"context"
	"net"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
)

var (
	// ベンチマーク終了後のタイムアウト
	ErrDeadline failure.StringCode = "critical"

	// 予期しないクリティカルなエラー
	ErrCritical failure.StringCode = "critical"

	// リクエスト失敗
	ErrRequestFailed failure.StringCode = "request_failed"
	// ステータスコードが不正
	ErrInvalidStatusCode failure.StringCode = "invalid_status_code"
	// Content-Typeが不正
	ErrInvalidContentType failure.StringCode = "invalid_content_type"
	// レスポンスボディがデコード不可
	ErrUndecodableBody failure.StringCode = "undecodable_body"
	// レスポンスボディが間違っている
	ErrInvalidBody failure.StringCode = "invalid_body"

	// タイムアウト
	ErrTimeout failure.StringCode = "timeout"
)

func IsErrCritical(err error) bool {
	return failure.IsCode(err, isucandar.ErrPrepare) ||
		failure.IsCode(err, ErrCritical)
}

func IsErrTimeout(err error) bool {
	if failure.IsCode(err, ErrTimeout) ||
		failure.IsCode(err, failure.TimeoutErrorCode) {
		return true
	}
	var nErr net.Error
	return failure.As(err, &nErr) && nErr.Timeout()
}

func IsErrCanceled(err error) bool {
	return failure.IsCode(err, ErrDeadline) ||
		failure.IsCode(err, failure.CanceledErrorCode) ||
		failure.Is(err, context.DeadlineExceeded)
}
