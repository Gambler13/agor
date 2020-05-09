package main

import "image"

type QuadTree struct {
	Root *Node
}

type Node struct {
	X      int
	Y      int
	Nodes  map[int]*Node
	Entity EntityImpl
}

const NW = 0
const SW = 1
const SE = 2
const NE = 3

func (q *QuadTree) insert(e EntityImpl) {
	q.Root = insert(q.Root, e)
}

func insert(n *Node, e EntityImpl) *Node {
	entity := e.getEntity()
	eX, eY := entity.PosInt()
	if n == nil {
		return &Node{
			X:      eX,
			Y:      eY,
			Nodes:  make(map[int]*Node),
			Entity: e,
		}
	} else if eX < n.X && eY < n.Y {
		n.Nodes[SW] = insert(n.Nodes[SW], e)
	} else if eX < n.X && eY > n.Y {
		n.Nodes[NW] = insert(n.Nodes[NW], e)
	} else if eX > n.X && eY < n.Y {
		n.Nodes[SE] = insert(n.Nodes[SE], e)
	} else if eX > n.X && eY > n.Y {
		n.Nodes[NE] = insert(n.Nodes[NE], e)
	}
	return n
}

func (q *QuadTree) query(rectangle image.Rectangle) []EntityImpl {
	return query(q.Root, rectangle)
}

func query(node *Node, rectangle image.Rectangle) []EntityImpl {

	var ents []EntityImpl

	if node == nil {
		return ents
	}

	xMin := rectangle.Min.X
	yMin := rectangle.Min.Y
	xMax := rectangle.Max.X
	yMax := rectangle.Max.Y

	p := image.Point{
		X: node.X,
		Y: node.Y,
	}

	if p.In(rectangle) {
		//Found
		ents = append(ents, node.Entity)
	}

	if xMin < node.X && yMin < node.Y {
		ents = append(ents, query(node.Nodes[SW], rectangle)...)
	}

	if xMin < node.X && yMax > node.Y {
		ents = append(ents, query(node.Nodes[NW], rectangle)...)
	}

	if xMax > node.X && yMin < node.Y {
		ents = append(ents, query(node.Nodes[SE], rectangle)...)
	}

	if xMax > node.X && yMax > node.Y {
		ents = append(ents, query(node.Nodes[NE], rectangle)...)
	}
	return ents
}

func (q *QuadTree) clear() {
	clear(q.Root)
	q.Root = nil
}

func clear(node *Node) {
	if node == nil {
		return
	}

	node.Entity = nil
	node.Nodes = make(map[int]*Node)
	for k := range node.Nodes {
		clear(node.Nodes[k])
	}
}
