#version 410

in vec3 Normal;
in vec3 FragPos;
in vec2 Frag_texture_coordinate;

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

#define NR_POINT_LIGHTS 4
uniform Light lights[NR_POINT_LIGHTS];

uniform vec3 viewPos;
uniform sampler2D materialDiffuse;
uniform sampler2D materialSpecular;
uniform float materialShininess;

vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir);

void main() {

    vec3 norm = normalize(Normal);
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
    // Specular shading - phong
//    vec3 reflectDir = reflect(-lightDir, normal);
//    float spec = pow(max(dot(viewDir, reflectDir), 0.0), materialShininess);
    //  Specular shading - blinn-phong
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), materialShininess);

    // Attenuation
    float distance = length(light.vector.xyz - fragPos);
    float attenuation = 1.0f / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    // Combine results
    vec3 ambient = light.ambient * vec3(texture(materialDiffuse, Frag_texture_coordinate));
    vec3 diffuse = light.diffuse * diff * vec3(texture(materialDiffuse, Frag_texture_coordinate));
    vec3 specular = light.specular * (spec * vec3(texture(materialSpecular, Frag_texture_coordinate)));

    ambient  *= attenuation;
    diffuse  *= attenuation;
    specular *= attenuation;

    return (ambient + diffuse + specular);
}

