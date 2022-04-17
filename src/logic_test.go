package main

import (
	"testing"
)

func TestCheckCorrectSolution(t *testing.T) {
	good := "145327698839654127672918543496185372218473956753296481367542819984761235521839764"
	ret, err := Check(good)
	if ret != true {
		t.Fatalf(`Check(%s) = %t, %v, want true, error`, good, ret, err)
	}
}

func TestCheckIncorrectSolution(t *testing.T) {
	bad := "INSERT_FLAG_HERE"
	ret, err := Check(bad)
	if ret == true {
		t.Fatalf(`Check(%s) = %t, %v, want false, error`, bad, ret, err)
	}
}
