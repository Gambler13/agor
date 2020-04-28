package main

import (
	"encoding/json"
	"fmt"
	"github.com/gambler13/agor/api"
	socketio "github.com/googollee/go-socket.io"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
	"time"
)

type GameLoop struct {
	tickRate       time.Duration
	quit           chan bool
	World          World
	PositionCh     chan PositionMsg
	AddPlayerCh    chan socketio.Conn
	RemovePlayerCh chan string
}

type PositionMsg struct {
	Position
	PlayerID string
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
		case p := <-g.PositionCh:
			g.World.handlePosition(p.PlayerID, p.Position)
		case c := <-g.AddPlayerCh:
			g.World.addNewPlayer(c)
		case i := <-g.RemovePlayerCh:
			g.World.removePlayer(i)
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
			c := p[i].Cells[j]
			c.move(delta)
			//TODO Check bounds after move

			c.eat(&g.World.FoodTree)
			c.eat(&g.World.CellTree)
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
	Mouse     Position
	Cells     []*Cell
	conn      socketio.Conn
	startTS   time.Time
	foodEaten int
}

func (p *Player) getCenter() Position {
	pos := make([]Position, len(p.Cells))
	for i := range p.Cells {
		pos[i] = p.Cells[i].Position
	}

	return centroid(pos)
}

func (p *Player) getMass() float64 {
	mass := 0.0
	for i := range p.Cells {
		mass += p.Cells[i].getEntity().Radius
	}
	return mass
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
		startTS:  time.Now(),
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

func (w *World) getLeaderboard() {

}

var nO = make(map[string]int)

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

		v, ok := nO[id]
		if ok && math.Abs(float64(v-len(ents))) > 5 {
			fmt.Printf("--- %d\n", len(ents))
		}

		nO[id] = len(ents)

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

		data, err := json.Marshal(results)
		if err != nil {
			fmt.Println("error while marhsalling results: " + err.Error())
		}
		posData, err := json.Marshal(api.Entity{Y: pos.Y, X: pos.X})
		if err != nil {
			fmt.Println("error while marhsalling results: " + err.Error())
		}

		gameData, err := json.Marshal(api.GameStats{
			PlayerId:   fmt.Sprintf("%d", p.Id),
			Mass:       p.getMass(),
			FoodEaten:  p.foodEaten,
			CellsEaten: 0,
			Rank:       1,
			NumPlayers: len(w.Players),
		})
		if err != nil {
			fmt.Println("error while marhsalling results: " + err.Error())
		}

		if p.conn != nil {
			p.conn.Emit("update", string(posData), string(data), string(gameData))
		}
	}
}
