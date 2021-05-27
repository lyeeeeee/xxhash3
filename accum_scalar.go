package xxhash3

import (
	"unsafe"
	"./internal/unalign"
)
var(
key64_128 uint64 = 0xc3ebd33483acc5ea
key64_136 uint64 = 0xeb6313faffa081c5
key64_144 uint64 = 0x49daf0b751dd0d17
key64_152 uint64 = 0x9e68d429265516d3
key64_160 uint64 = 0xfca1477d58be162b
key64_168 uint64 = 0xce31d07ad1b8f88f
key64_176 uint64 = 0x280416958f3acb45
key64_184 uint64 = 0x7e404bbbcafbd7af

	key64_121 uint64 = 0xea647378d9c97e9f
	key64_129 uint64 = 0xc5c3ebd33483acc5
	key64_137 uint64 = 0x17eb6313faffa081
	key64_145 uint64 = 0xd349daf0b751dd0d
	key64_153 uint64 = 0x2b9e68d429265516
	key64_161 uint64 = 0x8ffca1477d58be16
	key64_169 uint64 = 0x45ce31d07ad1b8f8
	key64_177 uint64 = 0xaf280416958f3acb

	key64_011 uint64 = 0x6dd4de1cad21f72c
	key64_019 uint64 = 0xa44072db979083e9
	key64_027 uint64 = 0xe679cb1f67b3b7a4
	key64_035 uint64 = 0xd05a8278e5c0cc4e
	key64_043 uint64 = 0x4608b82172ffcc7d
	key64_051 uint64 = 0x9035e08e2443f774
	key64_059 uint64 = 0x52283c4c263a81e6
	key64_067 uint64 = 0x65d088cb00c391bb
)
func accumScalar(accs *[8]uint64, p, key unsafe.Pointer, l uint64) {
	for l > _block {
		k := key

		// accs
		for i := 0; i < 16; i++ {
			dv0 := unalign.Read8(p, 8*0)
			dk0 := dv0 ^ unalign.Read8(k, 8*0)
			accs[1] += dv0
			accs[0] += (dk0 & 0xffffffff) * (dk0 >> 32)

			dv1 := unalign.Read8(p, 8*1)
			dk1 := dv1 ^ unalign.Read8(k, 8*1)
			accs[0] += dv1
			accs[1] += (dk1 & 0xffffffff) * (dk1 >> 32)

			dv2 := unalign.Read8(p, 8*2)
			dk2 := dv2 ^ unalign.Read8(k, 8*2)
			accs[3] += dv2
			accs[2] += (dk2 & 0xffffffff) * (dk2 >> 32)

			dv3 := unalign.Read8(p, 8*3)
			dk3 := dv3 ^ unalign.Read8(k, 8*3)
			accs[2] += dv3
			accs[3] += (dk3 & 0xffffffff) * (dk3 >> 32)

			dv4 := unalign.Read8(p, 8*4)
			dk4 := dv4 ^ unalign.Read8(k, 8*4)
			accs[5] += dv4
			accs[4] += (dk4 & 0xffffffff) * (dk4 >> 32)

			dv5 := unalign.Read8(p, 8*5)
			dk5 := dv5 ^ unalign.Read8(k, 8*5)
			accs[4] += dv5
			accs[5] += (dk5 & 0xffffffff) * (dk5 >> 32)

			dv6 := unalign.Read8(p, 8*6)
			dk6 := dv6 ^ unalign.Read8(k, 8*6)
			accs[7] += dv6
			accs[6] += (dk6 & 0xffffffff) * (dk6 >> 32)

			dv7 := unalign.Read8(p, 8*7)
			dk7 := dv7 ^ unalign.Read8(k, 8*7)
			accs[6] += dv7
			accs[7] += (dk7 & 0xffffffff) * (dk7 >> 32)

			l -= _stripe
			if l > 0 {
				p, k = unsafe.Pointer(uintptr(p)+_stripe), unsafe.Pointer(uintptr(k)+8)
			}
		}

		// scramble accs
		accs[0] ^= accs[0] >> 47
		accs[0] ^= key64_128
		accs[0] *= prime32_1

		accs[1] ^= accs[1] >> 47
		accs[1] ^= key64_136
		accs[1] *= prime32_1

		accs[2] ^= accs[2] >> 47
		accs[2] ^= key64_144
		accs[2] *= prime32_1

		accs[3] ^= accs[3] >> 47
		accs[3] ^= key64_152
		accs[3] *= prime32_1

		accs[4] ^= accs[4] >> 47
		accs[4] ^= key64_160
		accs[4] *= prime32_1

		accs[5] ^= accs[5] >> 47
		accs[5] ^= key64_168
		accs[5] *= prime32_1

		accs[6] ^= accs[6] >> 47
		accs[6] ^= key64_176
		accs[6] *= prime32_1

		accs[7] ^= accs[7] >> 47
		accs[7] ^= key64_184
		accs[7] *= prime32_1
	}

	if l > 0 {
		t, k := (l-1)/_stripe, key

		for i := uint64(0); i < t; i++ {
			dv0 := unalign.Read8(p, 8*0)
			dk0 := dv0 ^ unalign.Read8(k, 8*0)
			accs[1] += dv0
			accs[0] += (dk0 & 0xffffffff) * (dk0 >> 32)

			dv1 := unalign.Read8(p, 8*1)
			dk1 := dv1 ^ unalign.Read8(k, 8*1)
			accs[0] += dv1
			accs[1] += (dk1 & 0xffffffff) * (dk1 >> 32)

			dv2 := unalign.Read8(p, 8*2)
			dk2 := dv2 ^ unalign.Read8(k, 8*2)
			accs[3] += dv2
			accs[2] += (dk2 & 0xffffffff) * (dk2 >> 32)

			dv3 := unalign.Read8(p, 8*3)
			dk3 := dv3 ^ unalign.Read8(k, 8*3)
			accs[2] += dv3
			accs[3] += (dk3 & 0xffffffff) * (dk3 >> 32)

			dv4 := unalign.Read8(p, 8*4)
			dk4 := dv4 ^ unalign.Read8(k, 8*4)
			accs[5] += dv4
			accs[4] += (dk4 & 0xffffffff) * (dk4 >> 32)

			dv5 := unalign.Read8(p, 8*5)
			dk5 := dv5 ^ unalign.Read8(k, 8*5)
			accs[4] += dv5
			accs[5] += (dk5 & 0xffffffff) * (dk5 >> 32)

			dv6 := unalign.Read8(p, 8*6)
			dk6 := dv6 ^ unalign.Read8(k, 8*6)
			accs[7] += dv6
			accs[6] += (dk6 & 0xffffffff) * (dk6 >> 32)

			dv7 := unalign.Read8(p, 8*7)
			dk7 := dv7 ^ unalign.Read8(k, 8*7)
			accs[6] += dv7
			accs[7] += (dk7 & 0xffffffff) * (dk7 >> 32)

			l -= _stripe
			if l > 0 {
				p, k = unsafe.Pointer(uintptr(p)+_stripe), unsafe.Pointer(uintptr(k)+8)
			}
		}

		if l > 0 {
			p = unsafe.Pointer(uintptr(p) - uintptr(_stripe-l))

			dv0 := unalign.Read8(p, 8*0)
			dk0 := dv0 ^ key64_121
			accs[1] += dv0
			accs[0] += (dk0 & 0xffffffff) * (dk0 >> 32)

			dv1 := unalign.Read8(p, 8*1)
			dk1 := dv1 ^ key64_129
			accs[0] += dv1
			accs[1] += (dk1 & 0xffffffff) * (dk1 >> 32)

			dv2 := unalign.Read8(p, 8*2)
			dk2 := dv2 ^ key64_137
			accs[3] += dv2
			accs[2] += (dk2 & 0xffffffff) * (dk2 >> 32)

			dv3 := unalign.Read8(p, 8*3)
			dk3 := dv3 ^ key64_145
			accs[2] += dv3
			accs[3] += (dk3 & 0xffffffff) * (dk3 >> 32)

			dv4 := unalign.Read8(p, 8*4)
			dk4 := dv4 ^ key64_153
			accs[5] += dv4
			accs[4] += (dk4 & 0xffffffff) * (dk4 >> 32)

			dv5 := unalign.Read8(p, 8*5)
			dk5 := dv5 ^ key64_161
			accs[4] += dv5
			accs[5] += (dk5 & 0xffffffff) * (dk5 >> 32)

			dv6 := unalign.Read8(p, 8*6)
			dk6 := dv6 ^ key64_169
			accs[7] += dv6
			accs[6] += (dk6 & 0xffffffff) * (dk6 >> 32)

			dv7 := unalign.Read8(p, 8*7)
			dk7 := dv7 ^ key64_177
			accs[6] += dv7
			accs[7] += (dk7 & 0xffffffff) * (dk7 >> 32)
		}
	}
}
