/*
Unit tests for function computePoints in main.go
*/

package main

import (
    "testing"
)

// check point total for Target receipt
func Test_computePoints(t *testing.T) {
    
    test_cases := map[string]int {
        "examples/target.json": 28, 
        "examples/m-m-corner-market.json": 109, 
    }
    
    for fileName, points := range(test_cases) {
        got := pointsFromFile(fileName)
        want := points
        if got != want {
            t.Errorf("got %q on %s, wanted %q", got, fileName, want)
        }
    }
}