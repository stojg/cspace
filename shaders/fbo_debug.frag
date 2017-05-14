#version 330 core
out vec4 FragColor;
in  vec2 TexCoords;

uniform sampler2D fboAttachment;

void main()
{
    vec3 t = texture(fboAttachment, TexCoords).rgb;
    FragColor = vec4(t, 1.0);
}
