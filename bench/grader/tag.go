package grader

import (
	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/score"
)

var (
	ScoreGetMembers score.ScoreTag = "get_members"
	ScoreGetMember  score.ScoreTag = "get_member"
)

func setScore(result *isucandar.BenchmarkResult) {
	result.Score.Set(ScoreGetMembers, 5)
	result.Score.Set(ScoreGetMember, 1)
}
