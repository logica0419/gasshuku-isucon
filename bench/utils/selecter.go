package utils

import (
	"math/rand"
)

type Choice[V any] struct {
	Val    V
	Weight int
}

func WeightedSelect[V any](choices []Choice[V]) V {
	total := 0
	for _, choice := range choices {
		total += choice.Weight
	}

	r := rand.Intn(total)
	for _, choice := range choices {
		if r < choice.Weight {
			return choice.Val
		}
		r -= choice.Weight
	}

	return choices[len(choices)-1].Val
}
