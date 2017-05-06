// http://antongerdelan.net/opengl/glcontext2.html
package main

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const logFile = "gl.log"
const windowWidth = 800
const windowHeight = 600

var keys map[glfw.Key]bool
var cursor [2]float64

func main() {
	err := realMain()
	if err != nil {
		glError(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func realMain() error {

	keys = make(map[glfw.Key]bool)
	cursor[0] = windowWidth / 2
	cursor[1] = windowHeight / 2

	if err := restartLog(); err != nil {
		return err
	}
	defer glLogln("Program stopped")

	window, err := initWindow(windowWidth, windowHeight)
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	if err := initGL(); err != nil {
		return err
	}

	// this is pretty static for now. will need to be updated if window can change size
	//projection := mgl32.Perspective(mgl32.DegToRad(67.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	cam := newCamera()

	// load Mesh(es)
	cubeMesh := newCrateMesh()
	//lampMesh := newLightMesh()
	floor := newPlaneMesh()

	//ourShader, err := NewShader("material", "material")
	//if err != nil {
	//	return err
	//}

	//lampShader, err := NewShader("simple", "emissive")
	//if err != nil {
	//	return err
	//}

	//floorShader, err := NewShader("material", "floor")
	//if err != nil {
	//	return err
	//}

	//useNormalMapping := true

	//directionalColor := [3]float32{1.000, 0.99, 0.9}

	simpleDepthShader, err := NewShader("shadow_mapping_depth", "shadow_mapping_depth")
	if err != nil {
		return err
	}

	debugDepthQuad, err := NewShader("debugDepthQuad", "debugDepthQuad")
	if err != nil {
		return err
	}

	shader, err := NewShader("shadow", "shadow")
	if err != nil {
		return err
	}

	lightPos := [3]float32{-2.0, 4.0, -1.0}

	const SHADOW_WIDTH = 1024
	const SHADOW_HEIGHT = 1024
	var depthMapFBO uint32
	gl.GenFramebuffers(1, &depthMapFBO)

	var depthMap uint32
	gl.GenTextures(1, &depthMap)
	gl.BindTexture(gl.TEXTURE_2D, depthMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, SHADOW_WIDTH, SHADOW_HEIGHT, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// everything outside of the frustrum gets no shadow
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	borderColor := [4]float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, depthMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	woodTexture, err := newTexture(Diffuse, "textures/crate0/crate0_diffuse.png")
	if err != nil {
		return err
	}

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	previousTime := glfw.GetTime()
	for !window.ShouldClose() {

		// update timers
		now := glfw.GetTime()
		elapsed := float32(now - previousTime)
		previousTime = now

		glfw.PollEvents()
		fpsCounter(window)

		//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//if keys[glfw.Key1] {
		//	useNormalMapping = true
		//}
		//if keys[glfw.Key2] {
		//	useNormalMapping = false
		//}

		// update and get the camera view
		//view := cam.View(elapsed)

		// 1. Render depth of scene to texture (from light's perspective)
		// - Get light projection/view matrix.
		var near_plane float32 = 1
		var far_plane float32 = 7.5

		lightProjection := mgl32.Ortho(-10, 10, -10, 10, near_plane, far_plane)
		lightView := mgl32.LookAt(lightPos[0], lightPos[1], lightPos[2], 0, 0, 0, 0, 1, 0)
		lightSpaceMatrix := lightProjection.Mul4(lightView)
		// - render scene from light's point of view
		simpleDepthShader.Use()
		gl.UniformMatrix4fv(gl.GetUniformLocation(simpleDepthShader.Program, gl.Str("lightSpaceMatrix\x00")), 1, false, &lightSpaceMatrix[0])

		gl.Viewport(0, 0, SHADOW_WIDTH, SHADOW_HEIGHT)
		gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		drawScene(simpleDepthShader, cubeMesh, floor)
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

		// Reset viewport
		gl.Viewport(0, 0, windowWidth, windowHeight)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 2. Render scene as normal
		shader.Use()
		projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
		//glm::mat4 projection = glm::perspective(camera.Zoom, (float)SCR_WIDTH / (float)SCR_HEIGHT, 0.1f, 100.0f);
		view := cam.View(elapsed)
		setUniformMatrix4fv(shader, "projection", projection)
		setUniformMatrix4fv(shader, "view", view)
		// Set light uniforms
		gl.Uniform3fv(gl.GetUniformLocation(shader.Program, gl.Str("lightPos\x00")), 1, &lightPos[0])
		gl.Uniform3fv(gl.GetUniformLocation(shader.Program, gl.Str("viewPos\x00")), 1, &cam.position[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(shader.Program, gl.Str("lightSpaceMatrix\x00")), 1, false, &lightSpaceMatrix[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(uniformLocation(shader, "diffuseTexture"), 0)
		gl.BindTexture(gl.TEXTURE_2D, woodTexture.ID)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.Uniform1i(uniformLocation(shader, "shadowMap"), 1)
		gl.BindTexture(gl.TEXTURE_2D, depthMap)

		drawScene(shader, cubeMesh, floor)

		debugDepthQuad.Use()
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, depthMap)
		renderQuad()

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		window.SwapBuffers()

	}
	return nil
}

var quadVAO uint32 = 0
var quadVBO uint32 = 0

func renderQuad() {
	if quadVAO == 0 {
		quadVertices := []float32{
			// Positions        // Texture Coords
			0.5, 1.0, 0.0, 0.0, 1.0,
			0.5, 0.5, 0.0, 0.0, 0.0,
			1.0, 1.0, 0.0, 1.0, 1.0,
			1.0, 0.5, 0.0, 1.0, 0.0,
		}
		// Setup plane VAO
		gl.GenVertexArrays(1, &quadVAO)
		gl.GenBuffers(1, &quadVBO)
		gl.BindVertexArray(quadVAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, quadVBO)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	}
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
}

func drawScene(shader *Shader, cubeMesh, floorMesh *Mesh) {

	trans := mgl32.Translate3D(0, -0.5, 0)
	trans = trans.Mul4(mgl32.Scale3D(25, 0.1, 25))
	setUniformMatrix4fv(shader, "model", trans)
	floorMesh.Draw(shader)

	model := mgl32.Translate3D(0, 1.5, 0)
	setUniformMatrix4fv(shader, "model", model)
	cubeMesh.Draw(shader)

	model = mgl32.Translate3D(2.0, 0.0, 1)
	setUniformMatrix4fv(shader, "model", model)
	cubeMesh.Draw(shader)

	model = mgl32.Translate3D(-1.0, 0.0, 2)
	model = model.Mul4(mgl32.HomogRotate3D(45, mgl32.Vec3{1, 0, 1}.Normalize()))
	model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

	setUniformMatrix4fv(shader, "model", model)
	cubeMesh.Draw(shader)

}

func setDirectionalLight(shader *Shader, direction, color [3]float32) {
	name := fmt.Sprint("lights[0]")
	gl.Uniform4f(uniformLocation(shader, name+".vector"), direction[0], direction[1], direction[2], 0)
	gl.Uniform3f(uniformLocation(shader, name+".diffuse"), color[0], color[1], color[2])
	gl.Uniform3f(uniformLocation(shader, name+".ambient"), color[0]/10, color[1]/10, color[2]/10)
	gl.Uniform3f(uniformLocation(shader, name+".specular"), 1.0, 1.0, 1.0)
}

func setLights(shader *Shader, pos, color [][]float32) {
	for i := range pos {
		name := fmt.Sprintf("lights[%d]", i)
		gl.Uniform4f(uniformLocation(shader, name+".vector"), pos[i][0], pos[i][1], pos[i][2], 1)
		gl.Uniform3f(uniformLocation(shader, name+".diffuse"), color[i][0], color[i][1], color[i][2])
		gl.Uniform3f(uniformLocation(shader, name+".ambient"), color[i][0]/10, color[i][1]/10, color[i][2]/10)
		gl.Uniform3f(uniformLocation(shader, name+".specular"), color[i][0], color[i][1], color[i][2])
		gl.Uniform1f(uniformLocation(shader, name+".constant"), 1.0)
		gl.Uniform1f(uniformLocation(shader, name+".linear"), 0.14)
		gl.Uniform1f(uniformLocation(shader, name+".quadratic"), 0.07)
	}
}
