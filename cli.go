package main

import (
	"fmt"
	"math"
	"strconv"
)

func progressBar(progress float64, cellNumber int, roundPrec int) string {
	num := int(math.Round(float64(progress)/100 * float64(cellNumber)))

	res := "["
	for i := 0; i < num; i++ {
		res += "#"
	}
	for i := 0; i < cellNumber-num; i++ {
		res += " "
	}
	res += fmt.Sprintf("] %."+ strconv.Itoa(roundPrec) +"f%% ", progress)

	return res
}
// [#####     ] 50% 