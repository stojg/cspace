package main

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
