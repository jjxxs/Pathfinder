package problem

import (
	"math"
	"testing"
)

const (
	TestProblemDir           = "../samples/"
	TestProblemFileGermany   = "../samples/germany13.json"
	TestProblemFileWorkpiece = "../samples/workpiece.json"
)

func TestProblemLoadDir(t *testing.T) {
	// try to load directory
	problems, err := FromDir(TestProblemDir)
	if err != nil {
		t.Fatalf("failed to load problems from dir=%s, err=%s", TestProblemDir, err)
	}

	if len(problems) != 2 {
		t.Fatalf("failed to load problems from dir, invalid length")
	}
}

func TestProblemFileLoadGermany(t *testing.T) {
	// try to load problem
	problem, err := FromFile(TestProblemFileGermany)
	if err != nil {
		t.Fatalf("failed to load problem from file=%s, err=%s", TestProblemFileGermany, err)
	}

	// test Info
	if problem.Info.Name != "Germany 13" ||
		problem.Info.Description != "The thirteen biggest german cities by population" ||
		problem.Info.Type != Geographic {
		t.Fatalf("failed to load problem Info")
	}

	// test points length
	if len(problem.Points) != 13 {
		t.Fatalf("failed to load problem-points, invalid length")
	}

	// verify points
	expectedPoints := []Point{
		{X: 13.40514, Y: 52.5246, Name: "Berlin"},
		{X: 9.994583, Y: 53.5544, Name: "Hamburg"},
		{X: 11.5755, Y: 48.1374, Name: "München"},
		{X: 6.95000, Y: 50.9333, Name: "Köln"},
		{X: 8.68333, Y: 50.1167, Name: "Frankfurt"},
		{X: 9.1770, Y: 48.7823, Name: "Stuttgart"},
		{X: 6.8121, Y: 51.2205, Name: "Düsseldorf"},
		{X: 7.4660, Y: 51.5149, Name: "Dortmund"},
		{X: 7.0086, Y: 51.4624, Name: "Essen"},
		{X: 12.3713, Y: 51.3396, Name: "Leipzig"},
		{X: 8.8077, Y: 53.07516, Name: "Bremen"},
		{X: 13.7500, Y: 51.0500, Name: "Dresden"},
		{X: 9.7332, Y: 52.3705, Name: "Hannover"},
	}

	for _, actualPoint := range problem.Points {
		contains := false
		for _, expectedPoint := range expectedPoints {
			if expectedPoint.X == actualPoint.X ||
				expectedPoint.Y == actualPoint.Y ||
				expectedPoint.Name == actualPoint.Name {
				contains = true
			}
		}
		if !contains {
			t.Fatalf("failed to load problem-points, invalid Point-details")
		}
	}
}

func TestProblemFileLoadWerkstueck(t *testing.T) {
	// try to load problem
	problem, err := FromFile(TestProblemFileWorkpiece)
	if err != nil {
		t.Fatalf("failed to load problem from file=%s, err=%s", TestProblemFileWorkpiece, err)
	}

	// test Info
	if problem.Info.Name != "Workpiece" ||
		problem.Info.Description != "Technical drawing of a flat workpiece with prisms" ||
		problem.Info.Type != Euclidean {
		t.Fatalf("failed to load problem Info")
	}

	// test points length
	if len(problem.Points) != 30 {
		t.Fatalf("failed to load problem-points, invalid length")
	}

	// verify points
	expectedPoints := []Point{
		{X: 230, Y: 138},
		{X: 195, Y: 197},
		{X: 157, Y: 198},
		{X: 157, Y: 298},
		{X: 187, Y: 328},
		{X: 157, Y: 357},
		{X: 157, Y: 550},
		{X: 218, Y: 611},
		{X: 309, Y: 611},
		{X: 357, Y: 611},
		{X: 514, Y: 611},
		{X: 278, Y: 555},
		{X: 389, Y: 555},
		{X: 513, Y: 537},
		{X: 559, Y: 537},
		{X: 559, Y: 138},
		{X: 309, Y: 207},
		{X: 350, Y: 274},
		{X: 270, Y: 274},
		{X: 328, Y: 432},
		{X: 328, Y: 450},
		{X: 308, Y: 450},
		{X: 275, Y: 475},
		{X: 239, Y: 441},
		{X: 276, Y: 406},
		{X: 309, Y: 433},
		{X: 456, Y: 312},
		{X: 456, Y: 235},
		{X: 417, Y: 273},
		{X: 494, Y: 273},
	}

	for _, actualPoint := range problem.Points {
		contains := false
		for _, expectedPoint := range expectedPoints {
			if expectedPoint.X == actualPoint.X ||
				expectedPoint.Y == actualPoint.Y ||
				expectedPoint.Name == actualPoint.Name {
				contains = true
			}
		}
		if !contains {
			t.Fatalf("failed to load problem-points, invalid Point-details")
		}
	}
}

func TestProblemFromPointsEuclidean(t *testing.T) {
	points := []Point{
		{X: 13.23, Y: 52.31, Name: "Berlin"},
		{X: 10.0, Y: 53.33, Name: "Hamburg"},
		{X: 11.34, Y: 48.8, Name: "Munich"},
		{X: 6.57, Y: 50.56, Name: "Cologne"},
	}

	info := Info{
		Type: Euclidean,
	}

	cartesianProblem := Problem{Points: points, Info: info}
	cartesianProblem.calculateAdjacency()

	expectedAdj := [][]float64{
		{0, 3.387226003679116, 6.886080162182256, 3.9865022262630228},
		{3.387226003679116, 0, 4.408832044884447, 4.7240342928475885},
		{6.886080162182256, 4.408832044884447, 0, 5.0843386983953005},
		{3.9865022262630228, 4.7240342928475885, 5.0843386983953005, 0},
	}

	for i, row := range expectedAdj {
		for j := range row {
			if math.Round(cartesianProblem.Adjacency[i][j]*100)/100 != math.Round(expectedAdj[i][j]*100)/100 {
				t.Fatalf("failed to load euclidean problem")
			}
		}
	}
}

func TestProblemFromPointsGeographic(t *testing.T) {
	points := []Point{
		{X: 13.23, Y: 52.31, Name: "Berlin"},
		{X: 10.0, Y: 53.33, Name: "Hamburg"},
		{X: 11.34, Y: 48.8, Name: "Munich"},
		{X: 6.57, Y: 50.56, Name: "Cologne"},
	}

	info := Info{
		Type: Geographic,
	}

	geographicProblem := Problem{Points: points, Info: info}
	geographicProblem.calculateAdjacency()

	expectedAdj := [][]float64{
		{0, 375.94456892271836, 764.9363495938279, 435.41080570770396},
		{375.94456892271836, 0, 488.19591094433196, 516.9268227447614},
		{764.9363495938279, 488.19591094433196, 0, 564.5107479834015},
		{435.41080570770396, 516.9268227447614, 564.5107479834015, 0},
	}

	for i, row := range expectedAdj {
		for j := range row {
			if math.Round(geographicProblem.Adjacency[i][j]*100)/100 != math.Round(expectedAdj[i][j]*100)/100 {
				t.Fatalf("failed to load geographic problem: %f  == %f", geographicProblem.Adjacency[i][j], expectedAdj[i][j])
			}
		}
	}
}
