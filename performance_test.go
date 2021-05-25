package xxhash3

import (
	"../xxh3"
	"math/rand"
	"runtime"
	"testing"
)

var dat []byte

const capacity = 33554432 + 100000

func init() {
	dat = make([]byte, capacity)
	for i := 0; i < capacity; i++ {
		dat[i] = byte(rand.Int31())
	}
}

type benchTask struct {
	name   string
	action func(b []byte) uint64
}

func BenchmarkDefault(b *testing.B) {
	all := []benchTask{{
		name: "Target", action: func(b []byte) uint64 {
			return Hash(b)
		}}, {
		name: "Baseline", action: func(b []byte) uint64 {
			return xxh3.Hash(b)
		}},
	}

	benchLen0_16(b, all)
	benchLen17_128(b, all)
	benchLen129_240(b, all)
	benchLen241_1024(b, all)
	benchScalar(b, all)
	benchAVX2(b, all)
	benchSSE2(b, all)

}

func benchLen0_16(b *testing.B, benchTasks []benchTask) {
	for _, v := range benchTasks {
		b.Run("Len0_16/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				for i := 0; i <= 16; i++ {
					input := dat[b.N : b.N+i]
					a := v.action(input)
					runtime.KeepAlive(a)
				}
			}
		})
	}
}

func benchLen17_128(b *testing.B, benchTasks []benchTask) {
	for _, v := range benchTasks {
		b.Run("Len17_128/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				for i := 17; i <= 128; i++ {
					input := dat[b.N : b.N+i]
					a := v.action(input)
					runtime.KeepAlive(a)
				}
			}
		})
	}
}

func benchLen129_240(b *testing.B, benchTasks []benchTask) {
	for _, v := range benchTasks {
		b.Run("Len129_240/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				for i := 129; i <= 240; i++ {
					input := dat[b.N : b.N+i]
					a := v.action(input)
					runtime.KeepAlive(a)
				}
			}
		})
	}
}

func benchLen241_1024(b *testing.B, benchTasks []benchTask) {
	avx2, sse2 = false, false
	xxh3.Avx2, xxh3.Sse2 = false, false
	for _, v := range benchTasks {
		b.Run("Len241_1024/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				for i := 241; i <= 1024; i++ {
					input := dat[b.N : b.N+i]
					a := v.action(input)
					runtime.KeepAlive(a)
				}
			}
		})
	}
}

func benchScalar(b *testing.B, benchTasks []benchTask) {
	avx2, sse2 = false, false
	xxh3.Avx2, xxh3.Sse2 = false, false

	for _, v := range benchTasks {
		b.Run("Scalar/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				input := dat[n:33554432]
				a := v.action(input)
				runtime.KeepAlive(a)

			}
		})
	}
}

func benchAVX2(b *testing.B, benchTasks []benchTask) {
	avx2, sse2 = true, false
	xxh3.Avx2, xxh3.Sse2 = true, false

	for _, v := range benchTasks {
		b.Run("AVX2/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				input := dat[n:33554432]
				a := v.action(input)
				runtime.KeepAlive(a)
			}
		})
	}
}
func benchSSE2(b *testing.B, benchTasks []benchTask) {
	avx2, sse2 = false, true
	xxh3.Avx2, xxh3.Sse2 = false, true

	for _, v := range benchTasks {
		b.Run("SSE2/"+v.name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				input := dat[n:33554432]
				a := v.action(input)
				runtime.KeepAlive(a)
			}
		})
	}
}
