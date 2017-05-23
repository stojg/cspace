#version 330 core

in vec2 TexCoords;
in vec2 ViewRay;

out vec4 FragColor;

uniform sampler2D screenTexture;

void main()
{
    float Depth = texture(screenTexture, TexCoords).x;
    FragColor = vec4(vec3(Depth), 1.0);
}
