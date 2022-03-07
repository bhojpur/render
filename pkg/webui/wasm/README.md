# Bhojpur Render - WebAssembly Example

It shows method of building `WebAssembly` programs using `Go`.

## Pre-requisites

You need to install `tinygo` using following command (e.g. on macOS).

```bash
$ brew install tinygo
```

## Wasm Compilation

```bash
$ tinygo build -o main.wasm -target wasm main.go 
```
