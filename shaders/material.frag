#version 410
#define NR_POINT_LIGHTS 2

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
    sampler2D specular0;
    sampler2D diffuse0;
    sampler2D normal0;
    float shininess;
};

uniform Light lights[NR_POINT_LIGHTS];
uniform Material mat;
uniform vec3 viewPos;
uniform float useNormalMapping;

vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir);

vec3 CalcBumpedNormal(vec3 normal, vec3 tangent);

void main() {
    vec3 norm = Normal;
    if(useNormalMapping > 0) {
        norm = CalcBumpedNormal(Normal, Tangent);
    }

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
    vec3 lightDir;

    if(light.vector.w > 0) {
        lightDir = normalize(light.vector.xyz - fragPos);
    } else {
        lightDir = normalize(-light.vector.xyz);
    }

    // Diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);

    // Specular shading - blinn-phong
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), mat.shininess);


    // Combine results
    vec3 ambient =  light.ambient  * vec3(texture(mat.diffuse0, FragTexCoords));
    vec3 diffuse =  light.diffuse  * diff * vec3(texture(mat.diffuse0, FragTexCoords));
    vec3 specular = light.specular * spec * vec3(texture(mat.specular0, FragTexCoords));

    if(light.vector.w > 0) {
        // Attenuation
        float distance = length(light.vector.xyz - fragPos);
        float attenuation = 1.0f / (light.constant + light.linear * distance + light.quadratic * (distance * distance));
        ambient  *= attenuation;
        diffuse  *= attenuation;
        specular *= attenuation;
    }

    return clamp(ambient + diffuse + specular, 0, 1);
}

vec3 CalcBumpedNormal(vec3 normal, vec3 tangent)
{
    vec3 Normal = normalize(normal);
    vec3 Tangent = normalize(tangent);
    Tangent = normalize(Tangent - dot(Tangent, Normal) * Normal);
    vec3 Bitangent = cross(Tangent, Normal);
    vec3 BumpMapNormal = texture(mat.normal0, FragTexCoords).xyz;
    BumpMapNormal = 2.0 * BumpMapNormal - vec3(1.0, 1.0, 1.0);
    vec3 NewNormal;
    mat3 TBN = mat3(Tangent, Bitangent, Normal);
    NewNormal = TBN * BumpMapNormal;
    NewNormal = normalize(NewNormal);
    return NewNormal;
}

