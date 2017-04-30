#version 410

in vec3 vertex_position;

uniform mat4 matrix;

void main() {
	gl_Position = matrix * vec4(vertex_position, 1.0);
}
