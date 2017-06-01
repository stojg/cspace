#version 410 core

out vec3 Depth;

layout (std140) uniform Matrices
{
    mat4 projection;
    mat4 view;
    mat4 invProjection;
    mat4 invView;
    vec3 cameraPos;
};

uniform vec2 gScreenSize;
uniform sampler2D gDepth;

vec2 TexCoords = gl_FragCoord.xy / gScreenSize;

vec3 ViewPosFromDepth(float depth, vec2 TexCoords);

void main()
{
    vec2 pixelOffset = 1/gScreenSize;
    vec2 texCoord = vec2(TexCoords.x - 0.25f * pixelOffset.x, TexCoords.y - 0.25f * pixelOffset.y);
    Depth = ViewPosFromDepth(texture(gDepth, texCoord).r, TexCoords);
    return;
}

vec3 ViewPosFromDepth(float depth, vec2 TexCoords) {
    vec4 clipSpacePosition = vec4(TexCoords * 2.0 - 1.0, depth * 2.0 - 1.0, 1.0);
    vec4 viewSpacePosition = invProjection * clipSpacePosition;
    return (viewSpacePosition /= viewSpacePosition.w).xyz;
}
