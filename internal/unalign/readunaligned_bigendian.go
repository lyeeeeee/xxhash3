// +build ppc64 s390x mips mips64
//
// from golang-go/src/os/endian_little.go

package unalign

import (
	"unsafe"
)

func Read8(p unsafe.Pointer, offset uintptr) uint64 {
	p = unsafe.Pointer(uintptr(p) + offset)
	q := (*[8]byte)(p)
	return uint64(q[7]) | uint64(q[6])<<8 | uint64(q[5])<<16 | uint64(q[4])<<24 | uint64(q[3])<<32 | uint64(q[2])<<40 | uint64(q[1])<<48 | uint64(q[0])<<56
}

func Read4(p unsafe.Pointer, offset uintptr) uint64 {
	p = unsafe.Pointer(uintptr(p) + offset)
	q := (*[4]byte)(p)
	return uint64(q[3]) | uint64(q[2])<<8 | uint64(q[1])<<16 | uint64(q[0])<<24
}

func Read2(p unsafe.Pointer, offset uintptr) uint64 {
	p = unsafe.Pointer(uintptr(p) + offset)
	q := (*[2]byte)(p)
	return uint64(q[1]) | uint64(q[0])<<8
}
