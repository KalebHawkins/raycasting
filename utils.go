package main

import "math"

// NormalizeAngle will clamp and angle to always stay somewhere between 0 and 2PI.
func NormalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}
