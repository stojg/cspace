#version 410

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 color;
layout (location = 2) in vec2 vertTexCoord;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 transform;

out vec2 fragTexCoord;
out vec3 ourColor; // Output a color to the fragment shader

void main() {
    gl_Position = projection * camera * model * transform * vec4(position, 1.0);
    ourColor = color; // Set ourColor to the input color we got from the vertex data
    fragTexCoord = vertTexCoord;
}
