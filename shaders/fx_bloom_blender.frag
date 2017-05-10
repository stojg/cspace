#version 330 core

out vec4 FragColor;
in vec2 TexCoords;

uniform sampler2D screenTexture;
uniform sampler2D bloomTexture;

void main()
{
    vec3 hdrColor = texture(screenTexture, TexCoords).rgb;
    vec3 bloomColor = texture(bloomTexture, TexCoords).rgb;
    hdrColor += bloomColor; // additive blending
    FragColor = vec4(hdrColor, 1.0f);
}
