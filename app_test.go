package main

import (
	"fmt"
	"testing"
)

// first map challenge should be correctly populated
func TestFirstChallenge(t *testing.T) {
	megaverse := createMegaverse(11)

	solveFirstChallenge(&megaverse)
	megaverse.print()

	expected := CreateMegaverseFromFile(GoalFile1)

	compare(&expected, &megaverse, t)
}

func TestAllAstroTypesArePresentInTheMap(t *testing.T) {
	megaverse := CreateMegaverseFromFile(GoalFile2)
	megaverse.print()
	checkContent(megaverse.Content, 2, 2, astralObjectIdFor(Polyanet), t)
	checkContent(megaverse.Content, 2, 13, astralObjectIdFor(ComethUp), t)
	checkContent(megaverse.Content, 5, 12, astralObjectIdFor(ComethDown), t)
	checkContent(megaverse.Content, 4, 16, astralObjectIdFor(ComethLeft), t)
	checkContent(megaverse.Content, 1, 7, astralObjectIdFor(ComethRight), t)
	checkContent(megaverse.Content, 7, 15, astralObjectIdFor(SoloonRed), t)
	checkContent(megaverse.Content, 3, 20, astralObjectIdFor(SoloonWhite), t)
	checkContent(megaverse.Content, 5, 19, astralObjectIdFor(SoloonBlue), t)
	checkContent(megaverse.Content, 4, 8, astralObjectIdFor(SoloonPurple), t)
}

// check content of a single position; stop execution on first difference
func checkContent(content *[][]int, r int, c int, e int, t *testing.T) {
	if e != (*content)[r][c] {
		t.Fatalf("incorrect (%d,%d) expectted %d but got %d", r, c, e, (*content)[r][c])
	}
}

// utils

// print map to console
func print(megaverse Megaverse) {
	for i, a := range *megaverse.Content {
		fmt.Print(i)
		fmt.Print("  ")
		for j, _ := range a {
			fmt.Print((*megaverse.Content)[i][j])
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

// compare maps; stop execution on first difference
func compare(m1 *Megaverse, m2 *Megaverse, t *testing.T) {
	for row, columns := range *m1.Content {
		for column, _ := range columns {
			if (*m1.Content)[row][column] != (*m2.Content)[row][column] {
				t.Fatalf("incorrect (%d,%d) expected %d == %d", row, column, (*m1.Content)[row][column], (*m2.Content)[row][column])
			}
		}
	}
}
