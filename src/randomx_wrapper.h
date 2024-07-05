#ifndef RANDOMX_WRAPPER_H
#define RANDOMX_WRAPPER_H

#ifdef __cplusplus
extern "C" {
#endif

void* InitRandomX(const unsigned char* key, size_t keyLength);

#ifdef __cplusplus
}
#endif

#endif // RANDOMX_WRAPPER_H