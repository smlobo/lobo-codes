cmake_minimum_required(VERSION 3.30)

project(h-tree)

add_definitions(-std=c17)
set(CMAKE_C_STANDARD 17)

add_executable(${PROJECT_NAME} src/main.c
        src/h-tree.c)

target_compile_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2)

target_link_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2 -sEXPORTED_FUNCTIONS=["_mainf"])

add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E create_symlink "${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html" "${CMAKE_BINARY_DIR}/index.html"
        DEPENDS ${PROJECT_NAME}.html
        VERBATIM ON
)

install(FILES ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.js
        ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.wasm
        ${CMAKE_BINARY_DIR}/index.html
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
)
