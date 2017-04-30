#version 410

uniform sampler2D texture_diffuse1;
uniform sampler2D texture_diffuse2;

in vec3 ourColor;
in vec2 fragTexCoord;

out vec4 color;

void main() {
//    color = texture(tex, fragTexCoord) * vec4(ourColor, 1.0f);
    color = mix(texture(texture_diffuse1, fragTexCoord), texture(texture_diffuse2, fragTexCoord), 0.5);
}

