#version 410

uniform vec4 inputColor;

out vec4 frag_colour;

void main() {
	frag_colour = inputColor;
}

