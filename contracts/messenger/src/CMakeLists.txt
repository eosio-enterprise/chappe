project(messenger)

set(EOSIO_WASM_OLD_BEHAVIOR "Off")
find_package(eosio.cdt)

add_contract( messenger messenger messenger.cpp )
target_include_directories( messenger PUBLIC ${CMAKE_SOURCE_DIR}/../include )
