#version 410

uniform sampler2D texture_diffuse1;
uniform sampler2D texture_diffuse2;

in vec2 frag_texture_coordinate;

out vec4 color;

void main() {

    vec4 t2 = texture(texture_diffuse2, frag_texture_coordinate);
    float intensity = t2.r + t2.g + t2.b;
    if(intensity < 0.1) {
        color = texture(texture_diffuse1, frag_texture_coordinate);
        return;
    }
    color = mix(texture(texture_diffuse1, frag_texture_coordinate), texture(texture_diffuse2, frag_texture_coordinate), 0.1);
}

