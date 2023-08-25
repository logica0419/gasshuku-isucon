package utils

import (
	"errors"
	"math/rand"
)

// 重み付き選択肢
type Choice[V any] struct {
	Val    V
	Weight int
}

// 重み付きランダム選択
//
//	decがtrueの場合、選択された選択肢の重みを1減らす
func WeightedSelect[V any](choices []Choice[V], dec bool) (V, error) {
	total := 0
	for _, choice := range choices {
		total += choice.Weight
	}

	if total <= 0 {
		return choices[0].Val, errors.New("total weight is zero")
	}

	r := rand.Intn(total)
	for i, choice := range choices {
		if choice.Weight <= 0 {
			choice.Weight = 1
		}
		if r < choice.Weight {
			if dec {
				choices[i].Weight--
			}
			return choice.Val, nil
		}
		r -= choice.Weight
	}

	if dec {
		choices[len(choices)-1].Weight--
	}
	return choices[len(choices)-1].Val, nil
}
