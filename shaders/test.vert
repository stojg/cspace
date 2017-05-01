#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec2 vert_texture_coordinate;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform mat4 transform;

out vec2 frag_texture_coordinate;

void main() {
    gl_Position = projection * view * model * transform * vec4(position, 1.0);
    frag_texture_coordinate = vert_texture_coordinate;
}
