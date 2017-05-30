package main

import "github.com/go-gl/gl/v4.1-core/gl"

var u_lumaThreshold float32 = 0.6 // (1/3), (1/4), (1/8), (1/16)
var u_mulReduce float32 = 1 / 8.0 //
var u_minReduce float32 = 1 / 128.0
var u_maxSpan float32 = 8.0

var fxaaShader *DefaultShader
var fxaaTextureloc int32

var fxaaLocU_showEdges int32
var fxaaLocU_lumaThreshold int32
var fxaaLocU_mulReduce int32
var fxaaLocU_minReduce int32
var fxaaLocU_maxSpan int32
var fxaaLoc_enabled int32

func initFxaa() {
	fxaaShader = NewDefaultShader("fx", "fx_fxaa")
	fxaaTextureloc = uniformLocation(fxaaShader, "screenTexture")
	fxaaLocU_showEdges = pUniformLocation(fxaaShader.program, "u_showEdges")
	fxaaLocU_lumaThreshold = pUniformLocation(fxaaShader.program, "u_lumaThreshold")
	fxaaLocU_mulReduce = pUniformLocation(fxaaShader.program, "u_mulReduce")
	fxaaLocU_minReduce = pUniformLocation(fxaaShader.program, "u_minReduce")
	fxaaLocU_maxSpan = pUniformLocation(fxaaShader.program, "u_maxSpan")
	fxaaLoc_enabled = pUniformLocation(fxaaShader.program, "u_enabled")
}

func renderFxaa(inTexture uint32) {
	fxaaShader.Use()
	if fxaaOn {
		gl.Uniform1i(fxaaLoc_enabled, 1)
	} else {
		gl.Uniform1i(fxaaLoc_enabled, 0)
	}
	if showDebug {
		gl.Uniform1i(fxaaLocU_showEdges, 1)
	} else {
		gl.Uniform1i(fxaaLocU_showEdges, 0)
	}
	gl.Uniform1f(fxaaLocU_lumaThreshold, u_lumaThreshold)
	gl.Uniform1f(fxaaLocU_minReduce, u_minReduce)
	gl.Uniform1f(fxaaLocU_mulReduce, u_mulReduce)
	gl.Uniform1f(fxaaLocU_maxSpan, u_maxSpan)
	GLBindTexture(0, fxaaTextureloc, inTexture)
	renderQuad()
}
