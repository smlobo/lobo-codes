cmake_minimum_required(VERSION 3.30)
project(shortest-path)

set(CMAKE_CXX_STANDARD 17)

add_executable(${PROJECT_NAME}
        src/main.cpp
        src/edge-weighted-digraph.cpp
        src/vertex.cpp
        src/edge.cpp
        src/sdl-helpers/circle.cpp
)

target_compile_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2)

target_link_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2 -sEXPORTED_FUNCTIONS=["_mainf"]
        -sASYNCIFY -sTOTAL_MEMORY=512MB)

add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
        COMMENT "Symbolic link for index.html to ${PROJECT_NAME}.html"
        COMMAND ${CMAKE_COMMAND} -E create_symlink "${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html" "${CMAKE_BINARY_DIR}/index.html"
)

install(FILES ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.js ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.wasm
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
)

install(FILES ${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
        RENAME index.html
)

