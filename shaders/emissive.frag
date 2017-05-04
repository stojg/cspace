#version 410


uniform vec3 emissive;

out vec4 color;

void main() {
    color = vec4(emissive, 1.0f);
}
