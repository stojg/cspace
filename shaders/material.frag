#version 410 core

in vec3 Normal;
in vec3 FragPos;
in vec2 Frag_texture_coordinate;

out vec4 color;

uniform vec3 viewPos;

uniform vec3 lightPos;
uniform vec3 lightAmbient;
uniform vec3 lightDiffuse;
uniform vec3 lightSpecular;

uniform sampler2D materialDiffuse1;
uniform vec3 materialSpecular;
uniform float materialShininess;

void main() {
    // Ambient
    vec3 ambient = lightAmbient * vec3(texture(materialDiffuse1, Frag_texture_coordinate));

    // Diffuse
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = lightDiffuse * diff * vec3(texture(materialDiffuse1, Frag_texture_coordinate));

    // Specular
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), materialShininess);
    vec3 specular = lightSpecular * (spec * materialSpecular);

    vec3 result = ambient + diffuse + specular;
    color = vec4(result, 1.0f);
}

