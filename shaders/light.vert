#version 410

layout (location = 0) in vec3 position;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 transform;


void main() {
    gl_Position = projection * camera * model * transform * vec4(position, 1.0);
}
