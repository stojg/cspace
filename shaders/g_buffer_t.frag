#version 410 core

layout (location = 0) out vec4 gNormalRoughness;
layout (location = 1) out vec4 gAlbedoMetallic;

in vec2 TexCoords;
in vec3 Normal;
in mat3 TBN;

struct Material {
    sampler2D albedo;
    sampler2D metallic;
    sampler2D normal;
    sampler2D roughness;
};
uniform Material mat;

vec3 CalcBumpedNormal(vec3 normal);

void main()
{
    // store the per-fragment normals
    gNormalRoughness.rgb = CalcBumpedNormal(Normal).rgb;
    // store the per-fragment roughness
    gNormalRoughness.a = texture(mat.roughness, TexCoords).r;

    // And the diffuse per-fragment color
    gAlbedoMetallic.rgb = texture(mat.albedo, TexCoords).rgb;
    // Store specular intensity in gAlbedoSpec's alpha component
    gAlbedoMetallic.a = texture(mat.metallic, TexCoords).r;
}

vec3 CalcBumpedNormal(vec3 normal)
{
    vec3 BumpMapNormal = texture(mat.normal, TexCoords).xyz;
    BumpMapNormal = 2.0 * BumpMapNormal - vec3(1.0);
    return normalize(TBN * BumpMapNormal);
}
