package main

/*
#cgo CFLAGS: -I./include
#cgo linux LDFLAGS: -L./lib -lrandomx -lstdc++
#cgo darwin LDFLAGS: -L./lib -lrandomx -lstdc++
#cgo windows LDFLAGS: -L./lib -lrandomx -lstdc++ -lws2_32 -ladvapi32
#include <stdlib.h>
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

const nonceOffset = 39

// CheckSolution verifies a mining solution using RandomX

// VerifyEticaRandomXNonce verifies a mining solution using RandomX
//
//export VerifyEticaRandomXNonce
func VerifyEticaRandomXNonce(blockHeader *C.uchar, blockHeaderLength C.size_t,
	nonce *C.uchar, nonceLength C.size_t,
	target *C.uchar, targetLength C.size_t,
	seedHash *C.uchar, seedHashLength C.size_t,
	expectedHash *C.uchar, expectedHashLength C.size_t) C.bool {

	fmt.Println("*-*-*-*-**-*-*-*-*-*-Verifying with VerifyEticaRandomXNonce *-*-*-*-*-**-*-*-*-*-*-*-*-*-")

	// Convert C types to Go slices
	goBlockHeader := C.GoBytes(unsafe.Pointer(blockHeader), C.int(blockHeaderLength))
	goNonce := C.GoBytes(unsafe.Pointer(nonce), C.int(nonceLength))
	goTarget := C.GoBytes(unsafe.Pointer(target), C.int(targetLength))
	goSeedHash := C.GoBytes(unsafe.Pointer(seedHash), C.int(seedHashLength))
	goExpectedHash := C.GoBytes(unsafe.Pointer(expectedHash), C.int(expectedHashLength))

	fmt.Printf("Block Header (hex): %s\n", hex.EncodeToString(goBlockHeader))
	fmt.Printf("Nonce (hex): %s\n", hex.EncodeToString(goNonce))
	fmt.Printf("Target (hex): %s\n", hex.EncodeToString(goTarget))
	fmt.Printf("Seed Hash (hex): %s\n", hex.EncodeToString(goSeedHash))
	fmt.Printf("Expected Hash (hex): %s\n", hex.EncodeToString(goExpectedHash))

	// Create a copy of the block header and insert the nonce at the correct offset
	blobWithNonce := make([]byte, len(goBlockHeader))
	copy(blobWithNonce, goBlockHeader)
	copy(blobWithNonce[nonceOffset:nonceOffset+4], goNonce)

	// Log the original blob and the blob with the new nonce
	fmt.Printf("Original Blob: %s\n", hex.EncodeToString(goBlockHeader))
	fmt.Printf("Blob with Nonce: %s\n", hex.EncodeToString(blobWithNonce))

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

	const maxInputSize = 1024 * 1024 // 1 MB, adjust as needed
	if len(blobWithNonce) > maxInputSize {
		fmt.Printf("Input blobWithNonce size too large: %d bytes\n", len(blobWithNonce))
		return C.bool(false)
	}

	calculatedHash := calculateRandomXHash(blobWithNonce, goSeedHash)
	fmt.Printf("Calculated RandomX Hash (hex): %s\n", hex.EncodeToString(calculatedHash))

	if !bytes.Equal(calculatedHash, goExpectedHash) {
		fmt.Printf("calculated hash does not match expected hash\n")
		return C.bool(false)
	}

	valid, err := CheckSolutionWithTarget(vm, blobWithNonce, calculatedHash, goTarget)
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

// Function to initialize RandomX cache and VM, and calculate the hash
func calculateRandomXHash(blobWithNonce, seedHash []byte) []byte {
	flags := FlagDefault
	cache := InitRandomX(flags)
	if cache == nil {
		panic("Failed to allocate RandomX cache")
	}
	defer DestroyRandomX(cache)

	InitCache(cache, seedHash)

	vm := CreateVM(cache, flags)
	if vm == nil {
		panic("Failed to create RandomX VM")
	}
	defer DestroyVM(vm)

	hash := CalculateHash(vm, blobWithNonce)

	return hash
}

func CheckSolutionWithTarget(vm unsafe.Pointer, blobWithNonce []byte, calculatedHash []byte, target []byte) (bool, error) {
	if vm == nil {
		return false, errors.New("RandomX VM is not initialized")
	}

	// Log original calculatedHash
	fmt.Printf("Calculated RandomX Hash (original): %s\n", hex.EncodeToString(calculatedHash))

	// Reverse the calculatedHash
	reversedHash := reverseBytes(calculatedHash)
	fmt.Printf("Calculated RandomX Hash (reversed for diff checking): %s\n", hex.EncodeToString(reversedHash))

	if bytes.Compare(reversedHash, target) > 0 {
		fmt.Println("Hash does not meet target difficulty (reversedHash)")
		return false, errors.New("hash does not meet target difficulty")
	}

	return true, nil
}

func reverseBytes(data []byte) []byte {
	reversed := make([]byte, len(data))
	for i := range data {
		reversed[i] = data[len(data)-1-i]
	}
	return reversed
}

/* func CheckSolutionWithTarget(vm unsafe.Pointer, blockHeader []byte, nonce []byte, solution []byte, target []byte) (bool, error) {
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
} */

func main() {
	fmt.Println("Main function called!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
