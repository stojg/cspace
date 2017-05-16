#version 410 core

layout (location = 0) out vec4 gNormalRoughness;
layout (location = 1) out vec4 gAlbedoMetallig;

in vec2 TexCoords;
in vec3 Normal;
in vec3 Tangent;

struct Material {
    vec3 albedo;
    float metallic;
    float roughness;
    float specularExp;
};
uniform Material mat;

void main()
{
    // store the per-fragment normals
    gNormalRoughness.xyz = Normal;
    // store the per-fragment roughness
    gNormalRoughness.w = mat.roughness;
    // store the per-fragment albedo/diffuse
    gAlbedoMetallig.rgb = mat.albedo;
    // store the per-fragment metallicness
    gAlbedoMetallig.a = mat.metallic;


}
