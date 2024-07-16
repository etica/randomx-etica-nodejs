#include <node.h>
#include <node_buffer.h>
#include "src/randomx.h"
#include "src/randomx_wrapper.h"

namespace randomx_addon {

using v8::FunctionCallbackInfo;
using v8::Isolate;
using v8::Local;
using v8::Object;
using v8::String;
using v8::Value;
using v8::Exception;

// Declare the Go functions
extern "C" {
    void* InitRandomX(const unsigned char* key, size_t keyLength);
    bool VerifyEticaRandomXNonce(const unsigned char* blockHeader, size_t blockHeaderLength,
                                 const unsigned char* nonce, size_t nonceLength,
                                 const unsigned char* target, size_t targetLength,
                                 const unsigned char* seedHash, size_t seedHashLength);
}

void* randomxCache = nullptr;
void* randomxVM = nullptr;

void NodeInitRandomX(const FunctionCallbackInfo<Value>& args) {
    printf("NodeInitRandomX called\n");
    Isolate* isolate = args.GetIsolate();

    if (args.Length() < 1 || !node::Buffer::HasInstance(args[0])) {
        isolate->ThrowException(Exception::TypeError(
            String::NewFromUtf8(isolate, "Wrong arguments").ToLocalChecked()));
        return;
    }

    Local<Object> bufferObj = args[0]->ToObject(isolate->GetCurrentContext()).ToLocalChecked();
    char* bufferData = node::Buffer::Data(bufferObj);
    size_t bufferLength = node::Buffer::Length(bufferObj);

    randomxCache = InitRandomX(reinterpret_cast<const unsigned char*>(bufferData), bufferLength);
    
    if (randomxCache == nullptr) {
         printf("InitRandomX failed: cache is null\n");
        args.GetReturnValue().Set(false);
    } else {
        // Here you might want to create the VM if needed
        // randomxVM = CreateVM(randomxCache);
        printf("InitRandomX succeeded\n");
        args.GetReturnValue().Set(true);
    }
}


void VerifyEticaRandomXNonce(const FunctionCallbackInfo<Value>& args) {
    Isolate* isolate = args.GetIsolate();

    if (args.Length() < 4 || !node::Buffer::HasInstance(args[0]) || 
        !node::Buffer::HasInstance(args[1]) || !node::Buffer::HasInstance(args[2]) ||
        !node::Buffer::HasInstance(args[3])) {
        isolate->ThrowException(Exception::TypeError(
            String::NewFromUtf8(isolate, "Wrong arguments").ToLocalChecked()));
        return;
    }

    Local<Object> blockHeaderObj = args[0]->ToObject(isolate->GetCurrentContext()).ToLocalChecked();
    char* blockHeaderData = node::Buffer::Data(blockHeaderObj);
    size_t blockHeaderLength = node::Buffer::Length(blockHeaderObj);

    Local<Object> nonceObj = args[1]->ToObject(isolate->GetCurrentContext()).ToLocalChecked();
    char* nonceData = node::Buffer::Data(nonceObj);
    size_t nonceLength = node::Buffer::Length(nonceObj);

    Local<Object> targetObj = args[2]->ToObject(isolate->GetCurrentContext()).ToLocalChecked();
    char* targetData = node::Buffer::Data(targetObj);
    size_t targetLength = node::Buffer::Length(targetObj);

    Local<Object> seedHashObj = args[3]->ToObject(isolate->GetCurrentContext()).ToLocalChecked();
    char* seedHashData = node::Buffer::Data(seedHashObj);
    size_t seedHashLength = node::Buffer::Length(seedHashObj);

    bool result = VerifyEticaRandomXNonce(
        reinterpret_cast<const unsigned char*>(blockHeaderData), blockHeaderLength,
        reinterpret_cast<const unsigned char*>(nonceData), nonceLength,
        reinterpret_cast<const unsigned char*>(targetData), targetLength,
        reinterpret_cast<const unsigned char*>(seedHashData), seedHashLength);

    args.GetReturnValue().Set(result);
}

void Initialize(Local<Object> exports) {
    NODE_SET_METHOD(exports, "InitRandomX", NodeInitRandomX);
    NODE_SET_METHOD(exports, "VerifyEticaRandomXNonce", VerifyEticaRandomXNonce);
}

NODE_MODULE(NODE_GYP_MODULE_NAME, Initialize)

}  // namespace randomx_addon