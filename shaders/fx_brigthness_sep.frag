#version 330 core
layout (location = 0) out vec4 Color;

in vec2 TexCoords;

uniform sampler2D screenTexture;

void main()
{
    vec4 hdrColor = texture(screenTexture, TexCoords);
    Color = hdrColor * 0.1;
}
