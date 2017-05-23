#version 330 core
layout (location = 0) in vec2 position;
layout (location = 1) in vec2 texCoords;

float gAspectRatio = 1.6;
const float fov = 67.0;
const float fovR = 1.1693706;
float gTanHalfFOV = tan(fovR/2.0);

out vec2 TexCoords;
out vec2 ViewRay;

void main()
{
    gl_Position = vec4(position, 0.0f, 1.0f);
    TexCoords = texCoords;
    ViewRay.x = position.x * gAspectRatio * gTanHalfFOV;
    ViewRay.y = position.y * gTanHalfFOV;
}
