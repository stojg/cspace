package shaders

func NewFxaa() *Fxaa {
	c := buildShader("fx", "fx_fxaa")
	s := &Fxaa{
		Program:          c,
		LocInTexture:     loc(c, "screenTexture"),
		LocShowEdges:     loc(c, "u_showEdges"),
		LocLumaThreshold: loc(c, "u_lumaThreshold"),
		LocMulReduce:     loc(c, "u_mulReduce"),
		LocMinReduce:     loc(c, "u_minReduce"),
		LocMaxSpan:       loc(c, "u_maxSpan"),
	}

	//blockIndex := gl.GetUniformBlockIndex(c, gl.Str("Matrices\x00"))
	//gl.UniformBlockBinding(c, blockIndex, 0)

	return s
}

type Fxaa struct {
	Program uint32

	LocInTexture     int32
	LocShowEdges     int32
	LocLumaThreshold int32
	LocMulReduce     int32
	LocMinReduce     int32
	LocMaxSpan       int32
}
