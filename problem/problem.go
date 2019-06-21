package problem

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const (
	Geographic = "geographic"
	Euclidean  = "euclidean"
)

// contains the distances between each point on a route
type Adjacency [][]float64

// a cycle is a set of integers that are to be mapped to points
type Cycle []int

// a route is a set of points in a specific order, it is a high-level representation of a cycle
type Route []Point

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

// represents a 'tsp-problem' that is to be solved by the solver
type Problem struct {
	// info about the problem
	Info Info `json:"info"`

	// path to the image of the problem, if any
	Image Image `json:"image"`

	// route of the problem, e.g. a set of points that the solver has to bring into the right order
	// for it to be the shortest route possible
	Points []Point `json:"points"`

	ShortestRoute Route `json:"route"`

	ShortestDistance float64 `json:"shortestDistance"`

	// adjacency matrix, e.g. distances between the points
	Adjacency Adjacency `json:"adjacency"`
}

// contains information about a problem
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	// either 'geographic' or 'euclidean'
	// determines how distance between two points is calculated
	Type string `json:"type"`
}

type Image struct {
	Path   string  `json:"path"`
	X1     float64 `json:"x1"`
	Y1     float64 `json:"y1"`
	X2     float64 `json:"x2"`
	Y2     float64 `json:"y2"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
}

// a point in two-dimensional space
type Point struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Name string  `json:"name"`
}

type Status struct {
	Algorithm   string  `json:"algorithm"`
	Problem     string  `json:"problem"`
	Description string  `json:"description"`
	Elapsed     string  `json:"elapsed"`
	Shortest    float64 `json:"shortest"`
	Running     bool    `json:"running"`
}

func NewProblem(points []Point) *Problem {
	p := Problem{
		Points: points,
	}
	p.calculateAdjacency()
	return &p
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

func (p *Problem) UpdateRoute(cycle Cycle) {
	// set new route
	route := make(Route, len(cycle))
	for i, j := range cycle {
		route[i] = p.Points[j]
	}
	p.ShortestRoute = route

	// calculate new distance
	var distance float64
	for i := range cycle {
		if i == len(cycle)-1 {
			distance += p.Adjacency[cycle[i]][cycle[0]]
		} else {
			distance += p.Adjacency[cycle[i]][cycle[i+1]]
		}
	}
	p.ShortestDistance = distance
}

// calculates the adjacency matrix of the problem with given points
// 		- uses the haversine-formula to calculate distances for "geographic" problems
// 		- uses euclidean distance for "euclidean" problems
func (p *Problem) calculateAdjacency() {
	var calcDistance func(p1, p2 Point) float64

	switch pType := strings.ToLower(p.Info.Type); pType {
	case Geographic:
		calcDistance = haversine
	case Euclidean:
		calcDistance = euclidean
	default:
		calcDistance = euclidean
	}

	// shuffle before calculating adjacency
	rand.Shuffle(len(p.Points), func(i, j int) { p.Points[i], p.Points[j] = p.Points[j], p.Points[i] })

	// allocate adjacency and calculate distances
	p.Adjacency = make(Adjacency, len(p.Points))
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

// calculates the shortest point between two points located on a sphere (the earth)
func haversine(p1, p2 Point) float64 {
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

// calculates the shortest distance between two points in an euclidean system
func euclidean(p1, p2 Point) float64 {
	deltaX := math.Abs(p1.X - p2.X)
	deltaY := math.Abs(p1.Y - p2.Y)
	return math.Sqrt(math.Pow(deltaX, 2) + math.Pow(deltaY, 2))
}

func (p *Problem) MapRouteToImageCoordinates() []int {
	coordinates := make([]int, 2*len(p.ShortestRoute))

	if p.Info.Type == Geographic {
		xDiff := math.Abs(p.Image.X1 - p.Image.X2)
		yDiff := math.Abs(p.Image.Y1 - p.Image.Y2)

		xPixel := float64(p.Image.Width) / xDiff
		yPixel := float64(p.Image.Height) / yDiff

		for i, point := range p.ShortestRoute {
			x := (point.X - p.Image.X1) * xPixel
			y := (p.Image.Y1 - point.Y) * yPixel
			coordinates[i*2] = int(x)
			coordinates[i*2+1] = int(y)
		}
	} else {
		for i, point := range p.ShortestRoute {
			coordinates[i*2] = int(point.X)
			coordinates[i*2+1] = int(point.Y)
		}
	}

	return coordinates
}
