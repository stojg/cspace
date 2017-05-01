#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 vert_texture_coordinate;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 transform;

out vec3 Normal;
out vec3 FragPos;
out vec2 Frag_texture_coordinate;

void main() {
    gl_Position = projection * view * transform * vec4(position, 1.0);
    FragPos = vec3(transform * vec4(position, 1.0f));
    Normal = mat3(transpose(inverse(transform))) * normal;
    Frag_texture_coordinate = vert_texture_coordinate;
}
