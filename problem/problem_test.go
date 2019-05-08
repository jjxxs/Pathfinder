package problem

import (
	"reflect"
	"testing"
)

const (
	TestProblemDir           = "../samples/"
	TestProblemFileGermany   = "../samples/germany10.json"
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
	if problem.Info.Name != "Germany" ||
		problem.Info.Description != "Ten biggest german cities" ||
		problem.Info.Type != "geographic" {
		t.Fatalf("failed to load problem Info")
	}

	// test points length
	if len(problem.Points) != 10 {
		t.Fatalf("failed to load problem-points, invalid length")
	}

	// verify points
	expectedPoints := []Point{
		{X: 13.23, Y: 52.31, Name: "Berlin"},
		{X: 10.0, Y: 53.33, Name: "Hamburg"},
		{X: 11.34, Y: 48.8, Name: "Munich"},
		{X: 6.57, Y: 50.56, Name: "Cologne"},
		{X: 8.41, Y: 50.7, Name: "Frankfurt"},
		{X: 9.11, Y: 48.47, Name: "Stuttgart"},
		{X: 6.47, Y: 51.14, Name: "Düsseldorf"},
		{X: 7.28, Y: 51.31, Name: "Dortmund"},
		{X: 7.1, Y: 51.27, Name: "Essen"},
		{X: 12.23, Y: 51.2, Name: "Leipzig"},
	}

	for i, actualPoint := range problem.Points {
		expectedPoint := expectedPoints[i]
		if expectedPoint.X != actualPoint.X ||
			expectedPoint.Y != actualPoint.Y ||
			expectedPoint.Name != actualPoint.Name {
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
		problem.Info.Type != "euclidean" {
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

	for i, actualPoint := range problem.Points {
		expectedPoint := expectedPoints[i]
		if expectedPoint.X != actualPoint.X ||
			expectedPoint.Y != actualPoint.Y ||
			expectedPoint.Name != actualPoint.Name {
			t.Fatalf("failed to load problem-points, invalid Point-details")
		}
	}
}

func TestProblemFromPointsCartesian(t *testing.T) {
	points := []Point{
		{X: 13.23, Y: 52.31, Name: "Berlin"},
		{X: 10.0, Y: 53.33, Name: "Hamburg"},
		{X: 11.34, Y: 48.8, Name: "Munich"},
		{X: 6.57, Y: 50.56, Name: "Cologne"},
	}

	info := Info{
		Type: "euclidean",
	}

	cartesianProblem := NewProblem(points, info)

	expectedAdj := [][]float32{
		{0, 3.3872256, 3.9865038, 6.88608},
		{3.3872256, 0, 4.7240367, 4.408832},
		{3.9865038, 4.7240367, 0, 5.0843396},
		{6.88608, 4.408832, 5.0843396, 0},
	}

	if !reflect.DeepEqual(expectedAdj, cartesianProblem.Adjacency) {
		t.Fatalf("failed to load euclidean problem")
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
		Type: "geographic",
	}

	geographicProblem := NewProblem(points, info)

	expectedAdj := [][]float32{
		{0, 375.94455, 435.41098, 764.9363},
		{375.94455, 0, 516.9271, 488.19592},
		{435.41098, 516.9271, 0, 564.5108},
		{764.9363, 488.19592, 564.5108, 0},
	}

	if !reflect.DeepEqual(expectedAdj, geographicProblem.Adjacency) {
		t.Fatalf("failed to load geographic problem")
	}
}
