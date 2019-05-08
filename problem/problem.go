package problem

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type Problem struct {
	Info      info        `json:"info"`
	Points    []point     `json:"points"`
	Adjacency [][]float64 `json:"adjacency"`
}

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

		// load problem
		problem, err := FromFile(dir)
		if err != nil {
			continue
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

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

type info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type point struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Name string  `json:"name"`
}

type Cycle []point
type Cycles []Cycle

func (p *Problem) String() string {
	return p.Info.Name
}

func (p *Problem) calculateAdjacency() {
	var calcDistance func(p1, p2 point) float64

	switch pType := strings.ToLower(p.Info.Type); pType {
	case "geographic":
		calcDistance = haversine
	case "cartesian":
		calcDistance = cartesian
	default:
		calcDistance = cartesian
	}

	// allocate adjacency and calculate distances
	p.Adjacency = make([][]float64, len(p.Points))
	for i, rowPoint := range p.Points {

		adjRow := make([]float64, len(p.Points))
		for j, colPoint := range p.Points {
			adjRow[j] = calcDistance(rowPoint, colPoint)
		}

		p.Adjacency[i] = adjRow
	}
}

// the earths radius in kilometer, used to calculate distances on spheres using the haversine formula
const EarthRadius = 6371

func haversine(p1, p2 point) float64 {
	deg2rad := func(deg float64) float64 { return (math.Pi * deg) / 180 }

	lat1 := deg2rad(p1.X)
	lat2 := deg2rad(p2.X)
	long1 := deg2rad(p1.Y)
	long2 := deg2rad(p2.Y)

	deltaLong := long1 - long2
	deltaLat := lat1 - lat2

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

func cartesian(p1, p2 point) float64 {
	deltaX := math.Abs(p1.X - p2.X)
	deltaY := math.Abs(p1.Y - p2.Y)
	return math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2))
}
