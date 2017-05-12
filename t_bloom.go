package main

type BloomFBO struct {
	fbo      uint32
	textures []uint32
}

func NewBloom() *BloomFBO {
	frameBuffer := &BloomFBO{
		textures: make([]uint32, 2),
	}
	setFBO(&frameBuffer.fbo, frameBuffer.textures)
	return frameBuffer
}
