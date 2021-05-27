// Package main implements https://github.com/Cyan4973/xxHash/blob/dev/xxhash.h
package xxhash3

import (
	"./internal/unalign"
	"golang.org/x/sys/cpu"
	"math/bits"
	"unsafe"
)

func accumAVX2(acc *[8]uint64, xinput, xsecret unsafe.Pointer, len uint64)
func accumSSE2(acc *[8]uint64, xinput, xsecret unsafe.Pointer, len uint64)

var (
	avx2 = cpu.X86.HasAVX2
	sse2 = cpu.X86.HasSSE2
)

// Hash returns the hash value of the byte slice in 64bits.
func Hash(data []byte) uint64 {
	length := uint64(len(data))
	xinput := *(*unsafe.Pointer)(unsafe.Pointer(&data))

	if length > 240 {
		return hashLarge(xinput, length)
	} else if length > 128 {
		return xxh3Len129To240_64b(xinput, length)
	} else if length > 16 {
		return xxh3Len17To128_64b(xinput, length)
	} else {
		return xxh3Len0To16_64b(xinput, length)
	}
}

// HashString returns the hash value of the string in 64bits.
func HashString(s string) uint64 {
	return Hash([]byte(s))
}

// Hash128 returns the hash value of the byte slice in 128bits.
func Hash128(data []byte) [2]uint64 {
	length := uint64(len(data))
	xinput := *(*unsafe.Pointer)(unsafe.Pointer(&data))

	if length > 240 {
		return hashLarge128(xinput, length)
	} else if length > 128 {
		return xxh3Len129To240_128b(xinput, length)
	} else if length > 16 {
		return xxh3Len17To128_128b(xinput, length)
	} else {
		return xxh3Len0To16_128b(xinput, length)
	}
}

// Hash128String returns the hash value of the string in 128bits.
func Hash128String(s string) [2]uint64 {
	return Hash128([]byte(s))
}

func xxh3Len0To16_64b(xinput unsafe.Pointer, length uint64) uint64 {
	if length > 8 {
		inputlo := unalign.Read8(xinput, 0) ^ (unalign.Read8(xsecret, 24)) ^ unalign.Read8(xsecret, 32)
		inputhi := unalign.Read8(xinput, uintptr(length-8)) ^ (unalign.Read8(xsecret, 40)) ^ unalign.Read8(xsecret, 48)
		acc := length + bits.ReverseBytes64(inputlo) + inputhi + mix(inputlo, inputhi)
		return xxh3Avalanche(acc)
	} else if length >= 4 {
		input1 := unalign.Read4(xinput, 0)
		input2 := unalign.Read4(xinput, uintptr(length-4))
		input64 := input2 + input1<<32
		keyed := input64 ^ (unalign.Read8(xsecret, 8)) ^ unalign.Read8(xsecret, 16)
		return xxh3RRMXMX(keyed, length)
	} else if length > 0 {
		q := (*[4]byte)(xinput)
		combined := (uint64(q[0]) << 16) | (uint64(q[length>>1]) << 24) | (uint64(q[length-1]) << 0) | length<<8
		combined ^= unalign.Read4(xsecret, 0) ^ unalign.Read4(xsecret, 4)
		return xxh64Avalanche(combined)
	} else {
		return xxh64Avalanche(unalign.Read8(xsecret, 56) ^ unalign.Read8(xsecret, 64))
	}
}

func xxh3Len17To128_64b(xinput unsafe.Pointer, length uint64) uint64 {
	acc := length * prime64_1
	if length > 32 {
		if length > 64 {
			if length > 96 {
				acc += mix(unalign.Read8(xinput, 48)^unalign.Read8(xsecret, 96), unalign.Read8(xinput, 56)^unalign.Read8(xsecret, 104))
				acc += mix(unalign.Read8(xinput, uintptr(length-64))^unalign.Read8(xsecret, 112), unalign.Read8(xinput, uintptr(length-56))^unalign.Read8(xsecret, 120))
			}
			acc += mix(unalign.Read8(xinput, 32)^unalign.Read8(xsecret, 64), unalign.Read8(xinput, 40)^unalign.Read8(xsecret, 72))
			acc += mix(unalign.Read8(xinput, uintptr(length-48))^unalign.Read8(xsecret, 80), unalign.Read8(xinput, uintptr(length-40))^unalign.Read8(xsecret, 88))
		}
		acc += mix(unalign.Read8(xinput, 16)^unalign.Read8(xsecret, 32), unalign.Read8(xinput, 24)^unalign.Read8(xsecret, 40))
		acc += mix(unalign.Read8(xinput, uintptr(length-32))^unalign.Read8(xsecret, 48), unalign.Read8(xinput, uintptr(length-24))^unalign.Read8(xsecret, 56))
	}
	acc += mix(unalign.Read8(xinput, 0)^unalign.Read8(xsecret, 0), unalign.Read8(xinput, 8)^unalign.Read8(xsecret, 8))
	acc += mix(unalign.Read8(xinput, uintptr(length-16))^unalign.Read8(xsecret, 16), unalign.Read8(xinput, uintptr(length-8))^unalign.Read8(xsecret, 24))

	return xxh3Avalanche(acc)
}

func xxh3Len129To240_64b(xinput unsafe.Pointer, length uint64) uint64 {

	acc := length * prime64_1

	acc += mix(unalign.Read8(xinput, 16*0)^unalign.Read8(xsecret, 16*0), unalign.Read8(xinput, 16*0+8)^unalign.Read8(xsecret, 16*0+8))
	acc += mix(unalign.Read8(xinput, 16*1)^unalign.Read8(xsecret, 16*1), unalign.Read8(xinput, 16*1+8)^unalign.Read8(xsecret, 16*1+8))
	acc += mix(unalign.Read8(xinput, 16*2)^unalign.Read8(xsecret, 16*2), unalign.Read8(xinput, 16*2+8)^unalign.Read8(xsecret, 16*2+8))
	acc += mix(unalign.Read8(xinput, 16*3)^unalign.Read8(xsecret, 16*3), unalign.Read8(xinput, 16*3+8)^unalign.Read8(xsecret, 16*3+8))
	acc += mix(unalign.Read8(xinput, 16*4)^unalign.Read8(xsecret, 16*4), unalign.Read8(xinput, 16*4+8)^unalign.Read8(xsecret, 16*4+8))
	acc += mix(unalign.Read8(xinput, 16*5)^unalign.Read8(xsecret, 16*5), unalign.Read8(xinput, 16*5+8)^unalign.Read8(xsecret, 16*5+8))
	acc += mix(unalign.Read8(xinput, 16*6)^unalign.Read8(xsecret, 16*6), unalign.Read8(xinput, 16*6+8)^unalign.Read8(xsecret, 16*6+8))
	acc += mix(unalign.Read8(xinput, 16*7)^unalign.Read8(xsecret, 16*7), unalign.Read8(xinput, 16*7+8)^unalign.Read8(xsecret, 16*7+8))

	acc = xxh3Avalanche(acc)
	nbRounds := length >> 4

	for i := uint64(8); i < nbRounds; i++ {
		acc += mix(unalign.Read8(xinput, uintptr(16*i))^unalign.Read8(xsecret, uintptr(16*i-125)), unalign.Read8(xinput, uintptr(16*i+8))^unalign.Read8(xsecret, uintptr(16*i-117)))
	}

	acc += mix(unalign.Read8(xinput, uintptr(length-16))^unalign.Read8(xsecret, 119), unalign.Read8(xinput, uintptr(length-8))^unalign.Read8(xsecret, uintptr(127)))

	return xxh3Avalanche(acc)
}

func hashLarge(p unsafe.Pointer, length uint64) (acc uint64) {
	acc = length * prime64_1

	xacc := [8]uint64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	if avx2 {
		accumAVX2(&xacc, p, xsecret, length)
	} else if sse2 {
		accumSSE2(&xacc, p, xsecret, length)
	} else {
		accumScalar(&xacc, p, xsecret, length)
	}
	//merge xacc
	acc += mix(xacc[0]^key64_011, xacc[1]^key64_019)
	acc += mix(xacc[2]^key64_027, xacc[3]^key64_035)
	acc += mix(xacc[4]^key64_043, xacc[5]^key64_051)
	acc += mix(xacc[6]^key64_059, xacc[7]^key64_067)
	return xxh3Avalanche(acc)
}

func xxh3Len0To16_128b(xinput unsafe.Pointer, length uint64) [2]uint64 {

	if length > 8 {
		bitflipl := unalign.Read8(xsecret, 32) ^ unalign.Read8(xsecret, 40)
		bitfliph := unalign.Read8(xsecret, 48) ^ unalign.Read8(xsecret, 56)
		inputLow := unalign.Read8(xinput, 0)
		inputHigh := unalign.Read8(xinput, uintptr(length)-8)
		m128High64, m128Low64 := bits.Mul64(inputLow^inputHigh^bitflipl, prime64_1)

		m128Low64 += uint64(length-1) << 54
		inputHigh ^= bitfliph

		m128High64 += inputHigh + uint64(uint32(inputHigh))*(prime32_2-1)
		m128Low64 ^= bits.ReverseBytes64(m128High64)

		h128High64, h128Low64 := bits.Mul64(m128Low64, prime64_2)
		h128High64 += m128High64 * prime64_2

		h128Low64 = xxh3Avalanche(h128Low64)
		h128High64 = xxh3Avalanche(h128High64)

		return [2]uint64{h128High64, h128Low64}

	} else if length >= 4 {
		inputLow := unalign.Read4(xinput, 0)
		inputHigh := unalign.Read4(xinput, uintptr(length)-4)
		input64 := inputLow + (uint64(inputHigh) << 32)
		bitflip := unalign.Read8(xsecret, 16) ^ unalign.Read8(xsecret, 24)
		keyed := input64 ^ bitflip

		m128High64, m128Low64 := bits.Mul64(keyed, prime64_1+(length)<<2)
		m128High64 += m128Low64 << 1
		m128Low64 ^= m128High64 >> 3

		m128Low64 ^= m128Low64 >> 35
		m128Low64 *= 0x9fb21c651e98df25
		m128Low64 ^= m128Low64 >> 28

		m128High64 = xxh3Avalanche(m128High64)
		return [2]uint64{m128High64, m128Low64}

	} else if length >= 1 {
		q := (*[4]byte)(xinput)
		combinedl := (uint64(q[0]) << 16) | (uint64(q[length>>1]) << 24) | (uint64(q[length-1]) << 0) | length<<8
		combinedh := uint64(bits.RotateLeft32(bits.ReverseBytes32(uint32(combinedl)), 13))

		bitflipl := unalign.Read4(xsecret, 0) ^ unalign.Read4(xsecret, 4)
		bitfliph := unalign.Read4(xsecret, 8) ^ unalign.Read4(xsecret, 12)

		keyedLow := combinedl ^ bitflipl
		keyedHigh := combinedh ^ bitfliph

		keyedLow = combinedl ^ bitflipl
		keyedHigh = combinedh ^ bitfliph

		h128Low64 := xxh64Avalanche(keyedLow)
		h128High64 := xxh64Avalanche(keyedHigh)
		return [2]uint64{h128High64, h128Low64}
	}
	bitflipl := unalign.Read8(xsecret, 64) ^ unalign.Read8(xsecret, 72)
	bitfliph := unalign.Read8(xsecret, 80) ^ unalign.Read8(xsecret, 88)

	h128High64 := xxh64Avalanche(bitfliph)
	h128Low64 := xxh64Avalanche(bitflipl)

	return [2]uint64{h128High64, h128Low64}
}

func xxh3Len17To128_128b(xinput unsafe.Pointer, length uint64) [2]uint64 {

	accHigh := uint64(0)
	accLow := length * prime64_1

	if length > 32 {
		if length > 64 {
			if length > 96 {
				accLow += mix(unalign.Read8(xinput, 48)^unalign.Read8(xsecret, 96), unalign.Read8(xinput, 56)^unalign.Read8(xsecret, 104))
				accLow ^= unalign.Read8(xinput, uintptr(length-64)) + unalign.Read8(xinput, uintptr(length-56))
				accHigh += mix(unalign.Read8(xinput, uintptr(length-64))^unalign.Read8(xsecret, 112), unalign.Read8(xinput, uintptr(length-56))^unalign.Read8(xsecret, 120))
				accHigh ^= unalign.Read8(xinput, 48) + unalign.Read8(xinput, 56)
			}
			accLow += mix(unalign.Read8(xinput, 32)^unalign.Read8(xsecret, 64), unalign.Read8(xinput, 40)^unalign.Read8(xsecret, 72))
			accLow ^= unalign.Read8(xinput, uintptr(length-48)) + unalign.Read8(xinput, uintptr(length-40))
			accHigh += mix(unalign.Read8(xinput, uintptr(length-48))^unalign.Read8(xsecret, 80), unalign.Read8(xinput, uintptr(length-40))^unalign.Read8(xsecret, 88))
			accHigh ^= unalign.Read8(xinput, 32) + unalign.Read8(xinput, 40)
		}
		accLow += mix(unalign.Read8(xinput, 16)^unalign.Read8(xsecret, 32), unalign.Read8(xinput, 3*8)^unalign.Read8(xsecret, 40))
		accLow ^= unalign.Read8(xinput, uintptr(length-32)) + unalign.Read8(xinput, uintptr(length-3*8))
		accHigh += mix(unalign.Read8(xinput, uintptr(length-32))^unalign.Read8(xsecret, 48), unalign.Read8(xinput, uintptr(length-3*8))^unalign.Read8(xsecret, 56))
		accHigh ^= unalign.Read8(xinput, 16) + unalign.Read8(xinput, 3*8)
	}

	accLow += mix(unalign.Read8(xinput, 0)^unalign.Read8(xsecret, 0), unalign.Read8(xinput, 8)^unalign.Read8(xsecret, 8))
	accLow ^= unalign.Read8(xinput, uintptr(length-16)) + unalign.Read8(xinput, uintptr(length-8))
	accHigh += mix(unalign.Read8(xinput, uintptr(length-16))^unalign.Read8(xsecret, 16), unalign.Read8(xinput, uintptr(length-8))^unalign.Read8(xsecret, 24))
	accHigh ^= unalign.Read8(xinput, 0) + unalign.Read8(xinput, 8)

	h128Low := accHigh + accLow
	h128High := (accLow * prime64_1) + (accHigh * prime64_4) + (length * prime64_2)

	h128Low = xxh3Avalanche(h128Low)
	h128High = -xxh3Avalanche(h128High)

	return [2]uint64{h128High, h128Low}
}

func xxh3Len129To240_128b(xinput unsafe.Pointer, length uint64) [2]uint64 {
	nbRounds := length &^ 31 / 32
	accLow64 := length * prime64_1
	accHigh64 := uint64(0)

	for i := 0; i < 4; i++ {
		accLow64 += mix(unalign.Read8(xinput, uintptr(32*i))^unalign.Read8(xsecret, uintptr(32*i)), unalign.Read8(xinput, uintptr(32*i+8))^unalign.Read8(xsecret, uintptr(32*i+8)))
		accLow64 ^= unalign.Read8(xinput, uintptr(32*i+16)) + unalign.Read8(xinput, uintptr(32*i+24))
		accHigh64 += mix(unalign.Read8(xinput, uintptr(32*i+16))^unalign.Read8(xsecret, uintptr(32*i+16)), unalign.Read8(xinput, uintptr(32*i)+24)^unalign.Read8(xsecret, uintptr(32*i+24)))
		accHigh64 ^= unalign.Read8(xinput, uintptr(32*i)) + unalign.Read8(xinput, uintptr(32*i)+8)
	}

	accLow64 = xxh3Avalanche(accLow64)
	accHigh64 = xxh3Avalanche(accHigh64)

	for i := uint64(4); i < nbRounds; i++ {
		accHigh64 += mix(unalign.Read8(xinput, uintptr(32*i+16))^unalign.Read8(xsecret, uintptr(32*i-109)), unalign.Read8(xinput, uintptr(32*i)+24)^unalign.Read8(xsecret, uintptr(32*i-101)))
		accHigh64 ^= unalign.Read8(xinput, uintptr(32*i)) + unalign.Read8(xinput, uintptr(32*i)+8)

		accLow64 += mix(unalign.Read8(xinput, uintptr(32*i))^unalign.Read8(xsecret, uintptr(32*i-125)), unalign.Read8(xinput, uintptr(32*i+8))^unalign.Read8(xsecret, uintptr(32*i-117)))
		accLow64 ^= unalign.Read8(xinput, uintptr(32*i+16)) + unalign.Read8(xinput, uintptr(32*i+24))
	}

	// last 32 bytes
	accLow64 += mix(unalign.Read8(xinput, uintptr(length-16))^unalign.Read8(xsecret, 103), unalign.Read8(xinput, uintptr(length-8))^unalign.Read8(xsecret, 111))
	accLow64 ^= unalign.Read8(xinput, uintptr(length-32)) + unalign.Read8(xinput, uintptr(length-24))
	accHigh64 += mix(unalign.Read8(xinput, uintptr(length-32))^unalign.Read8(xsecret, 119), unalign.Read8(xinput, uintptr(length-24))^unalign.Read8(xsecret, 127))
	accHigh64 ^= unalign.Read8(xinput, uintptr(length-16)) + unalign.Read8(xinput, uintptr(length-8))

	accHigh64, accLow64 = (accLow64*prime64_1)+(accHigh64*prime64_4)+(length*prime64_2), accHigh64+accLow64

	accLow64 = xxh3Avalanche(accLow64)
	accHigh64 = -xxh3Avalanche(accHigh64)

	return [2]uint64{accHigh64, accLow64}
}

func hashLarge128(p unsafe.Pointer, length uint64) (acc [2]uint64) {
	acc[1] = length * prime64_1
	acc[0] = ^(length * prime64_2)

	xacc := [8]uint64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	accumScalar(&xacc, p, xsecret, length)
	// merge xacc
	acc[1] += mix(xacc[0]^unalign.Read8(xsecret, 11), xacc[1]^unalign.Read8(xsecret, 19))
	acc[1] += mix(xacc[2]^unalign.Read8(xsecret, 27), xacc[3]^unalign.Read8(xsecret, 35))
	acc[1] += mix(xacc[4]^unalign.Read8(xsecret, 43), xacc[5]^unalign.Read8(xsecret, 51))
	acc[1] += mix(xacc[6]^unalign.Read8(xsecret, 59), xacc[7]^unalign.Read8(xsecret, 67))

	acc[1] = xxh3Avalanche(acc[1])

	acc[0] += mix(xacc[0]^unalign.Read8(xsecret, 117), xacc[1]^unalign.Read8(xsecret, 125))
	acc[0] += mix(xacc[2]^unalign.Read8(xsecret, 133), xacc[3]^unalign.Read8(xsecret, 141))
	acc[0] += mix(xacc[4]^unalign.Read8(xsecret, 149), xacc[5]^unalign.Read8(xsecret, 157))
	acc[0] += mix(xacc[6]^unalign.Read8(xsecret, 165), xacc[7]^unalign.Read8(xsecret, 173))
	acc[0] = xxh3Avalanche(acc[0])

	return acc
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
	h64 ^= h64 >> 33
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
