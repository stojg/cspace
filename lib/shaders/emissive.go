package shaders

import "github.com/go-gl/gl/v4.1-core/gl"

type Emissive struct {
	Program  uint32
	LocModel int32
	LocColor int32
}

func NewEmissive() *Emissive {
	c := buildShader("emissive", "emissive")
	blockIndex := gl.GetUniformBlockIndex(c, gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(c, blockIndex, 0)

	return &Emissive{
		Program:  c,
		LocModel: loc(c, "model"),
		LocColor: loc(c, "emissive"),
	}
}
