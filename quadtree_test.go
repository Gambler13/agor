package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestBounds_intersect(t *testing.T) {
	var food []Food
	foodData, err := ioutil.ReadFile("./food.json")
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(foodData, &food)
	if err != nil {
		t.Error(err)
	}

	b := Bounds{
		X:      0,
		Y:      0,
		Width:  1200,
		Height: 900,
	}

	ft := Quadtree{
		Bounds:     b,
		MaxObjects: 20,
		MaxLevels:  10,
	}

	for i := range food {
		f := food[i]
		ft.Insert(&f)
	}

	p1 := Position{
		X: 377.50526094188143,
		Y: 225.02261706487732,
	}

	p2 := Position{
		X: 377.72807342796136,
		Y: 224.80757710738158,
	}

	e1 := ft.RetrieveViewIntersections(Bounds{
		X:      p1.X,
		Y:      p1.Y,
		Width:  300,
		Height: 300,
	})

	fmt.Printf("\n\n")

	e2 := ft.RetrieveViewIntersections(Bounds{
		X:      p2.X,
		Y:      p2.Y,
		Width:  300,
		Height: 300,
	})

	minX := p1.X - 150
	maxX := p1.X + 150
	minY := p1.Y - 150
	maxY := p1.Y + 150

	counter := 0

	for i := range food {
		p := food[i]
		if p.X > minX && p.X < maxX && p.Y > minY && p.Y < maxY {
			counter++
		}
	}

	fmt.Printf("counter %d\n", counter)

	if len(e1) != len(e2) {
		t.Errorf("lenght %d to %d", len(e1), len(e2))
	}

	for i := range e2 {
		contains := false
		for j := range e1 {
			if e2[i].getEntity().Id == e1[j].getEntity().Id {
				contains = true
			}
		}
		if !contains {
			fmt.Printf("%+v\n", e2[i].getEntity())
		}
	}

}
