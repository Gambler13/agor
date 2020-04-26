package main

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"golang.org/x/image/colornames"
	"hexhibit.xyz/agor/api"
	"image/color"
	"math/rand"
	"time"
)

type GameLoop struct {
	tickRate time.Duration
	quit     chan bool
	World    World
}

type World struct {
	Players  map[string]*Player
	CellTree Quadtree
	Food     []EntityImpl
	FoodTree Quadtree
	Bounds   Bounds
}

func InitWorld(b Bounds) World {

	var food []EntityImpl

	i := 0
	for i < 500 {

		var colName color.Color
		if i < 100 {
			colName = colornames.Bisque
		} else if i < 200 {
			colName = colornames.Aliceblue
		} else {
			colName = colornames.Greenyellow
		}

		f := &Food{Entity{
			Id: rand.Int(),
			Circle: Circle{
				Radius:   5,
				Position: getRandomPosition(b, 1),
				//Position: Position{X: float64(10.0 * i), Y: float64(10.0 * i)},
			},
			Killer: 0,
			Color:  colName,
		}}

		food = append(food, f)
		i++
	}

	p1 := &Player{
		Id:       rand.Int(),
		SocketId: "aaaa",
		Mouse: Position{
			X: 1,
			Y: 0,
		},
	}

	c := &Cell{
		Entity: Entity{
			Circle: Circle{
				Radius: 50,
				Position: Position{
					X: 100,
					Y: 100,
				},
			},
			Id:     1,
			Killer: 0,
			Color:  colornames.Red,
		},
		Owner: p1,
	}

	p1.Cells = append(p1.Cells, c)

	p2 := &Player{
		SocketId: "bbb",
		Id:       rand.Int(),
		Mouse: Position{
			X: -1,
			Y: 0,
		},
	}

	c2 := &Cell{
		Entity: Entity{
			Circle: Circle{
				Radius: 20,
				Position: Position{
					X: 300,
					Y: 100,
				},
			},
			Id:     11,
			Killer: 0,
			Color:  colornames.Beige,
		},
		Owner: p2,
	}

	p2.Cells = append(p2.Cells, c2)

	ct := Quadtree{
		Bounds:     b,
		MaxObjects: 20,
		MaxLevels:  10,
	}

	ft := Quadtree{
		Bounds:     b,
		MaxObjects: 20,
		MaxLevels:  10,
	}

	for i := range food {
		ft.Insert(food[i])
	}

	player := make(map[string]*Player)
	player[p1.SocketId] = p1
	player[p2.SocketId] = p2

	return World{
		Players:  player,
		CellTree: ct,
		Food:     food,
		FoodTree: ft,
		Bounds:   b,
	}

}

func (g *GameLoop) run() {
	tickInterval := time.Second / g.tickRate
	timeStart := time.Now().UnixNano()

	ticker := time.NewTicker(tickInterval)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			// DT in seconds
			delta := float64(now-timeStart) / 1000000000
			timeStart = now
			g.onUpdate(delta)

		case <-g.quit:
			ticker.Stop()
		}
	}
}

func (g *GameLoop) onUpdate(delta float64) {
	p := g.World.Players

	//Update cells
	for i := range p {
		for j := range p[i].Cells {
			p[i].Cells[j].move(delta)
			p[i].Cells[j].eat(&g.World.FoodTree)
			p[i].Cells[j].eat(&g.World.CellTree)
		}
	}

	g.World.CellTree.Clear()
	for i := range p {
		for j := range p[i].Cells {
			if p[i].Cells[j].getEntity().Killer == 0 {
				g.World.CellTree.Insert(p[i].Cells[j])
			}
		}
	}

	//Update food
	g.World.FoodTree.Clear()
	for i := range g.World.Food {
		f := g.World.Food[i]
		if f.getEntity().Killer == 0 {
			g.World.FoodTree.Insert(f)
		}
	}

	for i := range p {
		g.World.updatePlayers(p[i].SocketId)
	}

}

type Player struct {
	Id       int
	SocketId string
	//Normalized vector based on players center
	Mouse Position
	Cells []*Cell
	conn  socketio.Conn
}

func (p *Player) getCenter() Position {
	pos := make([]Position, len(p.Cells))
	for i := range p.Cells {
		pos[i] = p.Cells[i].Position
	}

	return centroid(pos)
}

//Return normalized vector from cell center to mouse position
func (p *Player) getMouseVector() Position {
	return sub(p.Mouse, p.getCenter())
}

//Distribute eaten food between players cell
func (p *Player) distributeFood(f EntityImpl) {
	r := f.getEntity().Radius * 0.075
	r = r / float64(len(p.Cells))
	for i := range p.Cells {
		p.Cells[i].Radius += r
	}
}

func (p *Player) splitCells() {
	newCells := make([]*Cell, len(p.Cells)*2)
	for i := range p.Cells {
		c1 := p.Cells[i]
		c2 := c1.split()
		newCells[i*2] = c1
		newCells[i*2+1] = c2
	}
	p.Cells = newCells
}

func (w *World) addNewPlayer(conn socketio.Conn) {

	player := &Player{
		Id:       rand.Int(),
		SocketId: conn.ID(),
		Mouse:    Position{},
		conn:     conn,
	}
	//TODO add cell method on player or somethinng like thaht
	c := w.NewCell(player)
	player.Cells = []*Cell{&c}

	w.Players[conn.ID()] = player
}

func (w *World) removePlayer(socketId string) {
	delete(w.Players, socketId)
}

func (w *World) handlePosition(id string, mousePos Position) {
	p := w.Players[id]
	p.Mouse = mousePos
}

func (w *World) handleSplit(id string) {
	p := w.Players[id]
	p.splitCells()
}

func (w *World) handleDiet(id string) {
	//TODO implement
}

func (w *World) updatePlayers(id string) {
	for _, p := range w.Players {

		pos := p.getCenter()

		ents := w.CellTree.RetrieveViewIntersections(Bounds{
			X:      pos.X,
			Y:      pos.Y,
			Width:  300,
			Height: 300,
		})
		ents2 := w.FoodTree.RetrieveViewIntersections(Bounds{
			X:      pos.X,
			Y:      pos.Y,
			Width:  300,
			Height: 300,
		})

		ents = append(ents, ents2...)

		results := make([]api.Entity, len(ents))

		for i := range ents {
			entityImpl := ents[i]
			e := entityImpl.getEntity()
			results[i] = api.Entity{
				X:      e.X,
				Y:      e.Y,
				Radius: e.Radius,
				Color:  hexColor(e.Color),
			}
		}

		data, _ := json.Marshal(results)
		posData, _ := json.Marshal(api.Entity{Y: pos.Y, X: pos.X})

		if p.conn != nil {
			p.conn.Emit("update", string(posData), string(data))
		}
	}
}
