cmake_minimum_required(VERSION 3.30)
project(collision-system)

set(CMAKE_CXX_STANDARD 17)

add_executable(${PROJECT_NAME} src/main.cpp
        src/collision-system.cpp
        src/collision-system.h
        src/main.h
        src/ball.cpp
        src/ball.h
        src/sdl-helpers/circle.cpp
        src/sdl-helpers/circle.h
        src/event.cpp
        src/event.h
)

target_compile_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2)

target_link_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2 -sEXPORTED_FUNCTIONS=["_mainf"]
        -sASYNCIFY -sTOTAL_MEMORY=512MB)

add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E create_symlink "${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html" "${CMAKE_BINARY_DIR}/index.html"
        DEPENDS ${PROJECT_NAME}.html
        VERBATIM ON
)

install(FILES ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.js
        ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.wasm
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
)

install(FILES ${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
        RENAME index.html
)
