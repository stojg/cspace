#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoords;
layout (location = 3) in vec3 tangent;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
};

//uniform mat4 projection;
//uniform mat4 view;
uniform mat4 transform;

out vec3 Normal;
out vec3 FragPos;
out vec2 FragTexCoords;
out mat3 TBN;

void main() {
    gl_Position = projection * view * transform * vec4(position, 1.0);
    FragPos = vec3(transform * vec4(position, 1.0f));
    Normal =  mat3(transpose(inverse(transform))) * normal;
    FragTexCoords = texCoords;

    vec3 T = normalize(vec3(transform * vec4(tangent, 0.0)));
    vec3 N = normalize(vec3(transform * vec4(normal, 0.0)));
    // re-orthogonalize T with respect to N
    T = normalize(T - dot(T, N) * N);
    // then retrieve perpendicular vector B with the cross product of T and N
    vec3 B = cross(N, T);

//    vec3 T = normalize(vec3(transform * vec4(tangent,   0.0)));
//    vec3 B = normalize(vec3(transform * vec4(bitangent, 0.0)));
//    vec3 N = normalize(vec3(transform * vec4(normal,    0.0)));
    TBN = mat3(T, B, N);
}
