package grader

import (
	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/score"
)

var ScoreGetMembes score.ScoreTag = "get_members"

func setScore(result *isucandar.BenchmarkResult) {
	result.Score.Set(ScoreGetMembes, 1)
}
