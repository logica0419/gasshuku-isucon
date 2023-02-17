package utils

import "math/rand"

type Choice struct {
	Value  any
	Weight int
}

func WeightedSelect(choices []Choice) any {
	total := 0
	for _, choice := range choices {
		total += choice.Weight
	}

	r := rand.Intn(total)
	for _, choice := range choices {
		if r < choice.Weight {
			return choice.Value
		}
		r -= choice.Weight
	}

	return nil
}
