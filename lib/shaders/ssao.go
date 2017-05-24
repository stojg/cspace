package shaders

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewSSAO() *SSAO {
	c := buildShader("ssao", "ssao")
	s := &SSAO{
		Program:       c,
		LocEnabled:    loc(c, "enabled"),
		LocGDepth:     loc(c, "gDepth"),
		LocScreenSize: loc(c, "gScreenSize"),
		// see todo in ssao.frag
		//LocGNormal:       loc(c, "gNormal"),
		//LocTexNoise:      loc(c, "texNoise"),
	}

	blockIndex := gl.GetUniformBlockIndex(c, gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(c, blockIndex, 0)

	for i := range s.LocSamples {
		s.LocSamples[i] = loc(c, fmt.Sprintf("samples[%d]", i))
	}
	return s
}

type SSAO struct {
	Program uint32

	LocEnabled    int32
	LocGAlbedo    int32
	LocGDepth     int32
	LocSamples    [64]int32
	LocScreenSize int32
	// see todo in ssao.frag
	//LocGNormal       int32
	//LocTexNoise      int32
}
