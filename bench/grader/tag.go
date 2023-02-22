package grader

import (
	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/score"
)

var (
	ScorePostMember score.ScoreTag = "post_member"
	ScoreGetMembers score.ScoreTag = "get_members"
	ScoreBanMember  score.ScoreTag = "ban_member"

	ScorePostBooks   score.ScoreTag = "post_books"
	ScoreSearchBooks score.ScoreTag = "search_books"

	ScorePostLendings   score.ScoreTag = "post_lendings"
	ScoreGetLendings    score.ScoreTag = "get_lendings"
	ScoreReturnLendings score.ScoreTag = "return_lendings"
)

func setScore(result *isucandar.BenchmarkResult) {
	result.Score.Set(ScorePostMember, 20)
	result.Score.Set(ScoreGetMembers, 15)
	result.Score.Set(ScoreBanMember, 1)

	result.Score.Set(ScorePostBooks, 20)
	result.Score.Set(ScoreSearchBooks, 15)

	result.Score.Set(ScorePostLendings, 25)
	result.Score.Set(ScoreGetLendings, 15)
	result.Score.Set(ScoreReturnLendings, 25)
}
