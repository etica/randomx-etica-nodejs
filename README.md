
Build radnomX library librandomx.a:
On another folder than this project
clone randomX repo https://github.com/tevador/RandomX
mkdir build && cd build
sudo apt-get update
sudo apt-get install build-essential cmake
cmake ..
make

add new functions Process:
> Implement new function in /src main go package
> Then add them to randomx_node.cpp following same pattern as initRandomX and VerifyEticaRandomXNonce


** Installation** 

copy paste librandomx.a in this project /lib
Then run:
> go build -o etica_randomx_checker ./src

create a C shared library:
> go build -buildmode=c-shared -o librandomx.so randomx_wrapper.go etica_verification.go 


Build the addon:
> npm install
> node example_usage.js


> cd /src
> go build -buildmode=c-shared -o librandomx_wrapper.so
> go mod init randomx_wrapper
> go mod tidy
> go build -buildmode=c-shared -o librandomx_wrapper.so


> cd /src
> go build -buildmode=c-shared -o librandomx_wrapper.so randomx_wrapper.go etica_verification.go
> cd ..
> npm install
> node example_usage.js