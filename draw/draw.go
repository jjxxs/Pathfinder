package draw

import (
	"fmt"
	"image"
	_ "image/png"
	"leistungsnachweis-graphiker/problem"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
)

type Drawer struct {
	imagePath   string
	imageBounds image.Rectangle
	drawerPath  string
}

func NewDrawer(imagePath, drawerPath string) (Drawer, error) {
	if _, err := os.Stat(imagePath); err != nil {
		return Drawer{}, err
	}

	if _, err := os.Stat(drawerPath); err != nil {
		return Drawer{}, err
	}

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return Drawer{imagePath: imagePath, drawerPath: drawerPath, imageBounds: img.Bounds()}, nil
}

func (d *Drawer) Draw(problem problem.Problem) {
	fmt.Printf("drawing..")
	coords := d.mapRouteToImageCoordinates(
		problem.ShortestRoute,
		problem.Image.X1,
		problem.Image.Y1,
		problem.Image.X2,
		problem.Image.Y2)
	points := fmt.Sprintf("%v", coords)
	points = strings.ReplaceAll(points, " ", ",")
	points = strings.ReplaceAll(points, "[", "")
	points = strings.ReplaceAll(points, "]", "")
	//args := fmt.Sprintf("--picture_in=\"%s\" --picture_out=\"%s\" --points=\"%v\"", problem.Image.Path, "/home/octav/blub.png", points)
	cmd := exec.Command("/home/octav/PycharmProjects/TSPic/TSPic.py",
		fmt.Sprintf("--picture_in=%s", problem.Image.Path),
		fmt.Sprintf("--picture_out=%s", "/home/octav/blub.png"),
		fmt.Sprintf("--points=%s", points))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Drawer) mapRouteToImageCoordinates(route problem.Route, x1, y1, x2, y2 float64) []int {
	coordinates := make([]int, 2*len(route))

	xDiff := math.Abs(x1 - x2)
	yDiff := math.Abs(y1 - y2)

	xPixel := float64(d.imageBounds.Max.X) / xDiff
	yPixel := float64(d.imageBounds.Max.Y) / yDiff

	for i, point := range route {
		x := (point.X - x1) * xPixel
		y := (point.Y - y2) * yPixel
		coordinates[i*2] = int(x)
		coordinates[i*2+1] = int(y)
	}

	return coordinates
}
