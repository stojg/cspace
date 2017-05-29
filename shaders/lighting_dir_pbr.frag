#version 410 core

out vec4 FragColor;

uniform sampler2D gDepth;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;
uniform sampler2D gAmbientOcclusion;
uniform samplerCube irradianceMap;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
    mat4 invProjection;
    mat4 invView;
    vec3 cameraPos;
};

struct Light {
    vec3 Direction;
    vec3 Color;
};
uniform Light dirLight;

uniform vec2 gScreenSize;

uniform sampler2D shadowMap;
uniform mat4 lightProjection;
uniform mat4 lightView;

const float PI = 3.14159265359;

vec4 viewPosFromDepth(float depth, vec2 TexCoords);
vec3 fresnelSchlick(float cosTheta, vec3 F0);
vec3 fresnelSchlickRoughness(float cosTheta, vec3 F0, float roughness);
float DistributionGGX(vec3 N, vec3 H, float roughness);
float GeometrySchlickGGX(float NdotV, float roughness);
float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness);
float ShadowCalculation(vec4 worldPos, vec3 normal);
vec3 LightCalculation(vec3 V, vec3 N, vec3 albedo, float roughness, float metallic, vec3 F0, vec3 lightPos, vec3 lightColor, float attenuation);

void main()
{
    vec2 TexCoords = gl_FragCoord.xy / gScreenSize;
    vec4 FragPos   = viewPosFromDepth(texture(gDepth, TexCoords).x, TexCoords);

    // Normal
    vec3 N = normalize(texture(gNormal, TexCoords).rgb);
    // since this is a directional light in view space
    vec3 V = normalize(-FragPos.xyz);

    // sample all the textures
    vec3 albedo     = texture(gAlbedoSpec, TexCoords).rgb;
    float metallic  = texture(gAlbedoSpec, TexCoords).a;
    float roughness = texture(gNormal, TexCoords).w;
    float ao        = texture(gAmbientOcclusion, TexCoords).r;
    float shadow    = ShadowCalculation(invView * FragPos, N);

    vec3 Lo = vec3(0.0);

    // calculate reflectance at normal incidence; if dia-electric (like plastic) use F0 of 0.04 and if it's a metal,
    // use the albedo color as F0 (metallic workflow)
    vec3 F0 = vec3(0.04);
    F0      = mix(F0, albedo, metallic);

    // light calculations start here
    vec3 lightPos = transpose(mat3(invView)) * normalize(dirLight.Direction);
    Lo += LightCalculation(V, N, albedo, roughness, metallic, F0, lightPos, dirLight.Color, 1.0) * (1 - shadow);

    // ambient lighting (we now use IBL as the ambient term)
    vec3 kS = fresnelSchlickRoughness(max(dot(N, V), 0.0), F0, roughness);
    vec3 kD = 1.0 - kS;
    kD *= 1.0 - metallic;
    vec3 irradiance = texture(irradianceMap, N).rgb;
    vec3 diffuse      = irradiance * albedo;
    vec3 ambient = (kD * diffuse) * ao;

    FragColor   = vec4(ambient + Lo ,1);
}

vec3 LightCalculation(vec3 V, vec3 N, vec3 albedo, float roughness, float metallic, vec3 F0, vec3 lightPos, vec3 lightColor, float attenuation) {
    vec3 L = normalize(lightPos);
    vec3 H = normalize(V + L);

    vec3 radiance = lightColor * attenuation;

    // Cook-Torrance BRDF
    float NDF = DistributionGGX(N, H, roughness);
    float G   = GeometrySmith(N, V, L, roughness);
    vec3 F  = fresnelSchlick(max(dot(H, V), 0.0), F0);

    vec3 nominator    = NDF * G * F;
    float denominator = 4 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0) + 0.001; // 0.001 to the denominator to prevent a divide by zero
    vec3 specular     = nominator / denominator;

    // kS is equal to Fresnel
    vec3 kS = F;
    // for energy conservation, the diffuse and specular light can't be above 1.0 (unless the surface emits light); to
    // preserve this relationship the diffuse component (kD) should equal 1.0 - kS.
    vec3 kD = vec3(1.0) - kS;
    // multiply kD by the inverse metalness such that only non-metals
    // have diffuse lighting, or a linear blend if partly metal (pure metals
    // have no diffuse light).
    kD *= 1.0 - metallic;
     // scale light by NdotL
    float NdotL = max(dot(N, L), 0.0);
     // add to outgoing radiance Lo
    return (kD * albedo / PI + specular) * radiance * NdotL; // note that we already multiplied the BRDF by the Fresnel (kS) so we won't multiply by kS again
}

float ShadowCalculation(vec4 worldPos, vec3 normal)
{
    // This will be in clip-space
    vec4 lightSpacePos = lightProjection * lightView * worldPos;
    // Transform it into NDC-space by dividing by w (perspective divide)
    lightSpacePos /= lightSpacePos.w;
    // Range is now [-1.0, 1.0], but we need [0.0, 1.0]
    lightSpacePos = lightSpacePos * vec4 (0.5) + vec4 (0.5);

    // dont shadow things outside the light frustrum far plane
    if(lightSpacePos.z > 1.0) {
        return 0.0;
    }

    float shadow = 0.0;
    // change the amount of bias based on the surface angle towards the light
    float bias = max(0.001 * (1.0 - dot(normal, dirLight.Direction)), 0.001);
    // float bias = 0.001;
    vec2 texelSize = 0.5 / textureSize(shadowMap, 0);
    // Percentage Closing Filter
    for(int x = -1; x <= 1; ++x) {
        for(int y = -1; y <= 1; ++y) {
            float pcfDepth = texture(shadowMap, lightSpacePos.xy + vec2(x, y) * texelSize).r;
            shadow += lightSpacePos.z - bias > pcfDepth ? 1.0 : 0.0;
        }
    }
    return shadow /= 9.0;
}

vec4 viewPosFromDepth(float depth, vec2 TexCoords) {
    float z = depth * 2.0 - 1.0;
    vec4 clipSpacePosition = vec4(TexCoords * 2.0 - 1.0, z, 1.0);
    vec4 viewSpacePosition = invProjection * clipSpacePosition;
    return viewSpacePosition /= viewSpacePosition.w;
}

// The Fresnel equation returns the ratio of light that gets reflected on a surface
vec3 fresnelSchlick(float cosTheta, vec3 F0)
{
    return F0 + (1.0 - F0) * pow(1.0 - cosTheta, 5.0);
}

vec3 fresnelSchlickRoughness(float cosTheta, vec3 F0, float roughness)
{
    return F0 + (max(vec3(1.0 - roughness), F0) - F0) * pow(1.0 - cosTheta, 5.0);
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
