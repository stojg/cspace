#version 410 core

out vec4 FragColor;

uniform sampler2D gDepth;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;

struct Light {
    vec3 Direction;
    vec3 Color;


};
uniform Light dirLight;

uniform vec3 viewPos;

uniform vec2 gScreenSize;
uniform mat4 projMatrixInv;
uniform mat4 viewMatrixInv;

uniform sampler2D shadowMap;
uniform mat4 lightProjection;
uniform mat4 lightView;

uniform float ao = 0.0;

const float PI = 3.14159265359;

vec2 CalcTexCoord();
vec4 WorldPosFromDepth(float depth, vec2 TexCoords);
vec3 fresnelSchlick(float cosTheta, vec3 F0);
float DistributionGGX(vec3 N, vec3 H, float roughness);
float GeometrySchlickGGX(float NdotV, float roughness);
float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness);
float ShadowCalculation(vec4 fragPosLightSpace);

void main()
{
    vec2 TexCoords = CalcTexCoord();
    float depth = texture(gDepth, TexCoords).x;
    vec4 FragPos = WorldPosFromDepth(depth, TexCoords);

    vec3 N = normalize(texture(gNormal, TexCoords).rgb);
    vec3 V = normalize(viewPos - FragPos.xyz);

    // shadow calc
    // This will be in clip-space
    vec4 lightSpacePos = lightProjection * lightView * FragPos;
    // Transform it into NDC-space by dividing by w
    lightSpacePos /= lightSpacePos.w;
    // Range is now [-1.0, 1.0], but we need [0.0, 1.0]
    lightSpacePos = lightSpacePos * vec4 (0.5) + vec4 (0.5);

    float shadow = 0.0;
    float closestDepth = texture(shadowMap, lightSpacePos.xy).r;
    float currentDepth = lightSpacePos.z;
    if(currentDepth < 1.0) {
    float bias = max(0.05 * (1.0 - dot(N, dirLight.Direction)), 0.005);
        shadow = currentDepth - bias > closestDepth  ? 1.0 : 0.0;
    }

    vec3 Lo = vec3(0.0);

    vec3 albedo = texture(gAlbedoSpec, TexCoords).rgb;
    float metallic = texture(gAlbedoSpec, TexCoords).a;
    float roughness = texture(gNormal, TexCoords).w;

    // foreach here
    vec3 L = dirLight.Direction;
    vec3 H = normalize(V + L);

    vec3 radiance = dirLight.Color;

    vec3 F0 = vec3(0.04);
    F0      = mix(F0, albedo, metallic);
    vec3 F  = fresnelSchlick(max(dot(H, V), 0.0), F0);

    float NDF = DistributionGGX(N, H, roughness);
    float G   = GeometrySmith(N, V, L, roughness);

    // Cook-Torrance BRDF
    vec3 nominator    = NDF * G * F;
    //  0.001 to the denominator to prevent a divide by zero
    float denominator = 4 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0) + 0.001;
    vec3 specular     = nominator / denominator;

    // and add to kD
    vec3 kS = F;
    vec3 kD = vec3(1.0) - kS;

    kD *= (1 - shadow);

    kD *= 1.0 - metallic;

    float NdotL = max(dot(N, L), 0.0);
    Lo += (kD * albedo / PI + specular) * radiance * NdotL;
    // end foreach

    // improvised ambient term
    vec3 ambient = vec3(0.03) * albedo * ao;
    FragColor   = vec4(ambient + Lo,1);
}

float ShadowCalculation(vec4 fragPosLightSpace)
{
    // perform perspective divide
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;

    projCoords = projCoords * 0.5 + 0.5;
    float closestDepth = texture(shadowMap, projCoords.xy).r;
    float currentDepth = projCoords.z;
    float shadow = currentDepth > closestDepth  ? 1.0 : 0.0;
    return shadow;
}

vec2 CalcTexCoord() {
   return gl_FragCoord.xy / gScreenSize;
}

vec4 WorldPosFromDepth(float depth, vec2 TexCoords) {
    float z = depth * 2.0 - 1.0;
    vec4 clipSpacePosition = vec4(TexCoords * 2.0 - 1.0, z, 1.0);
    vec4 viewSpacePosition = projMatrixInv * clipSpacePosition;
    // Perspective division
    viewSpacePosition /= viewSpacePosition.w;
    vec4 worldSpacePosition = viewMatrixInv * viewSpacePosition;
    return worldSpacePosition;
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
