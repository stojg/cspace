#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 transform;

out vec3 Normal;
out vec3 FragPos;

void main() {
    gl_Position = projection * view * transform * vec4(position, 1.0);
    FragPos = vec3(transform * vec4(position, 1.0f));
    Normal = mat3(transpose(inverse(transform))) * normal;
}
