#version 410

in vec3 Normal;
in vec3 FragPos;
in vec2 FragTexCoords;
in mat3 TBN;

out vec4 color;

struct Light {
    vec4 vector;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;

    // http://www.ogre3d.org/tikiwiki/tiki-index.php?page=-Point+Light+Attenuation
    float constant;
    float linear;
    float quadratic;
};

struct Material {
    sampler2D specular0;
    sampler2D diffuse0;
    sampler2D normal0;
    float shininess;
};

#define NR_POINT_LIGHTS 4
uniform Light lights[NR_POINT_LIGHTS];

uniform Material mat;

uniform vec3 viewPos;


vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir);

void main() {

    vec3 norm = texture(mat.normal0, FragTexCoords).rgb;
    norm = normalize(norm * 2.0 - 1.0);
    norm = normalize(TBN * norm);

    vec3 viewDir = normalize(viewPos - FragPos);

    vec3 result = vec3(0,0,0);
    for(int i = 0; i < NR_POINT_LIGHTS; i++) {
        result += CalcPointLight(lights[i], norm, FragPos, viewDir);
    }

    color = vec4(result, 1.0f);
}

// Calculates the color when using a point light.
vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir)
{
    vec3 lightDir = normalize(light.vector.xyz - fragPos);

    // Diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // Specular shading - blinn-phong
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), mat.shininess);

    // Attenuation
    float distance = length(light.vector.xyz - fragPos);
    float attenuation = 1.0f / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    // Combine results
    vec3 ambient = light.ambient * vec3(texture(mat.diffuse0, FragTexCoords));
    vec3 diffuse = light.diffuse * diff * vec3(texture(mat.diffuse0, FragTexCoords));
    vec3 specular = light.specular * (spec * vec3(texture(mat.specular0, FragTexCoords)));

    ambient  *= attenuation;
    diffuse  *= attenuation;
    specular *= attenuation;

    return (ambient + diffuse + specular);
}

