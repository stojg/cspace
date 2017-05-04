#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoords;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitangent;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 transform;

out vec3 Normal;
out vec3 FragPos;
out vec2 FragTexCoords;
out vec3 Tangent;

void main() {
    gl_Position = projection * view * transform * vec4(position, 1.0);
    FragPos = vec3(transform * vec4(position, 1.0f));
    Normal =  mat3(transpose(inverse(transform))) * normal;
    FragTexCoords = texCoords;
    Tangent = normalize(vec3(transform * vec4(tangent,   0.0)));
}
