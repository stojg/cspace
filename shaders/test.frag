#version 410

uniform sampler2D tex;

in vec3 ourColor;
in vec2 fragTexCoord;

out vec4 color;

void main() {
    color = texture(tex, fragTexCoord) * vec4(ourColor, 1.0f);
//	color = vec4(ourColor, 1.0f);
}

