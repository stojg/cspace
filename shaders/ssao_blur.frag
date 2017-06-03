#version 330 core
out vec4 FragColor;
in vec2 TexCoords;

uniform sampler2D screenTexture;

uniform bool horizontal;

uniform float weight[3] = float[] (0.2270270270, 0.3162162162, 0.0702702703 );
uniform float offset[3] = float[]( 0.0, 1.3846153846, 3.2307692308 );

void main()
{
     vec2 tex_offset = 1.0 / textureSize(screenTexture, 0); // gets size of single texel
     vec3 result = texture(screenTexture, TexCoords).rgb * weight[0];
     if(horizontal)
     {

        for(int i = 1; i<3; i++) {
            vec2 off = vec2(tex_offset.x * offset[i], 0.0);
            result += texture(screenTexture, TexCoords + off).rgb * weight[i];
            result += texture(screenTexture, TexCoords - off).rgb * weight[i];
        }
     }
     else
     {
        for(int i = 1; i<3; i++) {
            vec2 off = vec2(0.0, tex_offset.y * offset[i]);
            result += texture(screenTexture, TexCoords + off).rgb * weight[i];
            result += texture(screenTexture, TexCoords - off).rgb * weight[i];
        }
     }
     FragColor = vec4(result, 1.0);
}
