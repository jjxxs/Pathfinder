package problem

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// represents a 'tsp-problem' that is to be solved by the solver
type Problem struct {
	// info about the problemset
	Info Info `json:"Info"`

	// set of points on the route
	Points []Point `json:"points"`

	// adjacency matrix, e.g. distances between the points
	Adjacency [][]float32 `json:"adjacency"`
}

// loads a set of problems from a directory
func FromDir(dir string) ([]Problem, error) {
	// stat directory to test if it's accessible
	if file, err := os.Stat(dir); err != nil {
		return []Problem{}, err
	} else if !file.IsDir() {

		// if it's not a directory try to load it as a file
		problem, err := FromFile(dir)
		if err != nil {
			return []Problem{}, err
		}

		return []Problem{problem}, nil
	}

	// try to read directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []Problem{}, err
	}

	// try to load every file in directory as problem
	problems := make([]Problem, 0)
	for _, file := range files {

		// is subdirectory, ignore
		if file.IsDir() {
			continue
		}

		// skip files that don't end with "json"
		absFilePath := filepath.Join(dir, file.Name())
		if filepath.Ext(absFilePath) != ".json" {
			continue
		}

		// load problem
		problem, err := FromFile(absFilePath)
		if err != nil {
			continue
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

// loads a problem from a file
func FromFile(file string) (Problem, error) {
	// stat file to test if it's accessible
	if file, err := os.Stat(file); err != nil {
		return Problem{}, err
	} else if file.IsDir() {
		return Problem{}, errors.New("expected file but provided directory")
	}

	// read file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return Problem{}, err
	}

	// parse json to problem
	var problem = Problem{}
	err = json.Unmarshal(bytes, &problem)
	if err != nil {
		return Problem{}, err
	}

	// calculate adjacency and return
	problem.calculateAdjacency()
	return problem, nil
}

// creates a new problem from given points and info
func NewProblem(points []Point, info Info) *Problem {
	p := Problem{Points: points, Info: info}
	p.calculateAdjacency()
	return &p
}

// contains information about a problem
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // either 'geographic' or 'euclidean'
}

type Point struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Name string  `json:"name"`
}

type Cycle []Point
type Cycles []Cycle

func (p *Problem) String() string {
	return p.Info.Name
}

// calculates the adjacency matrix of the problem with given points
// uses the haversine-formula to calculate distances for "geographic" problems
// uses euclidean distance for "euclidean" problems
func (p *Problem) calculateAdjacency() {
	var calcDistance func(p1, p2 Point) float32

	switch pType := strings.ToLower(p.Info.Type); pType {
	case "geographic":
		calcDistance = haversine
	case "euclidean":
		calcDistance = euclidean
	default:
		calcDistance = euclidean
	}

	// allocate adjacency and calculate distances
	p.Adjacency = make([][]float32, len(p.Points))
	for i, rowPoint := range p.Points {

		adjRow := make([]float32, len(p.Points))
		for j, colPoint := range p.Points {
			adjRow[j] = calcDistance(rowPoint, colPoint)
		}

		p.Adjacency[i] = adjRow
	}
}

// the earths radius in kilometer, used to calculate distances on spheres using the haversine formula
const EarthRadius = 6371

// calculates the shortest point between two points located on a sphere (the earth)
func haversine(p1, p2 Point) float32 {
	deg2rad := func(deg float32) float64 { return (math.Pi * float64(deg)) / 180 }

	lat1 := deg2rad(p1.X)
	lat2 := deg2rad(p2.X)
	long1 := deg2rad(p1.Y)
	long2 := deg2rad(p2.Y)

	deltaLong := long1 - long2
	deltaLat := lat1 - lat2

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return float32(EarthRadius * c)
}

// calculates the shortest distance between two points in an euclidean system
func euclidean(p1, p2 Point) float32 {
	deltaX := math.Abs(float64(p1.X - p2.X))
	deltaY := math.Abs(float64(p1.Y - p2.Y))
	return float32(math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2)))
}
