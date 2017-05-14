#version 330 core
layout (location = 0) out vec4 FragColor;
layout (location = 1) out vec4 BrightColor;

in vec2 TexCoords;

uniform sampler2D screenTexture;

void main()
{
    vec4 color = texture(screenTexture, TexCoords);
    float brightness = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    FragColor = color;
    if(brightness > 0.9) {
        BrightColor = color;
    } else {
        BrightColor = vec4(0,0,0,1);
    }
}
