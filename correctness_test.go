package xxhash3

import (
	"../xxh3"
	"testing"
)

func TestLen0_16(t *testing.T) {
	for i := 0; i <= 16; i++ {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen17_128(t *testing.T) {
	for i := 17; i <= 128; i++ {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen129_240(t *testing.T) {

	for i := 129; i <= 240; i++ {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen240_1024(t *testing.T) {
	avx2, sse2 = false, false

	for i := 240; i <= 1024; i++ {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen1025_1048576_scalar(t *testing.T) {
	avx2, sse2 = false, false
	for i := 1025; i < 1048576; i = i << 1 {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer", i)
		}
	}
}

func TestLen1024_1048576_AVX2(t *testing.T) {
	avx2, sse2 = true, false

	for i := 1024; i < 1048576; i = i << 1 {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer", i)
		}
	}
}

func TestLen1024_1048576_SSE2(t *testing.T) {
	avx2, sse2 = false, true

	for i := 1024; i < 1048576; i = i << 1 {
		input := dat[:i]
		res1 := xxh3.Hash(input)
		res2 := Hash(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen128_0_16(t *testing.T) {
	for i := 0; i <= 16; i++ {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen128_17_128(t *testing.T) {
	for i := 17; i <= 128; i++ {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen128_129_240(t *testing.T) {

	for i := 129; i <= 240; i++ {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen128_240_1024(t *testing.T) {
	avx2, sse2 = false, false

	for i := 240; i <= 1024; i++ {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}

func TestLen128_1025_1048576_scalar(t *testing.T) {
	avx2, sse2 = false, false
	for i := 1025; i < 1048576; i = i << 1 {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer", i)
		}
	}
}

func TestLen128_1024_1048576_AVX2(t *testing.T) {
	avx2, sse2 = true, false

	for i := 1024; i < 1048576; i = i << 1 {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer", i)
		}
	}
}

func TestLen128_1024_1048576_SSE2(t *testing.T) {
	avx2, sse2 = false, true

	for i := 1024; i < 1048576; i = i << 1 {
		input := dat[:i]
		res1 := xxh3.Hash128(input)
		res2 := Hash128(input)

		if res1 != res2 {
			t.Fatal("Wrong answer")
		}
	}
}
