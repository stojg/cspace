package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/cspace/lib/shaders"
)

func NewSkymap(cubemap uint32) *Skybox {
	t := &Skybox{
		cubemap: cubemap,
		shader:  shaders.NewSkybox(),
	}
	return t
}

type Skybox struct {
	cubemap uint32
	shader  *shaders.Skybox
}

func (s *Skybox) Render(view mgl32.Mat4, cubeTexture uint32) {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.UseProgram(s.shader.Program)
	skyBoxView := view.Mat3().Mat4() // remove the rotation
	gl.UniformMatrix4fv(s.shader.LocSkyView, 1, false, &skyBoxView[0])
	GLBindCubeMap(0, s.shader.LocScreenTexture, cubeTexture)
	gl.BindVertexArray(s.shader.SkyboxVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

}
