package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Node struct {
	children  []*Node
	transform *mgl32.Mat4
	mesh      *Mesh
}

func (n *Node) Render(shader ModelShader) {

	for _, child := range n.children {
		if child.transform != nil {
			transform := child.transform.Mul4(*n.transform)
			gl.UniformMatrix4fv(shader.ModelUniform(), 1, false, &transform[0])
			child.mesh.Render(shader)
		}
		child.Render(shader)
	}
}

func (n *Node) Destroy() {
	n.children = make([]*Node, 0)
}

func (n *Node) Add(mesh *Mesh, transform mgl32.Mat4) *Node {
	child := &Node{
		mesh:      mesh,
		transform: &transform,
	}
	n.children = append(n.children, child)
	return child
}
