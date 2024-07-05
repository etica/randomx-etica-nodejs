package main

/*
#cgo CFLAGS: -I./include
#cgo linux LDFLAGS: -L./lib -lrandomx -lstdc++
#cgo darwin LDFLAGS: -L./lib -lrandomx -lstdc++
#cgo windows LDFLAGS: -L./lib -lrandomx -lstdc++ -lws2_32 -ladvapi32
#include <stdbool.h>
#include "randomx.h"
*/
import "C"
import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"unsafe"
)

// CheckSolution verifies a mining solution using RandomX

// VerifyEticaRandomXNonce verifies a mining solution using RandomX
//
//export VerifyEticaRandomXNonce
func VerifyEticaRandomXNonce(blockHeader *C.uchar, blockHeaderLength C.size_t,
	nonce *C.uchar, nonceLength C.size_t,
	target *C.uchar, targetLength C.size_t) C.bool {

	fmt.Println("*-*-*-*-**-*-*-*-*-*-Verifying with VerifyEticaRandomXNonce *-*-*-*-*-**-*-*-*-*-*-*-*-*-")

	// Convert C types to Go slices
	goBlockHeader := C.GoBytes(unsafe.Pointer(blockHeader), C.int(blockHeaderLength))
	goNonce := C.GoBytes(unsafe.Pointer(nonce), C.int(nonceLength))
	goTarget := C.GoBytes(unsafe.Pointer(target), C.int(targetLength))

	fmt.Printf("Block Header (hex): %s\n", hex.EncodeToString(goBlockHeader))
	fmt.Printf("Nonce (hex): %s\n", hex.EncodeToString(goNonce))
	fmt.Printf("Target (hex): %s\n", hex.EncodeToString(goTarget))

	// Initialize RandomX
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		fmt.Println("Failed to initialize RandomX cache")
		return C.bool(false)
	}
	defer DestroyRandomX(cache)

	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		fmt.Println("Failed to create RandomX VM")
		return C.bool(false)
	}
	defer DestroyVM(vm)

	input := append(goBlockHeader, goNonce...)
	fmt.Printf("Input for solution (hex): %s\n", hex.EncodeToString(input))

	const maxInputSize = 1024 * 1024 // 1 MB, adjust as needed
	if len(input) > maxInputSize {
		fmt.Printf("Input size too large: %d bytes\n", len(input))
		return C.bool(false)
	}

	correctSolution := CalculateHash(vm, input)
	fmt.Printf("Calculated solution (hex): %s\n", hex.EncodeToString(correctSolution))

	valid, err := CheckSolutionWithTarget(vm, goBlockHeader, goNonce, correctSolution, goTarget)
	if err != nil {
		fmt.Printf("RandomX verification error: %v\n", err)
		return C.bool(false)
	}

	if valid {
		fmt.Println("------***----------- RandomX verification passed -----------***-------")
		return C.bool(true)
	} else {
		fmt.Println("RandomX verification failed")
		return C.bool(false)
	}
}

func CheckSolutionWithTarget(vm unsafe.Pointer, blockHeader []byte, nonce []byte, solution []byte, target []byte) (bool, error) {
	if vm == nil {
		return false, errors.New("RandomX VM is not initialized")
	}

	input := append(blockHeader, nonce...)
	hash := CalculateHash(vm, input)

	if !bytes.Equal(hash, solution) {
		return false, errors.New("solution does not match calculated hash")
	}

	if bytes.Compare(hash, target) > 0 {
		return false, errors.New("hash does not meet target difficulty")
	}

	return true, nil
}

func main() {
	fmt.Println("Main function called!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
