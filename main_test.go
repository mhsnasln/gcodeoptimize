package main

import "testing"

func TestOptimize(t *testing.T) {

	var x float64 = 12

	var sample_quiz = []string{
		";FLAVOR:Marlin",
		";TIME:4347",
		";Filament used: 19.6817m",
		"G1 Z15.0 F6000 ;Move the platform down 15mm",
	}
	var answers = []string{
		"",
		"",
		"",
		"G1 Z15.0 F6000  the platform down 15mm",
	}

	for index, question := range sample_quiz {
		optimize, err := worker(question, x)

		if err != nil {
			panic("Birşeyler ters gitti.")
		}

		if optimize != answers[index] {
			t.Errorf(optimize)
		}

	}
}
