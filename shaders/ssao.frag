#version 330 core

out float FragColor;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
    mat4 invProjection;
    mat4 invView;
    vec3 cameraPos;
};

uniform sampler2D gDepth;
uniform sampler2D gNormal;
uniform sampler2D texNoise;

uniform vec2 gScreenSize;

uniform int enabled = 0;

uniform vec3 samples[64];

// parameters (you'd probably want to use them as uniforms to more easily tweak the effect)
const int kernelSize = 64;
const float radius = 0.9;
const float bias = 0.025;

// tile noise texture over screen based on screen dimensions divided by noise size
vec2 noiseScale = gScreenSize / 4.0;

vec2 TexCoords = gl_FragCoord.xy / gScreenSize;

void main() {

    if(enabled != 1.0) {
        FragColor = 1.0;
        return;
    }

    vec3 fragPos = texture(gDepth, TexCoords).xyz;

    vec3 normal = texture(gNormal, TexCoords).rgb;
    vec3 randomVec = normalize(texture(texNoise, TexCoords * noiseScale).xyz);
    // create TBN change-of-basis matrix: from tangent-space to view-space
    vec3 tangent = normalize(randomVec - normal * dot(randomVec, normal));
    vec3 bitangent = cross(normal, tangent);
    mat3 TBN = mat3(tangent, bitangent, normal);

    // iterate over the sample kernel and calculate occlusion factor
    float occlusion = 0.0;

    for(int i = 0; i < kernelSize; ++i) {
        // get sample position
        vec3 sampl = TBN * samples[i]; // from tangent to view-space
        sampl = sampl * radius + fragPos;
        // project sample position:
        vec4 offset = vec4(sampl, 1.0);
        offset = projection * offset;
        offset.xy /= offset.w;
        offset.xy = offset.xy * 0.5 + vec2(0.5);

        float sampleDepth = texture(gDepth, offset.xy).z;
        float rangeCheck = smoothstep(0.0, 1.0, radius / abs(fragPos.z - sampleDepth));
        occlusion       += (sampleDepth >= sampl.z + bias ? 1.0 : 0.0) * rangeCheck;
    }
    FragColor = 1.0 - (occlusion / kernelSize);
}

