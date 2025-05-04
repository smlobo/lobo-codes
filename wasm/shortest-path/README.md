# Build Instructions

## Debug

* Use conan to generate user presets
```commandline
conan install . --output-folder=build-debug --profile=debug
```

* CMake configure with `emcmake` wrapper
```commandline
emcmake cmake --preset conan-debug
```

* CMake build
```commandline
cmake --build --preset conan-debug 
```

* Start a webserver & navigate to the *binary* dir `index.html`
```commandline
python3 -m http.server 9000
```
http://localhost:9000/build-debug/

## Release

* Use conan to generate user presets (default profile is `Release`)
```commandline
conan install . --output-folder=build-release
```

* CMake configure with `emcmake` wrapper
```commandline
emcmake cmake --preset conan-release
```

* CMake build
```commandline
cmake --build --preset conan-release
```

* Start a webserver & navigate to the *binary* dir `index.html`
```commandline
python3 -m http.server 9000
```
http://localhost:9000/build-release/

## Install

Install the `build-release` build to the current dir
```commandline
cmake --install build-release --prefix .
```

* Start a webserver & navigate to the *dist* dir
```commandline
python3 -m http.server 9000
```
http://localhost:9000/dist/shortest-path/
