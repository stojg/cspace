#version 330 core

in vec2 TexCoords;
in vec2 ViewRay;

out vec4 FragColor;

uniform sampler2D screenTexture;

float near = 1;
float far = 100;

float S = -1.020202; // 2,2
float T = -2.020202; // 2.3

float CalcViewZ(vec2 Coords);

void main()
{

    S = (-near-far)/(near-far);
    T = (2*far*near)/(near-far);
    float ViewZ = CalcViewZ(TexCoords);
    float ViewX = ViewRay.x * ViewZ;
    float ViewY = ViewRay.y * ViewZ;

    vec3 Pos = vec3(ViewX, ViewY, ViewZ);
    FragColor = vec4(Pos, 1.0f);

//    float Depth = texture(screenTexture, TexCoords).x;
}

float CalcViewZ(vec2 Coords)
{
    float Depth = texture(screenTexture, TexCoords).x;
    float ViewZ = T / (2 * Depth -1 - S);
    return ViewZ;
}


