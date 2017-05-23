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
    vec3 Diffuse = texture(gAlbedoSpec, TexCoords).rgb;
    float Specular = texture(gAlbedoSpec, TexCoords).a;

    vec3 lightDir = dirLight.Direction;
    vec3 viewDir  = normalize(viewPos - FragPos);

    // ambient
    vec3 ambient = Diffuse * 0.001;

    // Diffuse
    vec3 diffuse = max(dot(Normal, lightDir), 0.0) * Diffuse * dirLight.Color;

    // Specular
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(Normal, halfwayDir), 0.0), 256);
    vec3 specular = dirLight.Color * spec * Specular;

    vec3 color = diffuse;
    FragColor = vec4(color, 1.0);
}
