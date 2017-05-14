#version 330 core

layout (location = 0) out vec3 gNormal;
layout (location = 1) out vec4 gAlbedoSpec;

in vec2 TexCoords;
in vec3 Normal;
in vec3 Tangent;

struct Material {
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float specularExp;
    float Transparency; // unused for now
};
uniform Material mat;

void main()
{
    // Also store the per-fragment normals into the gbuffer
    gNormal = Normal;
    // And the diffuse per-fragment color
    gAlbedoSpec.rgb = mat.diffuse;
    // Store specular intensity in gAlbedoSpec's alpha component
    gAlbedoSpec.a = mat.specularExp;
}
