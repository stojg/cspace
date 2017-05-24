#version 410 core

out vec4 FragColor;

const int NR_LIGHTS = 32;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
    mat4 invProjection;
    mat4 invView;
};

uniform sampler2D gDepth;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;
uniform sampler2D gAmbientOcclusion;

struct Light {
    vec3 Position;
    vec3 Color;
    float Linear;
    float Quadratic;
};
uniform Light pointLight[NR_LIGHTS];

uniform vec3 viewPos;

uniform vec2 gScreenSize;
uniform int numLights = 1;

const float PI = 3.14159265359;

vec2 CalcTexCoord();
vec3 WorldPosFromDepth(float depth, vec2 TexCoords);
vec3 fresnelSchlick(float cosTheta, vec3 F0);
float DistributionGGX(vec3 N, vec3 H, float roughness);
float GeometrySchlickGGX(float NdotV, float roughness);
float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness);

void main()
{
    vec2 TexCoords = CalcTexCoord();
    vec3 FragPos   = WorldPosFromDepth(texture(gDepth, TexCoords).x, TexCoords);

    vec3 N = texture(gNormal, TexCoords).rgb;
    vec3 V = normalize(viewPos - FragPos);

    vec3 albedo = texture(gAlbedoSpec, TexCoords).rgb;
    float metallic = texture(gAlbedoSpec, TexCoords).a;
    float roughness = texture(gNormal, TexCoords).w;
    float ao = texture(gAmbientOcclusion, TexCoords).r;

    vec3 Lo = vec3(0.0);

    vec3 ambient = vec3(0.0);
    for(int i = 0; i < numLights; i++){

        vec3 kD = vec3(0.0);
        vec3 L = normalize(pointLight[i].Position - FragPos);
        vec3 H = normalize(V + L);

        float distance = length(L);
        float attenuation = 1.0 / (1.0 + pointLight[i].Linear * distance + pointLight[i].Quadratic * distance * distance * distance);
        // float attenuation = 1.0 / (distance * distance);
        vec3 radiance     = pointLight[i].Color * attenuation;

        vec3 F0 = vec3(0.04);
        F0      = mix(F0, albedo, metallic);
        vec3 F  = fresnelSchlick(max(dot(H, V), 0.0), F0);

        float NDF = DistributionGGX(N, H, roughness);
        float G   = GeometrySmith(N, V, L, roughness);

        vec3 specular = vec3(0.0);
        // Cook-Torrance BRDF
        vec3 nominator    = NDF * G * F;
        //  0.001 to the denominator to prevent a divide by zero
        float denominator = 4 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0) + 0.001;
        specular          = nominator / denominator;

        // and add to kD
        vec3 kS = F;
        kD = vec3(1.0) - kS;

        kD *= 1.0 - metallic;

        float NdotL = max(dot(N, L), 0.0);
        Lo += (kD * albedo / PI + specular) * radiance * NdotL;

        ambient = vec3(0.0006) * albedo * ao;
    }
    // end foreach

    FragColor   = vec4(ambient + Lo, 1.0);
}

vec2 CalcTexCoord() {
   return gl_FragCoord.xy / gScreenSize;
}

vec3 WorldPosFromDepth(float depth, vec2 TexCoords) {
    float z = depth * 2.0 - 1.0;
    vec4 clipSpacePosition = vec4(TexCoords * 2.0 - 1.0, z, 1.0);
    vec4 viewSpacePosition = invProjection * clipSpacePosition;
    // Perspective division
    viewSpacePosition /= viewSpacePosition.w;
    vec4 worldSpacePosition = invView * viewSpacePosition;
    return worldSpacePosition.xyz;
}

// The Fresnel equation returns the ratio of light that gets reflected on a surface
vec3 fresnelSchlick(float cosTheta, vec3 F0)
{
    return F0 + (1.0 - F0) * pow(1.0 - cosTheta, 5.0);
}

float DistributionGGX(vec3 N, vec3 H, float roughness)
{
    float a      = roughness*roughness;
    float a2     = a*a;
    float NdotH  = max(dot(N, H), 0.0);
    float NdotH2 = NdotH*NdotH;

    float nom   = a2;
    float denom = (NdotH2 * (a2 - 1.0) + 1.0);
    denom = PI * denom * denom;

    return nom / denom;
}

float GeometrySchlickGGX(float NdotV, float roughness)
{
    float r = (roughness + 1.0);
    float k = (r*r) / 8.0;

    float nom   = NdotV;
    float denom = NdotV * (1.0 - k) + k;

    return nom / denom;
}
float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness)
{
    float NdotV = max(dot(N, V), 0.0);
    float NdotL = max(dot(N, L), 0.0);
    float ggx2  = GeometrySchlickGGX(NdotV, roughness);
    float ggx1  = GeometrySchlickGGX(NdotL, roughness);

    return ggx1 * ggx2;
}

