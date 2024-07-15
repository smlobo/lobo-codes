cmake_minimum_required(VERSION 3.30)

project(h-tree)

add_definitions(-std=c17)
set(CMAKE_C_STANDARD 17)

add_executable(${PROJECT_NAME} src/main.c
        src/h-tree.c)

target_compile_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2)

target_link_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2 -sEXPORTED_FUNCTIONS=["_mainf"])

add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E create_symlink "${CMAKE_SOURCE_DIR}/h-tree.html" "${CMAKE_BINARY_DIR}/h-tree.html"
        DEPENDS h-tree.html
        VERBATIM ON
)

# Does not work with emcmake
install(TARGETS
        ${PROJECT_NAME}
)
