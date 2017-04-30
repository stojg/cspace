#version 410

in vec3 vertex_position;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

void main() {
    gl_Position = projection * camera * model * vec4(vertex_position, 1);
}
