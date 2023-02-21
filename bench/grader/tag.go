package grader

import (
	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/score"
)

var (
	ScorePostMember score.ScoreTag = "post_member"
	ScoreGetMembers score.ScoreTag = "get_members"
	ScoreGetMember  score.ScoreTag = "get_member"

	ScorePostBooks score.ScoreTag = "post_books"
)

func setScore(result *isucandar.BenchmarkResult) {
	result.Score.Set(ScorePostMember, 20)
	result.Score.Set(ScoreGetMembers, 20)
	result.Score.Set(ScoreGetMember, 1)

	result.Score.Set(ScorePostBooks, 20)
}
