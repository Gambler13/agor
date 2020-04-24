package main

import "math"

// Quadtree - The quadtree data structure
type Quadtree struct {
	Bounds     Bounds
	MaxObjects int // Maximum objects a node can hold before splitting into 4 subnodes
	MaxLevels  int // Total max levels inside root Quadtree
	Level      int // Depth level, required for subnodes
	Objects    []EntityImpl
	Nodes      []Quadtree
	Total      int
}

// Intersects - Checks if a Bounds object intersects with another Bounds
func (b *Entity) Intersects(a Entity) bool {

	posA := a.Position
	posB := b.Position

	//Distance between center points
	distance := math.Sqrt(q(posB.X-posA.X) + q(posB.Y-posA.Y))

	if distance-a.Radius < b.Radius {
		// The two overlap
		return true
	}

	return false
}

// Bounds - A bounding box with a x,y origin and width and height
type Bounds struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (b Bounds) intersect(entity Entity) bool {
	return intersects(entity.Circle, Rectangle{
		Position: Position{
			X: b.X,
			Y: b.Y,
		},
		Width:  b.Width,
		Height: b.Height,
	})
}

// TotalNodes - Retrieve the total number of sub-Quadtrees in a Quadtree
func (qt *Quadtree) TotalNodes() int {

	total := 0

	if len(qt.Nodes) > 0 {
		for i := 0; i < len(qt.Nodes); i++ {
			total += 1
			total += qt.Nodes[i].TotalNodes()
		}
	}

	return total

}

// split - split the node into 4 subnodes
func (qt *Quadtree) split() {

	if len(qt.Nodes) == 4 {
		return
	}

	nextLevel := qt.Level + 1
	subWidth := qt.Bounds.Width / 2
	subHeight := qt.Bounds.Height / 2
	x := qt.Bounds.X
	y := qt.Bounds.Y

	//top right node (0)
	qt.Nodes = append(qt.Nodes, Quadtree{
		Bounds: Bounds{
			X:      x + subWidth,
			Y:      y,
			Width:  subWidth,
			Height: subHeight,
		},
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]EntityImpl, 0),
		Nodes:      make([]Quadtree, 0, 4),
	})

	//top left node (1)
	qt.Nodes = append(qt.Nodes, Quadtree{
		Bounds: Bounds{
			X:      x,
			Y:      y,
			Width:  subWidth,
			Height: subHeight,
		},
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]EntityImpl, 0),
		Nodes:      make([]Quadtree, 0, 4),
	})

	//bottom left node (2)
	qt.Nodes = append(qt.Nodes, Quadtree{
		Bounds: Bounds{
			X:      x,
			Y:      y + subHeight,
			Width:  subWidth,
			Height: subHeight,
		},
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]EntityImpl, 0),
		Nodes:      make([]Quadtree, 0, 4),
	})

	//bottom right node (3)
	qt.Nodes = append(qt.Nodes, Quadtree{
		Bounds: Bounds{
			X:      x + subWidth,
			Y:      y + subHeight,
			Width:  subWidth,
			Height: subHeight,
		},
		MaxObjects: qt.MaxObjects,
		MaxLevels:  qt.MaxLevels,
		Level:      nextLevel,
		Objects:    make([]EntityImpl, 0),
		Nodes:      make([]Quadtree, 0, 4),
	})

}

// getIndex - Determine which quadrant the object belongs to (0-3)
func (qt *Quadtree) getIndex(pRect Bounds) int {

	index := -1 // index of the subnode (0-3), or -1 if pRect cannot completely fit within a subnode and is part of the parent node

	verticalMidpoint := qt.Bounds.X + (qt.Bounds.Width / 2)
	horizontalMidpoint := qt.Bounds.Y + (qt.Bounds.Height / 2)

	//pRect can completely fit within the top quadrants
	topQuadrant := (pRect.Y < horizontalMidpoint) && (pRect.Y+pRect.Height < horizontalMidpoint)

	//pRect can completely fit within the bottom quadrants
	bottomQuadrant := (pRect.Y > horizontalMidpoint)

	//pRect can completely fit within the left quadrants
	if (pRect.X < verticalMidpoint) && (pRect.X+pRect.Width < verticalMidpoint) {

		if topQuadrant {
			index = 1
		} else if bottomQuadrant {
			index = 2
		}

	} else if pRect.X > verticalMidpoint {
		//pRect can completely fit within the right quadrants

		if topQuadrant {
			index = 0
		} else if bottomQuadrant {
			index = 3
		}

	}

	return index

}

// Insert - Insert the object into the node. If the node exceeds the capacity,
// it will split and add all objects to their corresponding subnodes.
func (qt *Quadtree) Insert(entity EntityImpl) {

	qt.Total++

	i := 0
	var index int

	// If we have subnodes within the Quadtree
	if len(qt.Nodes) > 0 == true {

		index = qt.getIndex(entity.getEntity().Bounds())

		if index != -1 {
			qt.Nodes[index].Insert(entity)
			return
		}
	}

	// If we don't subnodes within the Quadtree
	qt.Objects = append(qt.Objects, entity)

	// If total objects is greater than max objects and level is less than max levels
	if (len(qt.Objects) > qt.MaxObjects) && (qt.Level < qt.MaxLevels) {

		// split if we don't already have subnodes
		if len(qt.Nodes) > 0 == false {
			qt.split()
		}

		// Add all objects to there corresponding subNodes
		for i < len(qt.Objects) {

			index = qt.getIndex(qt.Objects[i].getEntity().Bounds())

			if index != -1 {

				splice := qt.Objects[i]                                  // Get the object out of the slice
				qt.Objects = append(qt.Objects[:i], qt.Objects[i+1:]...) // Remove the object from the slice

				qt.Nodes[index].Insert(splice)

			} else {

				i++

			}

		}

	}

}

// Retrieve - Return all objects that could collide with the given object
func (qt *Quadtree) Retrieve(rect Bounds) []EntityImpl {

	index := qt.getIndex(rect)

	// Array with all detected objects
	returnObjects := qt.Objects

	//if we have subnodes ...
	if len(qt.Nodes) > 0 {

		//if pRect fits into a subnode ..
		if index != -1 {

			returnObjects = append(returnObjects, qt.Nodes[index].Retrieve(rect)...)

		} else {

			//if pRect does not fit into a subnode, check it against all subnodes
			for i := 0; i < len(qt.Nodes); i++ {
				returnObjects = append(returnObjects, qt.Nodes[i].Retrieve(rect)...)
			}

		}
	}

	return returnObjects

}

// RetrieveIntersections - Bring back all the bounds in a Quadtree that intersect with a provided bounds
func (qt *Quadtree) RetrieveIntersections(find EntityImpl) []EntityImpl {

	var foundIntersections []EntityImpl

	potentials := qt.Retrieve(find.getEntity().Bounds())
	for o := 0; o < len(potentials); o++ {
		e := potentials[o].getEntity()
		if e.Id == find.getEntity().Id {
			continue
		}
		if e.Intersects(*find.getEntity()) {
			foundIntersections = append(foundIntersections, potentials[o])
		}
	}

	return foundIntersections

}

// RetrieveIntersections - Bring back all the bounds in a Quadtree that intersect with a provided bounds
func (qt *Quadtree) RetrieveViewIntersections(find Bounds) []EntityImpl {

	var foundIntersections []EntityImpl

	potentials := qt.Retrieve(find)
	for o := 0; o < len(potentials); o++ {
		e := potentials[o].getEntity()
		if find.intersect(*e) {
			foundIntersections = append(foundIntersections, potentials[o])
		}
	}

	return foundIntersections

}

//Clear - Clear the Quadtree
func (qt *Quadtree) Clear() {

	qt.Objects = []EntityImpl{}

	if len(qt.Nodes)-1 > 0 {
		for i := 0; i < len(qt.Nodes); i++ {
			qt.Nodes[i].Clear()
		}
	}

	qt.Nodes = []Quadtree{}
	qt.Total = 0

}
