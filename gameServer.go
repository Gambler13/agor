package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gambler13/agor/api"
	"github.com/gambler13/agor/conf"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/parser"
	"image"
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
	X        float64
	Y        float64
	PlayerID string
	SeqID    uint32
}

type World struct {
	Players  map[string]*Player
	CellTree QuadTree
	Food     []EntityImpl
	FoodTree QuadTree
	Bounds   image.Rectangle
}

func InitWorld(conf conf.Game) World {

	world := conf.World

	bounds := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: world.Width,
			Y: world.Height,
		},
	}

	var food []EntityImpl

	i := 0
	for i < world.Food {

		f := &Food{Entity{
			Id: rand.Int(),
			Circle: Circle{
				Radius:   5,
				Position: getRandomPosition(bounds, 1),
			},
			Killer: 0,
			color:  randomLutIndex(),
		}}

		food = append(food, f)
		i++
	}

	ct := QuadTree{}

	ft := QuadTree{}

	for i := range food {
		ft.insert(food[i])
	}

	player := make(map[string]*Player)

	return World{
		Players:  player,
		CellTree: ct,
		Food:     food,
		FoodTree: ft,
		Bounds:   bounds,
	}

}

func (g *GameLoop) run() {
	Log.Info("start game loop")

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
			g.World.handlePosition(p)
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
		g.updateCells(delta, p[i])
	}

	g.World.CellTree.clear()
	for i := range p {
		for j := range p[i].Cells {
			if p[i].Cells[j].getEntity().Killer == 0 {
				g.World.CellTree.insert(p[i].Cells[j])
			}
		}
	}

	//Update food
	g.World.FoodTree.clear()
	for i := range g.World.Food {
		f := g.World.Food[i]
		if f.getEntity().Killer == 0 {
			g.World.FoodTree.insert(f)
		}
	}

	for i := range p {
		g.World.updatePlayers(p[i].SocketId)
	}

}

func (g *GameLoop) updateCells(delta float64, p *Player) {
	for j := range p.Cells {
		c := p.Cells[j]
		c.move(delta)
		c.eat(&g.World.FoodTree)
		c.eat(&g.World.CellTree)
	}
	pos := p.getCenter()
	point := image.Point{
		X: int(pos.X),
		Y: int(pos.Y),
	}

	if !point.In(g.World.Bounds) {
		for j := range p.Cells {
			c := p.Cells[j]
			//Undo cell move
			c.move((-1) * delta)
		}
	}
}

type Player struct {
	Id        int
	SocketId  string
	Mouse     Position64
	MouseSeq  uint32
	Cells     []*Cell
	conn      socketio.Conn
	startTS   time.Time
	foodEaten int
	updateCh  chan string
	cancel    func()
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

func (p *Player) runUpdateLoop(ctx context.Context) {
	for {
		select {
		case <-p.updateCh:

		case <-ctx.Done():
			return
		}
	}
}

func (p *Player) stopUpdateLoop() {
	p.cancel()
}

func (w *World) addNewPlayer(conn socketio.Conn) {
	Log.Infof("add player with socket id: %s", conn.ID())

	ctx, cancel := context.WithCancel(context.Background())

	player := &Player{
		Id:       rand.Int(),
		SocketId: conn.ID(),
		Mouse:    Position64{},
		conn:     conn,
		startTS:  time.Now(),
		cancel:   cancel,
	}
	//TODO add cell method on player or somethinng like thaht
	c := w.NewCell(player)
	player.Cells = []*Cell{&c}

	w.Players[conn.ID()] = player

	player.runUpdateLoop(ctx)
}

func (w *World) removePlayer(socketId string) {
	Log.Infof("stopUpdateLoop player with socket id: %s", socketId)
	p, ok := w.Players[socketId]
	if !ok {
		Log.Warnf("could not stopUpdateLoop player with id %s: not found", socketId)
		return
	}
	p.stopUpdateLoop()
	delete(w.Players, socketId)
}

func (w *World) handlePosition(msg PositionMsg) {
	p, ok := w.Players[msg.PlayerID]
	if !ok {
		return
	}
	p.Mouse = Position64{
		X: msg.X,
		Y: msg.Y,
	}
	p.MouseSeq = msg.SeqID
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

func (w *World) updatePlayers(id string) {
	for _, p := range w.Players {

		pos := p.getCenter()

		view := image.Rectangle{
			Min: image.Point{
				X: int(pos.X - 150),
				Y: int(pos.Y - 150),
			},
			Max: image.Point{
				X: int(pos.X + 150),
				Y: int(pos.Y + 150),
			},
		}

		cells := w.CellTree.query(view)

		food := w.FoodTree.query(view)

		entities := append(cells, food...)

		entityData := make([]byte, len(entities)*entities[0].getByteSize())

		for i := range entities {
			entityImpl := entities[i]
			e := entityImpl.getEntity()
			for j := range e.getByte() {
				entityData[i*e.getByteSize()+j] = e.getByte()[j]
			}
		}

		var socketBuf parser.Buffer
		socketBuf.Data = entityData

		apiPos := &api.Position{
			X: uint32(pos.X),
			Y: uint32(pos.Y),
		}

		var positionBuf parser.Buffer
		positionBuf.Data = apiPos.GetBytes()

		gameData, err := json.Marshal(api.GameStats{
			PlayerID:   fmt.Sprintf("%d", p.Id),
			Mass:       p.getMass(),
			FoodEaten:  p.foodEaten,
			CellsEaten: 0,
			Rank:       1,
			NumPlayers: len(w.Players),
		})
		if err != nil {
			Log.Errorf("error while marshalling game data entityData: %v", err)
			continue
		}

		if p.conn != nil {
			//TODO deadline / timeout?
			p.conn.Emit("update", &positionBuf, &socketBuf, string(gameData))
		}

	}
}
