#version 330 core
in vec2 TexCoords;
out vec4 color;

uniform sampler2D screenTexture;

void main()
{
     const float gamma = 2.2;
     vec3 hdrColor = texture(screenTexture, TexCoords).rgb;

     // Reinhard tone mapping
     vec3 mapped = hdrColor / (hdrColor + vec3(1.0));
     // Gamma correction
     mapped = pow(mapped, vec3(1.0 / gamma));

     color = vec4(mapped, 1.0);
}
