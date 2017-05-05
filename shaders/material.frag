#version 410
#define NR_POINT_LIGHTS 1

in vec3 FragNormal;
in vec3 FragPos;
in vec2 FragTexCoords;
in vec3 TangentLightPos;
in vec3 TangentViewPos;
in vec3 TangentFragPos;

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
    sampler2D depth0;
    float shininess;
};

uniform Light lights[NR_POINT_LIGHTS];
uniform Material mat;
uniform vec3 viewPos;

vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir, vec2 texCoords);

vec2 ParallaxMapping(vec2 texCoords, vec3 viewDir);

void main() {

    vec3 viewDir   = normalize(TangentViewPos - TangentFragPos);
//    vec2 texCoords = FragTexCoords;
    vec2 texCoords = ParallaxMapping(FragTexCoords,  viewDir);
//    vec3 normal = FragNormal;
    vec3 normal = texture(mat.normal0, texCoords).rgb;

    vec3 result = vec3(0,0,0);
    for(int i = 0; i < NR_POINT_LIGHTS; i++) {

        vec3 color =  texture(mat.diffuse0, FragTexCoords).rgb;
        // Ambient
        vec3 ambient = 0.1 * color;
        // Diffuse
        vec3 lightDir = normalize(TangentLightPos - TangentFragPos);
        float diff = max(dot(lightDir, normal), 0.0);
        vec3 diffuse = diff * color;
        // Specular
        vec3 viewDir = normalize(TangentViewPos - TangentFragPos);
        vec3 reflectDir = reflect(-lightDir, normal);
        vec3 halfwayDir = normalize(lightDir + viewDir);
        float spec = pow(max(dot(normal, halfwayDir), 0.0), 32.0);
        vec3 specular = vec3(0.2) * spec;

        // Attenuation
        float distance = length(lights[i].vector.xyz - TangentFragPos);
        float attenuation = 1.0f / (lights[i].constant + lights[i].linear * distance + lights[i].quadratic * (distance * distance));

        // Combine results
        //      vec3 ambient =  light.ambient  * vec3(texture(mat.diffuse0, texCoords));
        // vec3 specular = light.specular * spec * vec3(texture(mat.specular0, texCoords));

//        diffuse  *= attenuation;
//        specular *= attenuation;

//      return clamp(ambient + diffuse + specular, 0, 1);

        result = ambient + diffuse + specular;
    }
    color = vec4(result, 1.0f);
}

vec2 ParallaxMapping(vec2 texCoords, vec3 viewDir)
{
    float height_scale = 1.0f;
    float height =  texture(mat.depth0, texCoords).r;
    vec2 p = viewDir.xy / viewDir.z * (height * height_scale);
    return texCoords - p;
}

// Calculates the color when using a point light.
vec3 CalcPointLight(Light light, vec3 normal, vec3 fragPos, vec3 viewDir, vec2 texCoords)
{
    vec3 lightDir = normalize(light.vector.xyz - fragPos);

    // Specular shading - blinn-phong
    float diff = max(dot(normal, lightDir), 0.0);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), mat.shininess);

    // Attenuation
    float distance = length(light.vector.xyz - fragPos);
    float attenuation = 1.0f / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    // Combine results
    vec3 ambient =  light.ambient  * vec3(texture(mat.diffuse0, texCoords));
    vec3 diffuse =  light.diffuse  * diff * vec3(texture(mat.diffuse0, texCoords));
    vec3 specular = light.specular * spec * vec3(texture(mat.specular0, texCoords));

    ambient  *= attenuation;
    diffuse  *= attenuation;
    specular *= attenuation;

    return clamp(ambient + diffuse + specular, 0, 1);
}
