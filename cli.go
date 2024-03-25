package main

import (
	"fmt"
	"golang.org/x/term"
	"math"
	"strconv"
)

func printLine(fileName string, proress float64) {
	if !term.IsTerminal(0) {
		return
	}

	width, _, err := term.GetSize(0)
	if width < 35 {
		width = 35
	}

	if err != nil {
		return
	}

	if len(fileName) > width/3*2 {
		fileName = "..." + fileName[len(fileName)-(width/3*2-3):]
	}

	progress := getBarBySize(width/3, proress, 0)

	spaces := width - len(progress) - len(fileName)

	result := fileName
	for i := 0; i < spaces-1; i++ {
		result += " "
	}
	result += progress
	fmt.Print(result)
}

func getBarBySize(size int, progress float64, roundPrec int) string {
	if size == 0 {
		return ""
	}

	zeroBar := progressBar(progress, 0, roundPrec)
	if len(zeroBar) > size {
		return "."
	}
	return progressBar(progress, size-len(zeroBar), roundPrec)
}

func progressBar(progress float64, cellNumber int, roundPrec int) string {
	num := int(math.Round(progress / 100 * float64(cellNumber)))

	res := "["
	for i := 0; i < num; i++ {
		res += "#"
	}
	for i := 0; i < cellNumber-num; i++ {
		res += "-"
	}
	res += fmt.Sprintf("] %."+strconv.Itoa(roundPrec)+"f%% ", progress)

	return res
}

func cleanLine() error {
	_, err := fmt.Print("\033[2K\r")
	return err
}

// [#####     ] 50%
