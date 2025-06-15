package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// Map related constants.
	MapColumns         = 15
	MapRows            = 11
	TileSize           = 64
	ScreenWidth        = TileSize * MapColumns
	ScreenHeight       = TileSize * MapRows
	MiniMapScaleFactor = 0.4
	// Tile type constants.
	TileEmpty = 0
	TileWall  = 1
	// Raycasting related constants
	FovAngle       = 60 * math.Pi / 180
	WallStripWidth = 4
	NumRays        = ScreenWidth / WallStripWidth
)

type Game struct {
	Map    *Map
	Player *Player
}

func (g *Game) Update() error {
	g.Player.Update()
	if err := g.Player.MoveWithCollision(g.Map); err != nil {
		return err
	}
	if err := g.Player.CastRays(FovAngle, NumRays, g.Map); err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{170, 170, 170, 255})
	g.Render3DProjectedWalls(screen)
	g.Map.Draw(screen)
	for _, r := range g.Player.rays {
		r.Draw(screen)
	}
	g.Player.Draw(screen)
}

func (g *Game) Layout(outWidth, outHeight int) (int, int) {
	return outWidth, outHeight
}

func NewGame() *Game {
	g := &Game{
		Map:    NewMap(MapRows, MapColumns, TileSize),
		Player: NewPlayer(ScreenWidth/2, ScreenHeight/2, 4, color.RGBA{255, 0, 0, 255}),
	}

	return g
}

func (g *Game) Render3DProjectedWalls(dst *ebiten.Image) {
	for i, r := range g.Player.rays {
		rayDistance := r.distance * math.Cos(r.angle-g.Player.rotationAngle)
		distanceProjectionPlane := (ScreenWidth / 2) / (math.Tan(FovAngle / 2))

		wallStripHeight := (TileSize / rayDistance) * distanceProjectionPlane

		const minDistance = 1.0
		const maxDistance = MapColumns * MapRows * 2 // tweak this for your map/scale

		clampedDistance := math.Max(minDistance, math.Min(rayDistance, maxDistance))
		shade := uint8(255 * (1 - (clampedDistance-minDistance)/(maxDistance-minDistance)))

		vector.DrawFilledRect(dst, float32(i*WallStripWidth),
			float32((ScreenHeight/2)-(wallStripHeight/2)),
			WallStripWidth,
			float32(wallStripHeight),
			color.RGBA{shade, shade, shade, 255},
			true,
		)
	}
}

func main() {
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
