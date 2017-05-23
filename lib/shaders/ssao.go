package shaders

import "fmt"

func NewSSAO() *SSAO {
	c := buildShader("ssao", "ssao")
	s := &SSAO{
		Program:          c,
		LocEnabled:       loc(c, "enabled"),
		LocProjection:    loc(c, "projection"),
		LocInvProjection: loc(c, "projMatrixInv"),
		LocGDepth:        loc(c, "gDepth"),
		LocScreenSize:    loc(c, "gScreenSize"),
		// see todo in ssao.frag
		//LocGNormal:       loc(c, "gNormal"),
		//LocTexNoise:      loc(c, "texNoise"),
	}

	for i := range s.LocSamples {
		s.LocSamples[i] = loc(c, fmt.Sprintf("samples[%d]", i))
	}
	return s
}

type SSAO struct {
	Program uint32

	LocEnabled       int32
	LocProjection    int32
	LocInvProjection int32
	LocGAlbedo       int32
	LocGDepth        int32
	LocSamples       [64]int32
	LocScreenSize    int32
	// see todo in ssao.frag
	//LocGNormal       int32
	//LocTexNoise      int32
}
