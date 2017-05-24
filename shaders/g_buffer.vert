#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoords;
layout (location = 3) in vec3 tangent;

out vec2 TexCoords;
out vec3 Normal;
out vec3 Tangent;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
};

uniform mat4 model;

void main()
{
    mat4 mvp = projection * view * model;
    gl_Position = mvp * vec4(position, 1.0);

    TexCoords = texCoords;
    Normal = transpose(inverse(mat3(model))) * normal;
    Tangent = normalize(vec3(model * vec4(tangent,   0.0)));
}
