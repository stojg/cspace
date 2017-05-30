#version 330 core
out vec4 FragColor;
in  vec2 TexCoords;

uniform sampler2D fboAttachment;

void main()
{
    float t = texture(fboAttachment, TexCoords).a;
    FragColor = vec4(vec3(t), 1.0);
}
