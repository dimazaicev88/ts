package filters

import "github.com/guregu/null/v6"

type Worker struct {
	Status null.Value[uint8]
	Name   null.Value[string]
	Skip   int
	Limit  int
}
