package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	x, y          float64
	radius        float64
	turnDirection int
	walkDirection int
	rotationAngle float64
	moveSpeed     float64
	rotationSpeed float64
	color         color.Color
	rays          []*Ray
}

func NewPlayer(x, y, radius float64, clr color.Color) *Player {
	p := &Player{
		x:             x,
		y:             y,
		radius:        radius,
		turnDirection: 0,
		walkDirection: 0,
		rotationAngle: math.Pi / 2,
		moveSpeed:     2.0,
		rotationSpeed: 2 * (math.Pi / 180),
		color:         clr,
		rays:          []*Ray{},
	}

	return p
}

func (p *Player) Update() error {
	// Handle player forward and backward movement.
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.walkDirection = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.walkDirection = -1
	} else {
		p.walkDirection = 0
	}

	// Handle player turning left / right rotation.
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.turnDirection = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.turnDirection = 1
	} else {
		p.turnDirection = 0
	}

	p.rotationAngle += float64(p.turnDirection) * p.rotationSpeed
	// mx, my := ebiten.CursorPosition()
	// dx := mx - int(p.x)
	// dy := my - int(p.y)

	// p.rotationAngle = NormalizeAngle(math.Atan2(float64(dy), float64(dx)))
	return nil
}

func (p *Player) Draw(dst *ebiten.Image) {
	vector.DrawFilledCircle(dst, MiniMapScaleFactor*float32(p.x), MiniMapScaleFactor*float32(p.y), MiniMapScaleFactor*float32(p.radius), p.color, true)
}

func (p *Player) MoveWithCollision(m *Map) error {
	moveStep := float64(p.walkDirection) * p.moveSpeed
	newX := p.x + math.Cos(p.rotationAngle)*moveStep
	newY := p.y + math.Sin(p.rotationAngle)*moveStep

	tile, err := m.AtPixel(int(newX), int(newY))
	if err != nil {
		return err
	}

	if tile == TileEmpty {
		cx, cy := ebiten.CursorPosition()
		dx := cx - int(newX)
		dy := cy - int(newY)

		if dist := math.Hypot(float64(dx), float64(dy)); dist > p.radius {
			p.x = newX
			p.y = newY
		}
	}

	return nil
}

func (p *Player) CastRays(fovAngle float64, numRays int, m *Map) error {
	rayAngle := p.rotationAngle - (fovAngle / 2)
	p.rays = []*Ray{}

	for i := 0; i < NumRays; i++ {
		ray := NewRay(p.x, p.y, rayAngle)
		if err := ray.Cast(m); err != nil {
			return err
		}
		p.rays = append(p.rays, ray)
		rayAngle += fovAngle / float64(numRays)
	}

	return nil
}
