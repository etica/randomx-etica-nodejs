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

var globalRandomXCache unsafe.Pointer
var globalRandomXVM unsafe.Pointer
var globalSeedHash []byte

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

	// Initialize RandomX system if needed
	if globalRandomXCache == nil || globalRandomXVM == nil || !bytes.Equal(globalSeedHash, goSeedHash) {
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!! Initializing RandomX system due to missing or outdated cache/VM. !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		if err := initRandomXSystem(FlagDefault, goSeedHash); err != nil {
			fmt.Printf("Error initializing RandomX system: %v\n", err)
			return C.bool(false)
		}
		globalSeedHash = goSeedHash // Update the global seedHash
	}

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

	const maxInputSize = 1024 * 1024 // 1 MB, adjust as needed
	if len(blobWithNonce) > maxInputSize {
		fmt.Printf("Input blobWithNonce size too large: %d bytes\n", len(blobWithNonce))
		return C.bool(false)
	}

	calculatedHash := calculateRandomXHash(globalRandomXVM, blobWithNonce, goSeedHash)
	fmt.Printf("Calculated RandomX Hash (hex): %s\n", hex.EncodeToString(calculatedHash))

	if !bytes.Equal(calculatedHash, goExpectedHash) {
		fmt.Printf("calculated hash does not match expected hash\n")
		return C.bool(false)
	}

	valid, err := CheckSolutionWithTarget(globalRandomXVM, blobWithNonce, calculatedHash, goTarget)
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
func calculateRandomXHash(vm unsafe.Pointer, blobWithNonce, seedHash []byte) []byte {

	if vm == nil {
		return nil
	}

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

func initRandomXSystem(flags RandomXFlags, seedHash []byte) error {

	if globalRandomXVM != nil {
		fmt.Printf("*999-*-*999*-**-*999*- ----- globalRandomXVM empty calling -------- DestroyVM ---------- *999-*-*999*-**-*999* *-*-*-*-*-*-**-*")
		DestroyVM(globalRandomXVM)
		globalRandomXVM = nil
	}
	if globalRandomXCache != nil {
		fmt.Printf("*999-*-*999*-**-*999*- ----- globalRandomXCache empty calling -------- DestroyRandomX ---------- *999-*-*999*-**-*999* *-*-*-*-*-*-**-*")
		DestroyRandomX(globalRandomXCache)
		globalRandomXCache = nil
	}

	// Reinitialize cache and VM
	globalRandomXCache = InitRandomX(flags)

	if globalRandomXCache == nil {
		fmt.Printf("*999-*-*999*-**-*999*- ----- InitRandomX error failed to allocate RandomX cache --- *999-*-*999*-**-*999* *-*-*-*-*-*-**-*")
		return fmt.Errorf("failed to allocate RandomX cache")
	}
	InitCache(globalRandomXCache, seedHash)

	globalRandomXVM = CreateVM(globalRandomXCache, flags)
	if globalRandomXVM == nil {
		fmt.Printf("*999-*-*999*-**-*999*- ----- InitRandomX error failed to create RandomX VM --- *999-*-*999*-**-*999* *-*-*-*-*-*-**-*")
		return fmt.Errorf("failed to create RandomX VM")
	}

	fmt.Printf("*-*-** 999999999 -**-*-*- ------------- initRandomXSystem() SUCCESS ---------- *-*-*-*999999999999 *-*-**-*")

	return nil
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
