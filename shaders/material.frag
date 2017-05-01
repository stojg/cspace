#version 410

in vec3 Normal;
in vec3 FragPos;
in vec2 Frag_texture_coordinate;

out vec4 color;

uniform vec3 viewPos;

uniform vec3 lightPos;
uniform vec3 lightAmbient;
uniform vec3 lightDiffuse;
uniform vec3 lightSpecular;

uniform sampler2D materialDiffuse;
uniform sampler2D materialSpecular;
uniform float materialShininess;

void main() {
    // Ambient
    vec3 ambient = lightAmbient * vec3(texture(materialDiffuse, Frag_texture_coordinate));

    // Diffuse
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = lightDiffuse * diff * vec3(texture(materialDiffuse, Frag_texture_coordinate));

    // Specular
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), materialShininess);
    vec3 specular = lightSpecular * (spec * vec3(texture(materialSpecular, Frag_texture_coordinate)));

    float distance = length(lightPos - FragPos);
    float attenuation = 1.0f / (1.0f + 0.07f * distance + 0.017f * (distance * distance));

    ambient  *= attenuation;
    diffuse  *= attenuation;
    specular *= attenuation;

    color = vec4(ambient + diffuse + specular, 1.0f);
}

