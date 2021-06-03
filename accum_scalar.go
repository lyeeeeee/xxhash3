package xxhash3

import (
	"unsafe"
)

func ReadUnaligned64(p unsafe.Pointer) uint64 {
	// Equal to runtime.readUnaligned64, but this function can be inlined
	// compared to  use runtime.readUnaligned64 via go:linkname.
	q := (*[8]byte)(p)
	return uint64(q[0]) | uint64(q[1])<<8 | uint64(q[2])<<16 | uint64(q[3])<<24 | uint64(q[4])<<32 | uint64(q[5])<<40 | uint64(q[6])<<48 | uint64(q[7])<<56
}

func ReadUnaligned32(p unsafe.Pointer) uint64 {
	q := (*[4]byte)(p)
	return uint64(uint32(q[0]) | uint32(q[1])<<8 | uint32(q[2])<<16 | uint32(q[3])<<24)
}

func ReadUnaligned16(p unsafe.Pointer) uint64 {
	q := (*[2]byte)(p)
	return uint64(uint32(q[0]) | uint32(q[1])<<8)
}

func accumScalar(xacc *[8]uint64, xinput, xsecret unsafe.Pointer, l uintptr) {
	j := uintptr(0)

	// Loops over blocunsafe.Pointer(uintptr(k)+ process 16*8*8=1024 bytes of data each iteration
	for ; j < (l-1)/1024; j++ {
		k := xsecret
		for i := 0; i < 16; i++ {
			dataVec0 := ReadUnaligned64(xinput)

			keyVec := dataVec0 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*0))
			xacc[1] += dataVec0
			xacc[0] += (keyVec & 0xffffffff) * (keyVec >> 32)

			dataVec1 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*1))
			keyVec1 := dataVec1 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*1))
			xacc[0] += dataVec1
			xacc[1] += (keyVec1 & 0xffffffff) * (keyVec1 >> 32)

			dataVec2 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*2))
			keyVec2 := dataVec2 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*2))
			xacc[3] += dataVec2
			xacc[2] += (keyVec2 & 0xffffffff) * (keyVec2 >> 32)

			dataVec3 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*3))
			keyVec3 := dataVec3 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*3))
			xacc[2] += dataVec3
			xacc[3] += (keyVec3 & 0xffffffff) * (keyVec3 >> 32)

			dataVec4 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*4))
			keyVec4 := dataVec4 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*4))
			xacc[5] += dataVec4
			xacc[4] += (keyVec4 & 0xffffffff) * (keyVec4 >> 32)

			dataVec5 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*5))
			keyVec5 := dataVec5 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*5))
			xacc[4] += dataVec5
			xacc[5] += (keyVec5 & 0xffffffff) * (keyVec5 >> 32)

			dataVec6 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*6))
			keyVec6 := dataVec6 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*6))
			xacc[7] += dataVec6
			xacc[6] += (keyVec6 & 0xffffffff) * (keyVec6 >> 32)

			dataVec7 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*7))
			keyVec7 := dataVec7 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*7))
			xacc[6] += dataVec7
			xacc[7] += (keyVec7 & 0xffffffff) * (keyVec7 >> 32)

			xinput, k = unsafe.Pointer(uintptr(xinput)+_stripe), unsafe.Pointer(uintptr(k)+8)
		}

		// scramble xacc
		xacc[0] ^= xacc[0] >> 47
		xacc[0] ^= xsecret_128
		xacc[0] *= prime32_1

		xacc[1] ^= xacc[1] >> 47
		xacc[1] ^= xsecret_136
		xacc[1] *= prime32_1

		xacc[2] ^= xacc[2] >> 47
		xacc[2] ^= xsecret_144
		xacc[2] *= prime32_1

		xacc[3] ^= xacc[3] >> 47
		xacc[3] ^= xsecret_152
		xacc[3] *= prime32_1

		xacc[4] ^= xacc[4] >> 47
		xacc[4] ^= xsecret_160
		xacc[4] *= prime32_1

		xacc[5] ^= xacc[5] >> 47
		xacc[5] ^= xsecret_168
		xacc[5] *= prime32_1

		xacc[6] ^= xacc[6] >> 47
		xacc[6] ^= xsecret_176
		xacc[6] *= prime32_1

		xacc[7] ^= xacc[7] >> 47
		xacc[7] ^= xsecret_184
		xacc[7] *= prime32_1

	}
	l -= _block * j

	// last partial block (1024 bytes)
	if l > 0 {
		k := xsecret
		i := uintptr(0)
		for ; i < (l-1)/_stripe; i++ {
			dataVec := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*0))
			keyVec := dataVec ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*0))
			xacc[1] += dataVec
			xacc[0] += (keyVec & 0xffffffff) * (keyVec >> 32)

			dataVec1 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*1))
			keyVec1 := dataVec1 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*1))
			xacc[0] += dataVec1
			xacc[1] += (keyVec1 & 0xffffffff) * (keyVec1 >> 32)

			dataVec2 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*2))
			keyVec2 := dataVec2 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*2))
			xacc[3] += dataVec2
			xacc[2] += (keyVec2 & 0xffffffff) * (keyVec2 >> 32)

			dataVec3 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*3))
			keyVec3 := dataVec3 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*3))
			xacc[2] += dataVec3
			xacc[3] += (keyVec3 & 0xffffffff) * (keyVec3 >> 32)

			dataVec4 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*4))
			keyVec4 := dataVec4 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*4))
			xacc[5] += dataVec4
			xacc[4] += (keyVec4 & 0xffffffff) * (keyVec4 >> 32)

			dataVec5 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*5))
			keyVec5 := dataVec5 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*5))
			xacc[4] += dataVec5
			xacc[5] += (keyVec5 & 0xffffffff) * (keyVec5 >> 32)

			dataVec6 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*6))
			keyVec6 := dataVec6 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*6))
			xacc[7] += dataVec6
			xacc[6] += (keyVec6 & 0xffffffff) * (keyVec6 >> 32)

			dataVec7 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*7))
			keyVec7 := dataVec7 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*7))
			xacc[6] += dataVec7
			xacc[7] += (keyVec7 & 0xffffffff) * (keyVec7 >> 32)

			xinput, k = unsafe.Pointer(uintptr(xinput)+_stripe), unsafe.Pointer(uintptr(k)+8)
		}
		l -= _stripe * i

		// last stripe (64 bytes)
		if l > 0 {
			xinput = unsafe.Pointer(uintptr(xinput) - uintptr(_stripe-l))
			k = unsafe.Pointer(uintptr(xsecret) + 121)

			dataVec := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*0))
			keyVec := dataVec ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*0))
			xacc[1] += dataVec
			xacc[0] += (keyVec & 0xffffffff) * (keyVec >> 32)

			dataVec1 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*1))
			keyVec1 := dataVec1 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*1))
			xacc[0] += dataVec1
			xacc[1] += (keyVec1 & 0xffffffff) * (keyVec1 >> 32)

			dataVec2 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*2))
			keyVec2 := dataVec2 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*2))
			xacc[3] += dataVec2
			xacc[2] += (keyVec2 & 0xffffffff) * (keyVec2 >> 32)

			dataVec3 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*3))
			keyVec3 := dataVec3 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*3))
			xacc[2] += dataVec3
			xacc[3] += (keyVec3 & 0xffffffff) * (keyVec3 >> 32)

			dataVec4 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*4))
			keyVec4 := dataVec4 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*4))
			xacc[5] += dataVec4
			xacc[4] += (keyVec4 & 0xffffffff) * (keyVec4 >> 32)

			dataVec5 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*5))
			keyVec5 := dataVec5 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*5))
			xacc[4] += dataVec5
			xacc[5] += (keyVec5 & 0xffffffff) * (keyVec5 >> 32)

			dataVec6 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*6))
			keyVec6 := dataVec6 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*6))
			xacc[7] += dataVec6
			xacc[6] += (keyVec6 & 0xffffffff) * (keyVec6 >> 32)

			dataVec7 := ReadUnaligned64(unsafe.Pointer(uintptr(xinput) + 8*7))
			keyVec7 := dataVec7 ^ ReadUnaligned64(unsafe.Pointer(uintptr(k)+8*7))
			xacc[6] += dataVec7
			xacc[7] += (keyVec7 & 0xffffffff) * (keyVec7 >> 32)
		}
	}
}
