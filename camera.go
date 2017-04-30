package main

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func newCamera() *camera {
	c := &camera{
		position:   mgl32.Vec3{0, 0, 10},
		front:      mgl32.Vec3{0, 0, -1},
		up:         mgl32.Vec3{0, 1, 0},
		lastX:      windowWidth / 2,
		lastY:      windowHeight / 2,
		yaw:        -90.0,
		pitch:      0,
		speed:      5.0,
		firstMouse: true,
	}
	return c
}

type camera struct {
	position   mgl32.Vec3
	front      mgl32.Vec3
	up         mgl32.Vec3
	lastX      float32
	lastY      float32
	yaw        float32
	pitch      float32
	firstMouse bool
	speed      float32
}

func (cam *camera) Update(shader *Shader, elapsed float32) {
	changed := false
	if cam.handleKeyboard(elapsed) {
		changed = true
	}
	if cam.handleCursor(elapsed) {
		changed = true
	}

	if changed {
		cam.Draw(shader)
	}
}

func (cam *camera) Draw(shader *Shader) {
	mat := mgl32.LookAtV(cam.position, cam.position.Add(cam.front), cam.up)
	cameraUniform := gl.GetUniformLocation(shader.Program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &mat[0])
}

func (cam *camera) handleKeyboard(elapsed float32) bool {
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

func (cam *camera) handleCursor(elapsed float32) bool {
	xpos := cursor[0]
	ypos := cursor[1]

	if float32(xpos) == cam.lastX && float32(ypos) == cam.lastY {
		return false
	}

	if cam.firstMouse {
		cam.lastX = float32(xpos)
		cam.lastY = float32(ypos)
		cam.firstMouse = false
	}

	xOffset := float32(xpos) - cam.lastX
	yOffset := cam.lastY - float32(ypos)
	cam.lastX = float32(xpos)
	cam.lastY = float32(ypos)

	sensitivity := float32(5)
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

func (cam *camera) updateVectors() {
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
