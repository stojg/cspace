package shaders

import "github.com/go-gl/gl/v4.1-core/gl"

func NewSSAODepthResampler() *SSAODepthResampler {
	c := buildShader("ssao", "ssao_depth_resample")
	s := &SSAODepthResampler{
		Program:       c,
		LocGDepth:     loc(c, "gDepth"),
		LocScreenSize: loc(c, "gScreenSize"),
	}

	blockIndex := gl.GetUniformBlockIndex(c, gl.Str("Matrices\x00"))
	gl.UniformBlockBinding(c, blockIndex, 0)
	return s
}

type SSAODepthResampler struct {
	Program uint32

	LocEnabled    int32
	LocGAlbedo    int32
	LocGDepth     int32
	LocSamples    [64]int32
	LocScreenSize int32
	LocGNormal    int32
	LocTexNoise   int32
}
