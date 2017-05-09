#version 330 core

out vec4 FragColor;

uniform sampler2D gPosition;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;
uniform sampler2D gDepth;

struct Light {
    vec3 Direction;
    vec3 Color;
};

uniform Light dirLight;
uniform vec3 viewPos;
uniform vec2 gScreenSize;

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

    vec3 lightDir = dirLight.Direction;
    vec3 viewDir  = normalize(viewPos - FragPos);

    // ambient
    vec3 ambient = Diffuse * 0.02;

    // Diffuse
    vec3 diffuse = max(dot(Normal, lightDir), 0.0) * Diffuse * dirLight.Color;

    // Specular
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(Normal, halfwayDir), 0.0), 32.0);
    vec3 specular = dirLight.Color * spec * Specular;

    vec3 color = ambient + diffuse + specular;
    FragColor = vec4(color, 1.0);
}
