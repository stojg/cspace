package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderType int

const (
	TexturedMesh ShaderType = iota
	MaterialMesh
)

type SceneNode interface {
	SimpleRender(ModelShader)
	Render(tShader *GbufferTShader, mShader *GbufferMShader)
	Add(mesh []*Mesh, transform mgl32.Mat4)
}

func NewBaseNode() SceneNode {
	origin := mgl32.Translate3D(0, 0, 0)
	q := &BaseNode{
		Node: Node{
			transform: &origin,
		},
	}
	return q
}

type BaseNode struct {
	Node
}

func (n *BaseNode) Render(tShader *GbufferTShader, mShader *GbufferMShader) {
	var tMeshes []*Node
	var mMeshes []*Node
	children := n.Node.Children()
	for _, child := range children {
		if child.mesh.MeshType == TexturedMesh {
			tMeshes = append(tMeshes, child)
		} else if child.mesh.MeshType == MaterialMesh {
			mMeshes = append(mMeshes, child)
		}
	}

	gl.UseProgram(tShader.Program())
	for i := range tMeshes {
		gl.UniformMatrix4fv(tShader.LocModel, 1, false, &tMeshes[i].transform[0])
		tMeshes[i].mesh.setTextures(tShader)
		tMeshes[i].mesh.Render()
	}

	gl.UseProgram(mShader.Program())
	for i := range mMeshes {
		gl.UniformMatrix4fv(mShader.locModel, 1, false, &mMeshes[i].transform[0])
		mMeshes[i].mesh.setMaterial(mShader)
		mMeshes[i].mesh.Render()
	}
}

func (n *BaseNode) SimpleRender(shader ModelShader) {
	for _, child := range n.Node.Children() {
		child.SimpleRender(shader)
	}
}

type Node struct {
	children  []*Node
	transform *mgl32.Mat4
	mesh      *Mesh
}

func (n *Node) SimpleRender(shader ModelShader) {
	gl.UniformMatrix4fv(shader.ModelUniform(), 1, false, &n.transform[0])
	n.mesh.Render()
	for _, child := range n.children {
		child.SimpleRender(shader)
	}
}

func (n *Node) Children() []*Node {
	var children []*Node
	for _, child := range n.children {
		children = append(children, child)
		children = append(children, child.Children()...)
	}
	return children
}

func (n *Node) Add(mesh []*Mesh, transform mgl32.Mat4) {
	for i := range mesh {
		child := &Node{
			mesh:      mesh[i],
			transform: &transform,
		}
		n.children = append(n.children, child)
	}
}
