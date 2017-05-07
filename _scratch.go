package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var near_plane float32 = 1
var far_plane float32 = 7.5

func blach() {
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
	//glm::mat4 projection = glm::perspective(Camera.Zoom, (float)SCR_WIDTH / (float)SCR_HEIGHT, 0.1f, 100.0f);
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
}

func emesiive() {

	// 2.5. Copy content of geometry's gDepth buffer to default framebuffer's gDepth buffer
	//gl.BindFramebuffer(gl.READ_FRAMEBUFFER, s.gbuffer.fbo)
	//gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0) // Write to default framebuffer
	////// blit to default framebuffer. Note that this may or may not work as the internal formats of both the FBO and default framebuffer have to match.
	////// the internal formats are implementation defined. This works on all of my systems, but if it doesn't on yours you'll likely have to write to the
	////// depth buffer in another shader stage (or somehow see to match the default framebuffer's internal format with the FBO's internal format).
	//gl.BlitFramebuffer(0, 0, s.width, s.height, 0, 0, s.width, s.height, gl.DEPTH_BUFFER_BIT, gl.NEAREST)
	//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	//
	//// 3. Render lights on top of scene, by blitting
	//s.shaderLightBox.UsePV(s.projection, view)
	//
	//for _, l := range s.pointLights {
	//	model := mgl32.Translate3D(l.Direction[0], l.Direction[1], l.Direction[2])
	//	model = model.Mul4(mgl32.Scale3D(0.02, 0.02, 0.02))
	//	//rad := l.Radius()
	//	//model = model.Mul4(mgl32.Scale3D(rad, rad, rad))
	//	setUniformMatrix4fv(s.shaderLightBox, "model", model)
	//	gl.Uniform3f(uniformLocation(s.shaderLightBox, "emissive"), l.Color[0], l.Color[1], l.Color[2])
	//	gl.BindVertexArray(s.lightMesh.vao)
	//	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.lightMesh.Vertices)))
	//	gl.BindVertexArray(0)
	//}
}
