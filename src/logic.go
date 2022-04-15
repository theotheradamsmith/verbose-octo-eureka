package main

import (
	"errors"
	"fmt"
	"strings"
)

const digits string = "123456789"
const rows string = "ABCDEFGHI"
const cols string = digits

var squares = cross(rows, cols)
var unitList = createUnitList(cols, rows)
var units = createUnits(squares, unitList)

//var peers = createPeers(units)

func Hello(name string) (string, error) {
	if name == "" {
		return name, errors.New("empty name")
	}

	message := fmt.Sprintf("Hello, %v! Thanks for stopping by to gander at the logic!", name)
	return message, nil
}

func cross(a, b string) []string {
	var ret []string
	for _, av := range a {
		for _, bv := range b {
			ret = append(ret, string(av)+string(bv))
		}
	}
	return ret
}

func createUnitList(cols, rows string) [][]string {
	ret := make([][]string, len(rows)*3)
	i := 0
	for _, c := range cols {
		ret[i] = cross(rows, string(c))
		i++
	}
	for _, r := range rows {
		ret[i] = cross(string(r), cols)
		i++
	}
	rs := []string{"ABC", "DEF", "GHI"}
	cs := []string{"123", "456", "789"}
	for _, r := range rs {
		for _, c := range cs {
			ret[i] = cross(r, string(c))
			i++
		}
	}
	return ret
}

func createUnits(squares []string, unitList [][]string) map[string][][]string {
	units := make(map[string][][]string, len(squares))
	for _, s := range squares {
		unit := make([][]string, 3)
		i := 0
		for _, u := range unitList {
			for _, su := range u {
				if s == su {
					unit[i] = u
					i++
					break
				}
			}
		}
		units[s] = unit
	}
	return units
}

func GridValues(grid string) (map[string]string, error) {
	// Convert grid into a dict of {square: char} with '0' or '.' for empties
	gridValues := make(map[string]string, len(squares))
	validChars := make([]string, 0, len(grid))

	for _, c := range grid {
		if strings.Contains(digits, string(c)) || strings.Contains(".0", string(c)) {
			validChars = append(validChars, string(c))
		}
	}

	if len(validChars) != 81 {
		return gridValues, errors.New("Invalid input grid")
	}

	for i, s := range squares {
		gridValues[s] = string(validChars[i])
	}

	return gridValues, nil
}

func Verify(values map[string]string) (bool, error) {
	unitSolved := func(unit []string) bool {
		digitsSet := make(map[string]bool, len(digits))
		for _, r := range digits {
			digitsSet[string(r)] = true
		}
		for _, s := range unit {
			key := string(values[s])
			if _, ok := digitsSet[key]; ok {
				delete(digitsSet, key)
			} else {
				return false
			}
		}
		return len(digitsSet) == 0
	}
	for _, unit := range unitList {
		if unitSolved(unit) != true {
			msg := fmt.Sprintf("Error in unit: [%s:%s]", unit[0], unit[8])
			return false, errors.New(msg)
		}
	}
	return true, nil
}

func Check(userInput string) (bool, error) {
	if gv, ok := GridValues(userInput); ok != nil {
		return false, ok
	} else {
		if _, ok := Verify(gv); ok != nil {
			return false, ok
		}
	}
	return true, nil
}
