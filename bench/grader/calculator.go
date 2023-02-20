package grader

import (
	"log"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

func CulcResult(result *isucandar.BenchmarkResult, finish bool) bool {
	passed := true
	reason := "pass"
	errors := result.Errors.All()

	setScore(result)
	scoreRaw := result.Score.Sum()

	for tag, count := range result.Score.Breakdown() {
		logger.Admin.Printf("SCORE: %s: %d", tag, count)
	}

	errCount := int64(0)
	timeoutCount := int64(0)
	for _, err := range errors {
		switch {
		case model.IsErrCritical(err):
			passed = false
			reason = "fail: critical"
		case model.IsErrTimeout(err):
			timeoutCount++
		default:
			errCount += 1
		}
	}
	deductionTotal := errCount*10 + timeoutCount/10

	score := scoreRaw - deductionTotal
	if score <= 0 && passed {
		passed = false
		reason = "fail: score"
	}

	var scoreLogger *log.Logger
	if finish {
		scoreLogger = logger.Contestant
	} else {
		scoreLogger = logger.Admin
	}

	scoreLogger.Printf("score: %d(%d - %d) : %s", score, scoreRaw, deductionTotal, reason)
	scoreLogger.Printf("deduction: %d / timeout: %d", errCount, timeoutCount)

	return passed
}
