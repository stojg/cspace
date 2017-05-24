#version 410

uniform vec3 emissive;

out vec4 Color;

void main() {
    Color = vec4(emissive, 1.0f);
}
