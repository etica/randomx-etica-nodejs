{
  "targets": [
    {
      "target_name": "randomx_addon",
      "sources": [ "randomx_node.cpp" ],
      "include_dirs": [
        "<!@(node -p \"require('node-addon-api').include\")",
        "src"
      ],
      "dependencies": [
        "<!(node -p \"require('node-addon-api').gyp\")"
      ],
      "cflags!": [ "-fno-exceptions" ],
      "cflags_cc!": [ "-fno-exceptions" ],
      "conditions": [
        ['OS=="linux"', {
          "libraries": [
            "-Wl,-rpath,'$$ORIGIN:$$ORIGIN/src'",
            "<!(pwd)/src/librandomx_wrapper.so",
            "<!(pwd)/src/lib/librandomx.a",
            "-lstdc++"
          ]
        }],
        ['OS=="mac"', {
          "libraries": [
            "-Wl,-rpath,@loader_path",
            "-L<!(pwd)/src",
            "-lrandomx_wrapper",
            "<!(pwd)/src/lib/librandomx.a",
            "-lstdc++"
          ],
          "xcode_settings": {
            "GCC_ENABLE_CPP_EXCEPTIONS": "YES",
            "CLANG_CXX_LIBRARY": "libc++",
            "MACOSX_DEPLOYMENT_TARGET": "10.7"
          }
        }],
        ['OS=="win"', {
          "libraries": [
            "<!(echo %cd%)/src/librandomx_wrapper.dll",
            "<!(echo %cd%)/src/lib/randomx.lib",
            "ws2_32.lib",
            "advapi32.lib"
          ]
        }]
      ]
    }
  ]
}