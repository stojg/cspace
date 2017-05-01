#version 410

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

in vec3 Normal;
in vec3 FragPos;
in vec2 Frag_texture_coordinate;

out vec4 color;

uniform vec3 viewPos;

uniform Light light;

uniform sampler2D materialDiffuse;
uniform sampler2D materialSpecular;
uniform float materialShininess;

void main() {

    vec3 lightDir;
    float attenuation = 1.0f;
    if (light.vector.w == 0.0f) {
        // directional light
        lightDir = normalize(light.vector.xyz);
    } else {
        // point light
        lightDir = normalize(light.vector.xyz - FragPos);
        float distance = length(light.vector.xyz - FragPos);
        attenuation = 1.0f / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

    }

    // Ambient
    vec3 ambient = light.ambient * vec3(texture(materialDiffuse, Frag_texture_coordinate));

    // Diffuse
    vec3 norm = normalize(Normal);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = light.diffuse * diff * vec3(texture(materialDiffuse, Frag_texture_coordinate));

    // Specular
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), materialShininess);
    vec3 specular = light.specular * (spec * vec3(texture(materialSpecular, Frag_texture_coordinate)));

    ambient  *= attenuation;
    diffuse  *= attenuation;
    specular *= attenuation;

    color = vec4(ambient + diffuse + specular, 1.0f);
}

