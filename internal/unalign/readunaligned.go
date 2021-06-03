// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm
//
// from golang-go/src/os/endian_big.go

package unalign

import (
	"unsafe"
)

func Read8(p unsafe.Pointer, offset uintptr) uint64 {
	p = unsafe.Pointer(uintptr(p) + offset)
	q := (*[8]byte)(p)
	return uint64(q[0]) | uint64(q[1])<<8 | uint64(q[2])<<16 | uint64(q[3])<<24 | uint64(q[4])<<32 | uint64(q[5])<<40 | uint64(q[6])<<48 | uint64(q[7])<<56
}

func Read4(p unsafe.Pointer, offset uintptr) uint64 {
	p = unsafe.Pointer(uintptr(p) + offset)
	q := (*[4]byte)(p)
	return uint64(q[0]) | uint64(q[1])<<8 | uint64(q[2])<<16 | uint64(q[3])<<24
}
