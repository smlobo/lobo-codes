cmake_minimum_required(VERSION 3.30)
project(shortest-path)

set(CMAKE_CXX_STANDARD 17)

add_executable(${PROJECT_NAME}
        src/main.cpp
        src/edge-weighted-digraph.cpp
        src/vertex.cpp
        src/edge.cpp
        src/dijkstra-shortest-path.cpp
        src/sdl-helpers/circle.cpp
)

target_compile_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2 -sUSE_SDL_TTF=2)

target_link_options(${PROJECT_NAME} PRIVATE -O2 -sUSE_SDL=2 -sEXPORTED_FUNCTIONS=["_mainf"]
        -sASYNCIFY -sTOTAL_MEMORY=512MB -sUSE_SDL_TTF=2 --preload-file fonts)

add_custom_command(TARGET ${PROJECT_NAME} PRE_LINK
        COMMENT "Pre link symbolic link for fonts to static/fonts"
        COMMAND ${CMAKE_COMMAND} -E create_symlink "${CMAKE_SOURCE_DIR}/static/fonts" "${CMAKE_BINARY_DIR}/fonts"
        DEPENDS static/fonts
        VERBATIM ON
)

add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
        COMMENT "Post build symbolic link for index.html to ${PROJECT_NAME}.html"
        COMMAND ${CMAKE_COMMAND} -E create_symlink "${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html" "${CMAKE_BINARY_DIR}/index.html"
        DEPENDS ${PROJECT_NAME}.html
        VERBATIM ON
)

install(FILES ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.js ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.wasm
        ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.data ${CMAKE_BINARY_DIR}/fonts
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
)

install(FILES ${CMAKE_SOURCE_DIR}/${PROJECT_NAME}.html
        CONFIGURATIONS Release
        DESTINATION dist/${PROJECT_NAME}
        RENAME index.html
)

