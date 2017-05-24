#version 330 core
layout (location = 0) in vec3 position;
out vec3 TexCoords;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
    mat4 invProjection;
    mat4 invView;
    vec3 cameraPos;
};

uniform mat4 skyView;

void main()
{
    vec4 pos =   projection * skyView * vec4(position, 1.0);
    gl_Position = pos.xyww;
    TexCoords = position;

}
