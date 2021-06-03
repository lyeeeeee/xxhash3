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

		h128High64, h128Low64 := bits.Mul64(m128Low64, prime64_2)
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

		bitflipl := uint64(xsecret32_000 ^ xsecret32_004)
		bitfliph := uint64(xsecret32_008^ xsecret32_012)

		keyedLow := combinedl ^ bitflipl
		keyedHigh := combinedh ^ bitfliph

		keyedLow = combinedl ^ bitflipl
		keyedHigh = combinedh ^ bitfliph

		h128Low64 := xxh64Avalanche(keyedLow)
		h128High64 := xxh64Avalanche(keyedHigh)
		return [2]uint64{h128High64, h128Low64}
	}

	return [2]uint64{0x99aa06d3014798d8, 0x6001c324468d497f}
}

func xxh3HashLarge128(xinput unsafe.Pointer, l int) (acc [2]uint64) {
	length := uint64(l)

	if length <= 128 {

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
	} else if length <= 240 {
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

	xacc = [8]uint64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	acc[1] = length * prime64_1
	acc[0] = ^(length * prime64_2)

	accumScalar(&xacc, xinput, xsecret, length)
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