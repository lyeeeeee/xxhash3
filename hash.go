// Package main implements https://github.com/Cyan4973/xxHash/blob/dev/xxhash.h
package xxhash3

import (
	"golang.org/x/sys/cpu"
	"math/bits"
	"unsafe"
)

func accumAVX2(acc *[8]uint64, xinput, xsecret unsafe.Pointer, len uintptr)
func accumSSE2(acc *[8]uint64, xinput, xsecret unsafe.Pointer, len uintptr)

var (
	avx2 = cpu.X86.HasAVX2
	sse2 = cpu.X86.HasSSE2
	xacc = [8]uint64{}
)

// Hash returns the hash value of the byte slice in 64bits.
func Hash(data []byte) uint64 {
	fn := xxh3HashLarge

	if len(data) <= 16 {
		fn = xxh3HashSmall
	}
	return fn(*(*unsafe.Pointer)(unsafe.Pointer(&data)), len(data))

}

// HashString returns the hash value of the string in 64bits.
func HashString(s string) uint64 {
	return Hash([]byte(s))
}

func xxh3HashSmall(xinput unsafe.Pointer, length int) uint64 {

	if length > 8 {
		inputlo := ReadUnaligned64(xinput) ^ xsecret_024 ^ xsecret_032
		inputhi := ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-8))) ^ xsecret_040 ^ xsecret_048
		return xxh3Avalanche(uint64(length) + bits.ReverseBytes64(inputlo) + inputhi + mix(inputlo, inputhi))
	} else if length >= 4 {
		input1 := ReadUnaligned32(xinput)
		input2 := ReadUnaligned32(unsafe.Pointer(uintptr(xinput) + uintptr(length-4)))
		input64 := input2 + input1<<32
		keyed := input64 ^ xsecret_008 ^ xsecret_016
		return xxh3RRMXMX(keyed, uint64(length))
	} else if length == 3 {
		c12 := ReadUnaligned16(xinput)
		c3 := uint64(*(*uint8)(unsafe.Pointer(uintptr(xinput) + 2)))
		acc := c12<<16 + c3 + 3<<8
		acc ^= uint64(xsecret32_000 ^ xsecret32_004)
		return xxh64Avalanche(acc)
	} else if length == 2 {
		c12 := ReadUnaligned16(xinput)
		acc := c12*(1<<24+1)>>8 + 2<<8
		acc ^= uint64(xsecret32_000 ^ xsecret32_004)
		return xxh64Avalanche(acc)
	} else if length == 1 {
		c1 := uint64(*(*uint8)(xinput))
		acc := c1*(1<<24+1<<16+1) + 1<<8
		acc ^= uint64(xsecret32_000 ^ xsecret32_004)
		return xxh64Avalanche(acc)
	}
	return 0x2d06800538d394c2
}

func xxh3HashLarge(xinput unsafe.Pointer, l int) (acc uint64) {
	length := uintptr(l)

	if length <= 128 {
		acc := uint64(length * prime64_1)
		if length > 32 {
			if length > 64 {
				if length > 96 {
					acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+48))^xsecret_096, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+56))^xsecret_104)
					acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-64)))^xsecret_112, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-56)))^xsecret_120)
				}
				acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+32))^xsecret_064, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+40))^xsecret_072)
				acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-48)))^xsecret_080, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-40)))^xsecret_088)
			}
			acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16))^xsecret_032, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+24))^xsecret_040)
			acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-32)))^xsecret_048, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-24)))^xsecret_056)
		}
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+0))^xsecret_000, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+8))^xsecret_008)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-16)))^xsecret_016, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-8)))^xsecret_024)

		return xxh3Avalanche(acc)

	} else if length <= 240 {
		acc := uint64(length * prime64_1)

		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*0))^xsecret_000, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*0+8))^xsecret_008)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*1))^xsecret_016, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*1+8))^xsecret_024)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*2))^xsecret_032, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*2+8))^xsecret_040)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*3))^xsecret_048, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*3+8))^xsecret_056)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*4))^xsecret_064, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*4+8))^xsecret_072)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*5))^xsecret_080, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*5+8))^xsecret_088)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*6))^xsecret_096, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*6+8))^xsecret_104)
		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*7))^xsecret_112, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+16*7+8))^xsecret_120)

		acc = xxh3Avalanche(acc)
		nbRounds := uint64(length >> 4 << 4)

		for i := uint64(8 * 16); i < nbRounds; i += 16 {
			acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(i)))^ReadUnaligned64(unsafe.Pointer(uintptr(xsecret)+uintptr(i-125))), ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(i+8)))^ReadUnaligned64(unsafe.Pointer(uintptr(xsecret)+uintptr(i-117))))
		}

		acc += mix(ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-16)))^xsecret_119, ReadUnaligned64(unsafe.Pointer(uintptr(xinput)+uintptr(length-8)))^xsecret_127)

		return xxh3Avalanche(acc)
	}

	xacc = [8]uint64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	acc = uint64(length * prime64_1)

	if avx2 {
		accumAVX2(&xacc, xinput, xsecret, length)
	} else if sse2 {
		accumSSE2(&xacc, xinput, xsecret, length)
	} else {
		accumScalar(&xacc, xinput, xsecret, length)
	}
	//merge xacc
	acc += mix(xacc[0]^xsecret_011, xacc[1]^xsecret_019)
	acc += mix(xacc[2]^xsecret_027, xacc[3]^xsecret_035)
	acc += mix(xacc[4]^xsecret_043, xacc[5]^xsecret_051)
	acc += mix(xacc[6]^xsecret_059, xacc[7]^xsecret_067)

	return xxh3Avalanche(acc)
}

func mix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}
func xxh3RRMXMX(h64 uint64, length uint64) uint64 {
	h64 ^= bits.RotateLeft64(h64, 49) ^ bits.RotateLeft64(h64, 24)
	h64 *= 0x9fb21c651e98df25
	h64 ^= (h64 >> 35) + length
	h64 *= 0x9fb21c651e98df25
	h64 ^= (h64 >> 28)
	return h64
}

func xxh64Avalanche(h64 uint64) uint64 {
	h64 *= prime64_2
	h64 ^= h64 >> 29
	h64 *= prime64_3
	h64 ^= h64 >> 32
	return h64
}

func xxh3Avalanche(x uint64) uint64 {
	x ^= x >> 37
	x *= 0x165667919e3779f9
	x ^= x >> 32
	return x
}
