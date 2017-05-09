#version 330 core

out vec4 FragColor;

uniform sampler2D gPosition;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;
uniform sampler2D gDepth;

struct Light {
    vec3 Position;
    vec3 Color;
    float Linear;
    float Quadratic;
};

uniform Light pointLight;
uniform vec3 viewPos;
uniform vec2 gScreenSize;

in vec3 PositionVS;

float far = 200;
float near = 0.1;

vec2 CalcTexCoord() {
   return gl_FragCoord.xy / gScreenSize;
}

void main()
{
    vec2 TexCoords = CalcTexCoord();
    vec3 FragPos = texture(gPosition, TexCoords).rgb;
    vec3 Normal = texture(gNormal, TexCoords).rgb;
    vec3 Diffuse = texture(gAlbedoSpec, TexCoords).rgb;
    float Specular = texture(gAlbedoSpec, TexCoords).a;

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
    lighting += diffuse + specular;

    FragColor = vec4(lighting, 1.0);
}


