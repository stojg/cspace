#version 410

layout (location = 0) in vec3 position;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
};
uniform mat4 transform;

void main() {
    gl_Position = projection * view * transform * vec4(position, 1.0);
}
