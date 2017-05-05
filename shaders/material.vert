#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoords;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitangent;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 transform;
uniform vec3 lightPos;
uniform vec3 viewPos;

out vec3 FragNormal;
out vec3 FragPos;
out vec2 FragTexCoords;
out vec3 TangentLightPos;
out vec3 TangentViewPos;
out vec3 TangentFragPos;

void main() {
    vec3 T   = normalize(mat3(transform) * tangent);
    vec3 B   = normalize(mat3(transform) * bitangent);
    vec3 N   = normalize(mat3(transform) * normal);
    if (dot(cross(N, T), B) < 0.0) {
        T = T * -1.0;
    }
    mat3 TBN = transpose(mat3(T, B, N));

    FragNormal = normal;
    FragPos = vec3(transform * vec4(position, 1.0f));
    FragTexCoords = texCoords;
    TangentLightPos = TBN * lightPos;
    TangentViewPos  = TBN * viewPos;
    TangentFragPos  = TBN * FragPos;

     gl_Position = projection * view * transform * vec4(position, 1.0);
}
