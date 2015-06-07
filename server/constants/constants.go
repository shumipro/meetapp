package constants

import (
	"sort"

	"github.com/mattn/natural"
)

type Constants map[string][]Constant

type Constant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ConstantSlice []Constant

func (p ConstantSlice) Len() int           { return len(p) }
func (p ConstantSlice) Less(i, j int) bool { return natural.NaturalCaseComp(p[i].ID, p[j].ID) < 0 }
func (p ConstantSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ConstantSlice) Sort()              { sort.Sort(p) }

var allConstants = Constants{}

func addConstants(key string, c []Constant) {
	allConstants[key] = c
}

func AllConstants() Constants {
	return allConstants
}

func convertConstant(m map[string]string) []Constant {
	c := []Constant{}
	for key, val := range m {
		c = append(c, Constant{key, val})
	}

	sort.Sort(ConstantSlice(c))
	return c
}
