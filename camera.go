package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func NewCamera() *Camera {
	c := &Camera{
		position:   mgl32.Vec3{10, 5, 5},
		front:      mgl32.Vec3{0, 0, -1},
		up:         mgl32.Vec3{0, 1, 0},
		lastX:      float32(windowWidth / 2),
		lastY:      float32(windowHeight / 2),
		yaw:        -135,
		pitch:      -1,
		speed:      12.0,
		firstMouse: true,
	}
	c.updateVectors()
	c.view = mgl32.LookAtV(c.position, [3]float32{0, 0, 0}, c.up)
	return c
}

type Camera struct {
	position   mgl32.Vec3
	front      mgl32.Vec3
	up         mgl32.Vec3
	lastX      float32
	lastY      float32
	yaw        float32
	pitch      float32
	firstMouse bool
	speed      float32
	view       mgl32.Mat4
}

func (cam *Camera) View(elapsed float32) mgl32.Mat4 {
	changed := false
	if cam.handleKeyboard(elapsed) {
		changed = true
	}
	if cam.handleCursor(elapsed) {
		changed = true
	}

	if changed {
		cam.view = mgl32.LookAtV(cam.position, cam.position.Add(cam.front), cam.up)
	}
	return cam.view
}

func (cam *Camera) handleKeyboard(elapsed float32) bool {
	changed := false
	if keys[glfw.KeyW] {
		change := cam.front.Mul(cam.speed * elapsed)
		cam.position = cam.position.Add(change)
		changed = true
	}
	if keys[glfw.KeyS] {
		change := cam.front.Mul(cam.speed * elapsed)
		cam.position = cam.position.Sub(change)
		changed = true
	}
	if keys[glfw.KeyA] {
		change := cam.front.Cross(cam.up).Normalize().Mul(cam.speed * elapsed)
		cam.position = cam.position.Sub(change)
		changed = true
	}
	if keys[glfw.KeyD] {
		change := cam.front.Cross(cam.up).Normalize().Mul(cam.speed * elapsed)
		cam.position = cam.position.Add(change)
		changed = true
	}
	return changed
}

func (cam *Camera) handleCursor(elapsed float32) bool {
	xpos := cursor[0]
	ypos := cursor[1]

	if float32(xpos) == cam.lastX && float32(ypos) == cam.lastY {
		return false
	}

	if cam.firstMouse {
		cam.lastX = float32(xpos)
		cam.lastY = float32(ypos)
		cam.firstMouse = false
		return false
	}

	xOffset := float32(xpos) - cam.lastX
	yOffset := cam.lastY - float32(ypos)
	cam.lastX = float32(xpos)
	cam.lastY = float32(ypos)

	sensitivity := float32(10)
	xOffset *= sensitivity * elapsed
	yOffset *= sensitivity * elapsed

	cam.yaw += xOffset
	cam.pitch += yOffset

	if cam.pitch > 89 {
		cam.pitch = 89
	} else if cam.pitch < -89 {
		cam.pitch = -89
	}

	cam.updateVectors()

	return true
}

func (cam *Camera) updateVectors() {
	pitchRad := float64(mgl32.DegToRad(cam.pitch))
	yawRad := float64(mgl32.DegToRad(cam.yaw))
	cam.front[0] = float32(math.Cos(yawRad) * math.Cos(pitchRad))
	cam.front[1] = float32(math.Sin(pitchRad))
	cam.front[2] = float32(math.Sin(yawRad) * math.Cos(pitchRad))
	// normalise
	l := 1.0 / cam.front.Len()
	cam.front[0] *= l
	cam.front[1] *= l
	cam.front[2] *= l
}
