#version 330 core
// lightning_point.frag

out vec4 FragColor;

uniform sampler2D gDepth;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;

struct Light {
    vec3 Position;
    vec3 Color;
    float Linear;
    float Quadratic;
};

uniform Light pointLight;
uniform vec3 viewPos;
uniform vec2 gScreenSize;
uniform mat4 projMatrixInv;
uniform mat4 viewMatrixInv;

in vec3 Position;

vec2 CalcTexCoord() {
   return gl_FragCoord.xy / gScreenSize;
}

vec3 WorldPosFromDepth(float depth, vec2 TexCoords) {
    float z = depth * 2.0 - 1.0;
    vec4 clipSpacePosition = vec4(TexCoords * 2.0 - 1.0, z, 1.0);
    vec4 viewSpacePosition = projMatrixInv * clipSpacePosition;
    // Perspective division
    viewSpacePosition /= viewSpacePosition.w;
    vec4 worldSpacePosition = viewMatrixInv * viewSpacePosition;
    return worldSpacePosition.xyz;
}

void main()
{
    vec2 TexCoords = CalcTexCoord();

    float depth = texture(gDepth, TexCoords).x;

    vec3 FragPos = WorldPosFromDepth(depth, TexCoords);

    vec3 Normal = texture(gNormal, TexCoords).rgb;
    vec4 albedo = texture(gAlbedoSpec, TexCoords);
    vec3 Diffuse = albedo.rgb;
    float Specular = albedo.a;

    vec3 lighting  = vec3(0.0);
    vec3 viewDir  = normalize(viewPos - FragPos);

    float distance = length(pointLight.Position - FragPos);

    // Diffuse
    vec3 lightDir = normalize(pointLight.Position - FragPos);
    vec3 diffuse = max(dot(Normal, lightDir), 0.0) * Diffuse * pointLight.Color;

    // Specular
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(Normal, halfwayDir), 0.0), 64.0);
    vec3 specular = pointLight.Color * spec * Specular;

    // Attenuation
    float attenuation = 1.0 / (1.0 + pointLight.Linear * distance + pointLight.Quadratic * distance * distance);
    diffuse *= attenuation;
    specular *= attenuation;
    // ambient
    lighting += diffuse + specular;
    FragColor = vec4(lighting, 1.0);
}



