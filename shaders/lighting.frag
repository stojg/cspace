#version 410 core

out vec4 FragColor;

uniform sampler2D gPosition;
uniform sampler2D gNormal;
uniform sampler2D gAlbedoSpec;

struct Light {
    vec3 Position;
    vec3 Color;
    float DiffuseIntensity;
    float Linear;
    float Quadratic;
};

uniform Light pointLight;
uniform vec3 viewPos;
uniform vec2 gScreenSize;

vec2 CalcTexCoord() {
   return gl_FragCoord.xy / gScreenSize;
}

void main()
{
//    FragColor = vec4(1.0f);
//    return;
    vec2 TexCoords = CalcTexCoord();
    // Retrieve data from gbuffer
    vec3 FragPos = texture(gPosition, TexCoords).rgb;
    vec3 Normal = texture(gNormal, TexCoords).rgb;
    Normal = normalize(Normal);
    vec3 Diffuse = texture(gAlbedoSpec, TexCoords).rgb;
    float Specular = texture(gAlbedoSpec, TexCoords).a;



//    FragColor = vec4(Diffuse, 1.0f);
//    return;

    // Then calculate lighting as usual
    vec3 lighting  = Diffuse * 0.02; // hard-coded ambient component
    vec3 viewDir  = normalize(viewPos - FragPos);

    float distance = length(pointLight.Position - FragPos);

    // Diffuse
    vec3 lightDir = normalize(pointLight.Position - FragPos);
    vec3 diffuse = max(dot(Normal, lightDir), 0.0) * Diffuse * pointLight.Color;

    // Specular
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(Normal, halfwayDir), 0.0), 32.0);
    vec3 specular = pointLight.Color * spec * Specular;

    // Attenuation
    float attenuation = 1.0 / (1.0 + pointLight.Linear * distance + pointLight.Quadratic * distance * distance);
    diffuse *= attenuation * pointLight.DiffuseIntensity;
    specular *= attenuation;
    lighting += diffuse + specular;

    FragColor = vec4(lighting, 1.0);
}


