package problem

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Problem struct {
	Info   info    `json:"info"`
	Points []point `json:"points"`
}

func FromDir(dir string) ([]Problem, error) {
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

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []Problem{}, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		log.Println(file.Name())
	}

	return []Problem{}, nil
}

func FromFile(file string) (Problem, error) {
	if file, err := os.Stat(file); err != nil {
		return Problem{}, err
	} else if file.IsDir() {
		return Problem{}, errors.New("expected file but provided directory")
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return Problem{}, err
	}

	var problem = Problem{}
	err = json.Unmarshal(bytes, &problem)
	if err != nil {
		return Problem{}, err
	}

	return problem, nil
}

type info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type point struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Name string  `json:"name"`
}
