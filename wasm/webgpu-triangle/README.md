# WebGPU Triangle

A minimal C++ WebGPU example for Emscripten's current `emdawnwebgpu` port. It draws one triangle with an inline WGSL shader.

## Build

Set `EMSCRIPTEN` to the Emscripten installation directory containing `emcc` and
`cmake/Modules/Platform/Emscripten.cmake`. The supplied CMake presets are
recognized by CLion and standard CMake tooling. They also use a local,
gitignored `.emscripten-cache` directory, avoiding write access to a
package-manager-owned Emscripten cache.

```sh
export EMSCRIPTEN=/path/to/emscripten
cmake --preset emscripten-debug
cmake --build --preset emscripten-debug
python3 -m http.server 9000
```

Then open [http://localhost:9000/build-debug/](http://localhost:9000/build-debug/). Use a browser with WebGPU support enabled.

In CLion, open this directory and select the `emscripten-debug` or
`emscripten-release` CMake profile. Set the same `EMSCRIPTEN` environment
variable in the profile's CMake options/environment settings.
