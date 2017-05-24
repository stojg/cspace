#version 410

layout (location = 0) in vec3 position;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
    mat4 invProjection;
    mat4 invView;
    vec3 cameraPos;
};

uniform mat4 model;

void main() {
    gl_Position = projection * view * model * vec4(position, 1.0);
}
