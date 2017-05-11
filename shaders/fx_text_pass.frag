#version 330 core
in vec2 TexCoords;
layout(location = 0) out vec4 FragColor;

uniform sampler2D screenTexture;

void main()
{
    FragColor = texture(screenTexture, TexCoords);
}
