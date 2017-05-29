#version 410 core

layout (location = 0) out vec4 gNormalRoughness;
layout (location = 1) out vec4 gAlbedoMetallic;

in vec2 TexCoords;
in vec3 Normal;

struct Material {
    vec3 albedo;
    float metallic;
    float roughness;
};
uniform Material mat;

void main()
{
    // store the per-fragment normals
    gNormalRoughness.xyz = Normal;
    // store the per-fragment roughness
    gNormalRoughness.w = mat.roughness;
    // store the per-fragment albedo/diffuse
    gAlbedoMetallic.rgb = mat.albedo;
    // store the per-fragment metallicness
    gAlbedoMetallic.a = mat.metallic;
}
