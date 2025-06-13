package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ray represents a single ray cast from the player's position.
// It stores the origin of the ray, its angle, direction flags, and hit distance.
type Ray struct {
	originX, originY float64 // Origin point (usually player position)
	angle            float64 // Normalized angle of the ray in radians
	distance         float64 // Distance from origin to wall hit
	isRayFacingDown  bool
	isRayFacingUp    bool
	isRayFacingRight bool
	isRayFacingLeft  bool
	isHitVertical    bool
}

// NewRay creates a new ray from a given position and angle, and calculates direction flags.
func NewRay(x, y, angle float64) *Ray {
	r := &Ray{
		originX: x,
		originY: y,
		angle:   NormalizeAngle(angle),
	}

	r.isRayFacingDown = r.angle > 0 && r.angle < math.Pi
	r.isRayFacingUp = !r.isRayFacingDown
	r.isRayFacingRight = r.angle < 0.5*math.Pi || r.angle > 1.5*math.Pi
	r.isRayFacingLeft = !r.isRayFacingRight

	return r
}

// Draw renders the ray as a line based on its distance and angle.
// This is a debug visual — not used for wall rendering.
func (r *Ray) Draw(dst *ebiten.Image) {
	vector.StrokeLine(
		dst,
		MiniMapScaleFactor*float32(r.originX),
		MiniMapScaleFactor*float32(r.originY),
		MiniMapScaleFactor*(float32(r.originX)+float32(math.Cos(r.angle)*r.distance)),
		MiniMapScaleFactor*(float32(r.originY)+float32(math.Sin(r.angle)*r.distance)),
		1,
		color.RGBA{0, 0, 200, 128},
		true,
	)
}

// Cast performs a raycast in both horizontal and vertical directions,
// then stores the shortest distance found.
func (r *Ray) Cast(m *Map) error {
	hx, hy, err := r.findHorizontalIntersects(m)
	if err != nil {
		return err
	}

	vx, vy, err := r.findVerticalIntersects(m)
	if err != nil {
		return err
	}

	hd := math.Hypot(hx-r.originX, hy-r.originY) // horizontal hit distance
	vd := math.Hypot(vx-r.originX, vy-r.originY) // vertical hit distance

	epsilon := 0.0001
	if math.Abs(vd-hd) < epsilon {
		r.isHitVertical = vd < hd
	} else {
		r.isHitVertical = vd < hd
	}

	r.distance = math.Min(hd, vd)

	return nil
}

// findHorizontalIntersects calculates where the ray first hits a horizontal grid line.
// It steps vertically through rows and calculates X based on Y and the ray angle.
func (r *Ray) findHorizontalIntersects(m *Map) (float64, float64, error) {
	yIntercept := math.Floor(r.originY/TileSize) * TileSize
	if r.isRayFacingDown {
		yIntercept += TileSize
	}

	tanAngle := math.Tan(r.angle)
	if tanAngle == 0 {
		tanAngle = 0.0001
	}
	xIntercept := r.originX + (yIntercept-r.originY)/tanAngle

	yStep := float64(TileSize)
	if r.isRayFacingUp {
		yStep *= -1 // check just above the grid line
	}

	xStep := TileSize / tanAngle
	if r.isRayFacingLeft && xStep > 0 || r.isRayFacingRight && xStep < 0 {
		xStep *= -1
	}

	offsetY := 0.0
	if r.isRayFacingUp {
		offsetY = -1
	}

	return r.traceRay(
		xIntercept, yIntercept,
		xStep, yStep,
		0, offsetY, func(x, y float64) (int, error) {
			return m.AtPixel(int(x), int(y))
		},
	)
}

// findVerticalIntersects calculates where the ray first hits a vertical grid line.
// It steps horizontally through columns and calculates Y based on X and the ray angle.
func (r *Ray) findVerticalIntersects(m *Map) (float64, float64, error) {
	xIntercept := math.Floor(r.originX/TileSize) * TileSize
	if r.isRayFacingRight {
		xIntercept += TileSize
	}

	tanAngle := math.Tan(r.angle)
	if tanAngle == 0 {
		tanAngle = 0.0001
	}
	yIntercept := r.originY + (xIntercept-r.originX)*tanAngle

	// Calculate the xDelta and yDelta
	xStep := float64(TileSize)
	if r.isRayFacingLeft {
		xStep *= -1
	}

	yStep := TileSize * tanAngle
	if r.isRayFacingUp && yStep > 0 || r.isRayFacingDown && yStep < 0 {
		yStep *= -1
	}

	offsetX := 0.0
	if r.isRayFacingLeft {
		offsetX = -1 // check just before grid line
	}

	return r.traceRay(
		xIntercept, yIntercept,
		xStep, yStep,
		offsetX, 0,
		func(x, y float64) (int, error) {
			return m.AtPixel(int(x), int(y))
		},
	)
}

// traceRay steps from an initial intersection point using x/y deltas,
// checking the map at each step until it hits a wall or goes out of bounds.
// `offsetX` and `offsetY` adjust the sample point so rays hitting from
// above/left sample the correct adjacent tile (since gridlines lie between tiles).
func (r *Ray) traceRay(startX, startY, stepX, stepY, offsetX, offsetY float64, sampleTile func(x, y float64) (int, error)) (float64, float64, error) {
	currX := startX
	currY := startY

	for currX >= 0 && currX < ScreenWidth && currY >= 0 && currY < ScreenHeight {
		sampleX := currX + offsetX
		sampleY := currY + offsetY

		tile, err := sampleTile(sampleX, sampleY)
		if err != nil {
			return 0, 0, err
		}

		if tile == TileWall {
			return currX, currY, nil
		}

		currX += stepX
		currY += stepY
	}

	return 0, 0, nil
}
