#version 330 core
out vec4 FragColor;
in vec2 TexCoords;

uniform sampler2D screenTexture;

uniform bool horizontal;

uniform float weight[5] = float[] (0.2270270270, 0.1945945946, 0.1216216216, 0.0540540541, 0.0162162162);

void main()
{
     vec2 tex_offset = 1.0 / textureSize(screenTexture, 0); // gets size of single texel
     vec3 result = texture(screenTexture, TexCoords).rgb * weight[0];
     if(horizontal)
     {
        result += texture(screenTexture, TexCoords + vec2(tex_offset.x * 1, 0.0)).rgb * weight[1];
        result += texture(screenTexture, TexCoords - vec2(tex_offset.x * 1, 0.0)).rgb * weight[1];

        result += texture(screenTexture, TexCoords + vec2(tex_offset.x * 2, 0.0)).rgb * weight[2];
        result += texture(screenTexture, TexCoords - vec2(tex_offset.x * 2, 0.0)).rgb * weight[2];

        result += texture(screenTexture, TexCoords + vec2(tex_offset.x * 3, 0.0)).rgb * weight[3];
        result += texture(screenTexture, TexCoords - vec2(tex_offset.x * 3, 0.0)).rgb * weight[3];

        result += texture(screenTexture, TexCoords + vec2(tex_offset.x * 4, 0.0)).rgb * weight[4];
        result += texture(screenTexture, TexCoords - vec2(tex_offset.x * 4, 0.0)).rgb * weight[4];
     }
     else
     {
        result += texture(screenTexture, TexCoords + vec2(0.0, tex_offset.y * 1)).rgb * weight[1];
        result += texture(screenTexture, TexCoords - vec2(0.0, tex_offset.y * 1)).rgb * weight[1];

        result += texture(screenTexture, TexCoords + vec2(0.0, tex_offset.y * 2)).rgb * weight[2];
        result += texture(screenTexture, TexCoords - vec2(0.0, tex_offset.y * 2)).rgb * weight[2];

        result += texture(screenTexture, TexCoords + vec2(0.0, tex_offset.y * 3)).rgb * weight[3];
        result += texture(screenTexture, TexCoords - vec2(0.0, tex_offset.y * 3)).rgb * weight[3];

        result += texture(screenTexture, TexCoords + vec2(0.0, tex_offset.y * 4)).rgb * weight[4];
        result += texture(screenTexture, TexCoords - vec2(0.0, tex_offset.y * 4)).rgb * weight[4];
     }
     FragColor = vec4(result, 1.0);
}
