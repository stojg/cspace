#version 410
#define NR_POINT_LIGHTS 1

in vec3 Normal;
in vec3 FragPos;
in vec2 FragTexCoords;
in vec3 Tangent;

out vec4 color;

struct Light {
    vec4 vector;
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float constant;
    float linear;
    float quadratic;
};

struct Material {
    sampler2D diffuse0;
    float shininess;
};

uniform Light lights[NR_POINT_LIGHTS];
uniform Material mat;
uniform vec3 viewPos;

vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir);

void main() {
    vec3 norm = Normal;
    vec3 viewDir = normalize(viewPos - FragPos);

    vec3 result = vec3(0,0,0);
    for(int i = 0; i < NR_POINT_LIGHTS; i++) {
        result += CalcPointLight(lights[i], norm, FragPos, viewDir);
    }
//    color = vec4(0.1f,0.5f,0.2f, 1.0f);
    color = vec4(result, 1.0f);
}

// Calculates the color when using a point light.
vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir)
{
    vec3 lightDir = normalize(light.vector.xyz - fragPos);

    // Diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // Attenuation
    float distance = length(light.vector.xyz - fragPos);
    float attenuation = 1.0f / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    // Combine results
    vec3 ambient =  vec3(0.02f);
    vec3 diffuse =  light.diffuse  * diff * vec3(0.1f);

    diffuse  *= attenuation;

    return (ambient + diffuse);
}

