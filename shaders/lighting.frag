#version 410 core
#define NR_LIGHTS 32

out vec4 FragColor;
in vec2 TexCoords;

uniform sampler2D gPosition;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;

struct Light {
    vec3 Position;
    vec3 Color;
    float Radius;
    float DiffuseIntensity;
    float Linear;
    float Quadratic;
};

uniform Light lights[NR_LIGHTS];
uniform vec3 viewPos;

vec4 CalcPointLight(Light light, vec3 fragPos, vec3 normal, vec3 lightDir, float distance);

vec4 CalcLightInternal(Light light,  vec3 LightDirection, vec3 WorldPos, vec3 Normal);

void main()
{
    // Retrieve data from gbuffer
    vec3 FragPos = texture(gPosition, TexCoords).rgb;
    vec3 Normal = texture(gNormal, TexCoords).rgb;
    vec3 Diffuse = texture(gAlbedoSpec, TexCoords).rgb;
    float Specular = texture(gAlbedoSpec, TexCoords).a;

    // Then calculate lighting as usual
    vec3 lighting  = Diffuse * 0.02; // hard-coded ambient component
    vec3 viewDir  = normalize(viewPos - FragPos);

    for(int i = 0; i < NR_LIGHTS; ++i)
    {
        float distance = length(lights[i].Position - FragPos);
        if(distance > lights[i].Radius) {
            continue;
        }

        // Diffuse
        vec3 lightDir = normalize(lights[i].Position - FragPos);
        vec3 diffuse = max(dot(Normal, lightDir), 0.0) * Diffuse * lights[i].Color;

        // Specular
        vec3 halfwayDir = normalize(lightDir + viewDir);
        float spec = pow(max(dot(Normal, halfwayDir), 0.0), 16.0);
        vec3 specular = lights[i].Color * spec * Specular;

        // Attenuation
        float attenuation = 1.0 / (1.0 + lights[i].Linear * distance + lights[i].Quadratic * distance * distance);
        diffuse *= attenuation * lights[i].DiffuseIntensity;
        specular *= attenuation * lights[i].DiffuseIntensity;
        lighting += diffuse + specular;
    }
    FragColor = vec4(lighting, 1.0);
}

vec4 CalcPointLight(Light light, vec3 fragPos, vec3 normal, vec3 lightDir, float distance) {

    vec4 Color = CalcLightInternal(light, lightDir, fragPos, normal);
    float att = 1.0f + light.Linear * distance + light.Quadratic * distance * distance;
    att = max(1.0, att);
    return Color / att;
}

vec4 CalcLightInternal(Light light,  vec3 LightDirection, vec3 WorldPos, vec3 Normal) {

    vec3 dColor = texture(gAlbedoSpec, TexCoords).rgb;
    float gSpecularPower = 0;
    float gMatSpecularIntensity = 0;

    vec4 AmbientColor = vec4(light.Color * 0.01 * dColor, 1.0);
    float DiffuseFactor = dot(Normal, -LightDirection);

    vec4 DiffuseColor  = vec4(0, 0, 0, 0);
    vec4 SpecularColor = vec4(0, 0, 0, 0);

    if (DiffuseFactor > 0.0) {
        DiffuseColor = vec4(light.Color * light.DiffuseIntensity * DiffuseFactor * dColor, 1.0);
        vec3 VertexToEye = normalize(viewPos - WorldPos);
        vec3 LightReflect = normalize(reflect(LightDirection, Normal));
        float SpecularFactor = dot(VertexToEye, LightReflect);
        if (SpecularFactor > 0.0) {
            float sColor = texture(gAlbedoSpec, TexCoords).a;
            SpecularFactor = pow(SpecularFactor, gSpecularPower);
            SpecularColor = vec4(light.Color * gMatSpecularIntensity * SpecularFactor * sColor, 1.0);
        }
    }

    return (AmbientColor + DiffuseColor + SpecularColor);
}



