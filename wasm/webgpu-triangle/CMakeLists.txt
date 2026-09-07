cmake_minimum_required(VERSION 3.30)
project(webgpu-triangle)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_executable(${PROJECT_NAME} src/main.cpp)

# Emscripten 6 uses the emdawnwebgpu port for the current WebGPU C API.
target_compile_options(${PROJECT_NAME} PRIVATE --use-port=emdawnwebgpu)

target_link_options(${PROJECT_NAME} PRIVATE
        -O2
        --use-port=emdawnwebgpu
        -sALLOW_MEMORY_GROWTH=1
)

add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
        COMMENT "Create an index.html link beside the WebAssembly output"
        COMMAND ${CMAKE_COMMAND} -E create_symlink
                "${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html"
                "${CMAKE_BINARY_DIR}/index.html"
        VERBATIM
)

install(FILES
        ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.js
        ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.wasm
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
)

install(FILES ${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
        RENAME index.html
)
