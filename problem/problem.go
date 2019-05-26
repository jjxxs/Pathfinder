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

// a cycle is a set of integers/point-ids that are to be mapped to points
// this way they represent a light-weight representation of a route that is
// more performant to deal with
type Cycle []int
type Cycles []Cycle

// a route is a set of points in a specific order
// they are the high-level representation of a cycle that was produced by an algorithm
type Route []Point
type Routes []Route

func (r Route) String() string {
	routeStr := ""
	for i, p := range r {
		if i == len(r)-1 {
			routeStr += p.Name
		} else {
			routeStr += p.Name + " <-> "
		}
	}
	return routeStr
}

func (r Routes) String() string {
	routesStr := ""
	for i, route := range r {
		if i == len(r)-1 {
			routesStr += route.String()
		} else {
			routesStr += route.String() + " :: "
		}
	}
	return routesStr
}

// contains the distances between each point on a route
type Adjacency [][]float64

// contains information about a problem
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	// either 'geographic' or 'euclidean'
	// determines how distance between two points is calculated
	Type string `json:"type"`
}

type Point struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Name string  `json:"name"`
}

// represents a 'tsp-problem' that is to be solved by the solver
type Problem struct {
	// info about the problem
	Info Info `json:"Info"`

	// route of the problem, e.g. a set of points that the solver has to bring into the right order
	// for it to be the shortest route possible
	Route Route `json:"points"`

	// adjacency matrix, e.g. distances between the points
	Adjacency Adjacency `json:"adjacency"`
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
func NewProblem(route Route, info Info) *Problem {
	p := Problem{Route: route, Info: info}
	p.calculateAdjacency()
	return &p
}

// converts given cycles to routes
func (p *Problem) GetRoutesFromCycles(cycles Cycles) Routes {
	routes := make(Routes, len(cycles))

	for i, cycle := range cycles {
		ps := make(Route, len(cycle))

		// find minimum in cycle
		minIndex, minValue := 0, -1
		for j := range cycle {
			if minValue == -1 {
				minValue = cycle[j]
			}

			if minValue > cycle[j] {
				minIndex, minValue = j, cycle[j]
			}
		}

		// direction
		minLeft, minRight := minIndex-1, minIndex+1
		if minLeft < 0 {
			minLeft = len(cycle) - 1
		}
		if minRight > len(cycle)-1 {
			minRight = 0
		}

		walkRight := cycle[minRight] < cycle[minLeft]

		// start cycle with minimum-index
		orderedCycle := make(Cycle, len(cycle))
		for j := 0; j < len(cycle); j++ {
			if walkRight {
				orderedCycle[j] = cycle[(minIndex+j)%len(cycle)]
			} else {
				i := minIndex - j
				if i < 0 {
					i = len(cycle) - 1 + i
				}
				orderedCycle[j] = cycle[i]
			}
		}

		// create route
		for j, id := range orderedCycle {
			ps[j] = p.Route[id]
		}

		routes[i] = ps
	}

	return routes
}

func (p *Problem) GetDistancesFromCycles(cycles Cycles) []float64 {
	distances := make([]float64, len(cycles))
	for j, cycle := range cycles {

		var distance float64
		for k, i := range cycle {
			if k == len(cycle)-1 {
				distance += p.Adjacency[i][0]
			} else {
				distance += p.Adjacency[i][i+1]
			}
		}

		distances[j] = distance
	}

	return distances
}

// name of the problem
func (p Problem) String() string {
	return p.Info.Name
}

// calculates the adjacency matrix of the problem with given points
// uses the haversine-formula to calculate distances for "geographic" problems
// uses euclidean distance for "euclidean" problems
func (p *Problem) calculateAdjacency() {
	var calcDistance func(p1, p2 Point) float64

	switch pType := strings.ToLower(p.Info.Type); pType {
	case "geographic":
		calcDistance = haversine
	case "euclidean":
		calcDistance = euclidean
	default:
		calcDistance = euclidean
	}

	// allocate adjacency and calculate distances
	p.Adjacency = make(Adjacency, len(p.Route))
	for i, rowPoint := range p.Route {

		adjRow := make([]float64, len(p.Route))
		for j, colPoint := range p.Route {
			adjRow[j] = calcDistance(rowPoint, colPoint)
		}

		p.Adjacency[i] = adjRow
	}
}

// the earths radius in kilometer, used to calculate distances on spheres using the haversine formula
const EarthRadius = 6371

// calculates the shortest point between two points located on a sphere (the earth)
func haversine(p1, p2 Point) float64 {
	deg2rad := func(deg float32) float64 { return (math.Pi * float64(deg)) / 180 }

	lat1 := deg2rad(p1.X)
	lat2 := deg2rad(p2.X)
	long1 := deg2rad(p1.Y)
	long2 := deg2rad(p2.Y)

	deltaLong := long1 - long2
	deltaLat := lat1 - lat2

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return float64(EarthRadius * c)
}

// calculates the shortest distance between two points in an euclidean system
func euclidean(p1, p2 Point) float64 {
	deltaX := math.Abs(float64(p1.X - p2.X))
	deltaY := math.Abs(float64(p1.Y - p2.Y))
	return float64(math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2)))
}
