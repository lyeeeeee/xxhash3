package xxhash3

import (
	"./internal/unalign"
	"math/bits"
	"unsafe"
)

// Hash128 returns the hash value of the byte slice in 128bits.
func Hash128(data []byte) [2]uint64 {
	fn := xxh3HashLarge128
	if len(data) <= 16 {
		fn = xxh3HashSmall128
	}
	return fn(*(*unsafe.Pointer)(unsafe.Pointer(&data)), len(data))
}

// Hash128String returns the hash value of the string in 128bits.
func Hash128String(s string) [2]uint64 {
	return Hash128([]byte(s))
}

func xxh3HashSmall128(xinput unsafe.Pointer, l int) [2]uint64 {
	length := uint64(l)
	var h128Low64 uint64
	var h128High64 uint64

	if length > 8 {
		bitflipl := xsecret_032 ^ xsecret_040
		bitfliph := xsecret_048 ^ xsecret_056
		inputLow := unalign.Read8(xinput, 0)
		inputHigh := unalign.Read8(xinput, uintptr(length)-8)
		m128High64, m128Low64 := bits.Mul64(inputLow^inputHigh^bitflipl, prime64_1)

		m128Low64 += uint64(length-1) << 54
		inputHigh ^= bitfliph

		m128High64 += inputHigh + uint64(uint32(inputHigh))*(prime32_2-1)
		m128Low64 ^= bits.ReverseBytes64(m128High64)

		h128High64, h128Low64 = bits.Mul64(m128Low64, prime64_2)
		h128High64 += m128High64 * prime64_2

		h128Low64 = xxh3Avalanche(h128Low64)
		h128High64 = xxh3Avalanche(h128High64)

		return [2]uint64{h128High64, h128Low64}
	} else if length >= 4 {
		inputLow := unalign.Read4(xinput, 0)
		inputHigh := unalign.Read4(xinput, uintptr(length)-4)
		input64 := inputLow + (uint64(inputHigh) << 32)
		bitflip := xsecret_016 ^ xsecret_024
		keyed := input64 ^ bitflip

		h128High64, h128Low64 = bits.Mul64(keyed, prime64_1+(length)<<2)
		h128High64 += h128Low64 << 1
		h128Low64 ^= h128High64 >> 3

		h128Low64 ^= h128Low64 >> 35
		h128Low64 *= 0x9fb21c651e98df25
		h128Low64 ^= h128Low64 >> 28

		h128High64 = xxh3Avalanche(h128High64)
		return [2]uint64{h128High64, h128Low64}
	} else if length == 3 {
		c12 := unalign.Read2(xinput, 0)
		c3 := uint64(*(*uint8)(unsafe.Pointer(uintptr(xinput) + 2)))
		h128Low64 = c12<<16 + c3 + 3<<8
	} else if length == 2 {
		c12 := unalign.Read2(xinput, 0)
		h128Low64 = c12*(1<<24+1)>>8 + 2<<8
	} else if length == 1 {
		c1 := uint64(*(*uint8)(xinput))
		h128Low64 = c1*(1<<24+1<<16+1) + 1<<8
	} else if length == 0 {
		return [2]uint64{0x99aa06d3014798d8, 0x6001c324468d497f}
	}
	h128High64 = uint64(bits.RotateLeft32(bits.ReverseBytes32(uint32(h128Low64)), 13))
	bitflipl := uint64(xsecret32_000 ^ xsecret32_004)
	bitfliph := uint64(xsecret32_008 ^ xsecret32_012)

	h128Low64 = h128Low64 ^ bitflipl
	h128High64 = h128High64 ^ bitfliph

	h128Low64 = xxh64Avalanche(h128Low64)
	h128High64 = xxh64Avalanche(h128High64)
	return [2]uint64{h128High64, h128Low64}
}

func xxh3HashLarge128(xinput unsafe.Pointer, l int) (acc [2]uint64) {
	length := uint64(l)

	if length <= 128 {

		accHigh := uint64(0)
		accLow := length * prime64_1

		if length > 32 {
			if length > 64 {
				if length > 96 {
					accLow += mix(unalign.Read8(xinput, 48)^xsecret_096, unalign.Read8(xinput, 56)^xsecret_104)
					accLow ^= unalign.Read8(xinput, uintptr(length-64)) + unalign.Read8(xinput, uintptr(length-56))
					accHigh += mix(unalign.Read8(xinput, uintptr(length-64))^xsecret_112, unalign.Read8(xinput, uintptr(length-56))^xsecret_120)
					accHigh ^= unalign.Read8(xinput, 48) + unalign.Read8(xinput, 56)
				}
				accLow += mix(unalign.Read8(xinput, 32)^xsecret_064, unalign.Read8(xinput, 40)^xsecret_072)
				accLow ^= unalign.Read8(xinput, uintptr(length-48)) + unalign.Read8(xinput, uintptr(length-40))
				accHigh += mix(unalign.Read8(xinput, uintptr(length-48))^xsecret_080, unalign.Read8(xinput, uintptr(length-40))^xsecret_088)
				accHigh ^= unalign.Read8(xinput, 32) + unalign.Read8(xinput, 40)
			}
			accLow += mix(unalign.Read8(xinput, 16)^xsecret_032, unalign.Read8(xinput, 3*8)^xsecret_040)
			accLow ^= unalign.Read8(xinput, uintptr(length-32)) + unalign.Read8(xinput, uintptr(length-3*8))
			accHigh += mix(unalign.Read8(xinput, uintptr(length-32))^xsecret_048, unalign.Read8(xinput, uintptr(length-3*8))^xsecret_056)
			accHigh ^= unalign.Read8(xinput, 16) + unalign.Read8(xinput, 3*8)
		}

		accLow += mix(unalign.Read8(xinput, 0)^xsecret_000, unalign.Read8(xinput, 8)^xsecret_008)
		accLow ^= unalign.Read8(xinput, uintptr(length-16)) + unalign.Read8(xinput, uintptr(length-8))
		accHigh += mix(unalign.Read8(xinput, uintptr(length-16))^xsecret_016, unalign.Read8(xinput, uintptr(length-8))^xsecret_024)
		accHigh ^= unalign.Read8(xinput, 0) + unalign.Read8(xinput, 8)

		h128Low := accHigh + accLow
		h128High := (accLow * prime64_1) + (accHigh * prime64_4) + (length * prime64_2)

		h128Low = xxh3Avalanche(h128Low)
		h128High = -xxh3Avalanche(h128High)

		return [2]uint64{h128High, h128Low}
	} else if length <= 240 {
		accLow64 := length * prime64_1
		accHigh64 := uint64(0)

		accLow64 += mix(unalign.Read8(xinput, uintptr(32*0))^xsecret_000, unalign.Read8(xinput, uintptr(32*0+8))^xsecret_008)
		accLow64 ^= unalign.Read8(xinput, uintptr(32*0+16)) + unalign.Read8(xinput, uintptr(32*0+24))
		accHigh64 += mix(unalign.Read8(xinput, uintptr(32*0+16))^xsecret_016, unalign.Read8(xinput, uintptr(32*0)+24)^xsecret_024)
		accHigh64 ^= unalign.Read8(xinput, uintptr(32*0)) + unalign.Read8(xinput, uintptr(32*0)+8)

		accLow64 += mix(unalign.Read8(xinput, uintptr(32*1))^xsecret_032, unalign.Read8(xinput, uintptr(32*1+8))^xsecret_040)
		accLow64 ^= unalign.Read8(xinput, uintptr(32*1+16)) + unalign.Read8(xinput, uintptr(32*1+24))
		accHigh64 += mix(unalign.Read8(xinput, uintptr(32*1+16))^xsecret_048, unalign.Read8(xinput, uintptr(32*1)+24)^xsecret_056)
		accHigh64 ^= unalign.Read8(xinput, uintptr(32*1)) + unalign.Read8(xinput, uintptr(32*1)+8)

		accLow64 += mix(unalign.Read8(xinput, uintptr(32*2))^xsecret_064, unalign.Read8(xinput, uintptr(32*2+8))^xsecret_072)
		accLow64 ^= unalign.Read8(xinput, uintptr(32*2+16)) + unalign.Read8(xinput, uintptr(32*2+24))
		accHigh64 += mix(unalign.Read8(xinput, uintptr(32*2+16))^xsecret_080, unalign.Read8(xinput, uintptr(32*2)+24)^xsecret_088)
		accHigh64 ^= unalign.Read8(xinput, uintptr(32*2)) + unalign.Read8(xinput, uintptr(32*2)+8)

		accLow64 += mix(unalign.Read8(xinput, uintptr(32*3))^xsecret_096, unalign.Read8(xinput, uintptr(32*3+8))^xsecret_104)
		accLow64 ^= unalign.Read8(xinput, uintptr(32*3+16)) + unalign.Read8(xinput, uintptr(32*3+24))
		accHigh64 += mix(unalign.Read8(xinput, uintptr(32*3+16))^xsecret_112, unalign.Read8(xinput, uintptr(32*3)+24)^xsecret_120)
		accHigh64 ^= unalign.Read8(xinput, uintptr(32*3)) + unalign.Read8(xinput, uintptr(32*3)+8)

		accLow64 = xxh3Avalanche(accLow64)
		accHigh64 = xxh3Avalanche(accHigh64)

		nbRounds := length >>5 << 5
		for i := uint64(4*32); i < nbRounds; i+=32 {
			accHigh64 += mix(unalign.Read8(xinput, uintptr(i+16))^unalign.Read8(xsecret, uintptr(i-109)), unalign.Read8(xinput, uintptr(i)+24)^unalign.Read8(xsecret, uintptr(i-101)))
			accHigh64 ^= unalign.Read8(xinput, uintptr(i)) + unalign.Read8(xinput, uintptr(i)+8)

			accLow64 += mix(unalign.Read8(xinput, uintptr(i))^unalign.Read8(xsecret, uintptr(i-125)), unalign.Read8(xinput, uintptr(i+8))^unalign.Read8(xsecret, uintptr(i-117)))
			accLow64 ^= unalign.Read8(xinput, uintptr(i+16)) + unalign.Read8(xinput, uintptr(i+24))
		}

		// last 32 bytes
		accLow64 += mix(unalign.Read8(xinput, uintptr(length-16))^xsecret_103, unalign.Read8(xinput, uintptr(length-8))^xsecret_111)
		accLow64 ^= unalign.Read8(xinput, uintptr(length-32)) + unalign.Read8(xinput, uintptr(length-24))
		accHigh64 += mix(unalign.Read8(xinput, uintptr(length-32))^xsecret_119, unalign.Read8(xinput, uintptr(length-24))^xsecret_127)
		accHigh64 ^= unalign.Read8(xinput, uintptr(length-16)) + unalign.Read8(xinput, uintptr(length-8))

		accHigh64, accLow64 = (accLow64*prime64_1)+(accHigh64*prime64_4)+(length*prime64_2), accHigh64+accLow64

		accLow64 = xxh3Avalanche(accLow64)
		accHigh64 = -xxh3Avalanche(accHigh64)

		return [2]uint64{accHigh64, accLow64}
	}

	xacc = [8]uint64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	acc[1] = length * prime64_1
	acc[0] = ^(length * prime64_2)

	if avx2 {
		accumAVX2(&xacc, xinput, xsecret, length)
	} else if sse2 {
		accumSSE2(&xacc, xinput, xsecret, length)
	} else {
		accumScalar(&xacc, xinput, xsecret, length)
	}
	// merge xacc
	acc[1] += mix(xacc[0]^xsecret_011, xacc[1]^xsecret_019)
	acc[1] += mix(xacc[2]^xsecret_027, xacc[3]^xsecret_035)
	acc[1] += mix(xacc[4]^xsecret_043, xacc[5]^xsecret_051)
	acc[1] += mix(xacc[6]^xsecret_059, xacc[7]^xsecret_067)
	acc[1] = xxh3Avalanche(acc[1])

	acc[0] += mix(xacc[0]^xsecret_117, xacc[1]^xsecret_125)
	acc[0] += mix(xacc[2]^xsecret_133, xacc[3]^xsecret_141)
	acc[0] += mix(xacc[4]^xsecret_149, xacc[5]^xsecret_157)
	acc[0] += mix(xacc[6]^xsecret_165, xacc[7]^xsecret_173)
	acc[0] = xxh3Avalanche(acc[0])

	return acc
}
