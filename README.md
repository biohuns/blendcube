# BlendCube

![Build and Deploy](https://github.com/biohuns/blendcube/workflows/Build%20and%20Deploy/badge.svg)

Simple API Server for Generating Rubik's Cube 3D Model from URL

## Usage

```bash
make run
```

## Parameters

| Variable | Description        | Value Range   | Comment                        |
| -------- | ------------------ | ------------- | ------------------------------ |
| .        | extension          | .gltf \| .glb | e.g. `/cube.gltf?alg=U2+F2+R2` |
| alg      | algorithm to apply | [UDFBLR'2 ]\* |                                |
| is_unlit | unlit model flag   | true \| false |                                |
