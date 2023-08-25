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
func WeightedSelect[V any](choices []Choice[V]) (V, error) {
	total := 0
	for _, choice := range choices {
		if choice.Weight <= 0 {
			choice.Weight = 1
		}
		total += choice.Weight
	}

	if total <= 0 {
		return choices[0].Val, errors.New("total weight is zero")
	}

	r := rand.Intn(total)
	for _, choice := range choices {
		if choice.Weight <= 0 {
			choice.Weight = 1
		}
		if r < choice.Weight {
			return choice.Val, nil
		}
		r -= choice.Weight
	}

	return choices[len(choices)-1].Val, nil
}
