// Package t1ha implements the t1ha hash function
/*

   https://github.com/leo-yuriev/t1ha

*/
package t1ha

import "encoding/binary"

const (
	// magic primes
	p0 = 17048867929148541611
	p1 = 9386433910765580089
	p2 = 15343884574428479051
	p3 = 13662985319504319857
	p4 = 11242949449147999147
	p5 = 13862205317416547141
	p6 = 14653293970879851569

	/* rotations */
	s0 = 41
	s1 = 17
	s2 = 31
)

func rot64(v uint64, s uint) uint64 {
	return (v >> s) | (v << (64 - s))
}

func fetch64(v []byte, idx int) uint64 {
	return binary.LittleEndian.Uint64(v[idx*8:])
}

func fetch_tail(v []byte) uint64 {
	switch len(v) {
	case 1:
		return uint64(v[0])
	case 2:
		return uint64(binary.LittleEndian.Uint16(v))
	case 3:
		return uint64(binary.LittleEndian.Uint16(v)) | uint64(v[2])<<16
	case 4:
		return uint64(binary.LittleEndian.Uint32(v))
	case 5:
		return uint64(binary.LittleEndian.Uint32(v)) | uint64(v[4])<<32
	case 6:
		return uint64(binary.LittleEndian.Uint32(v)) | uint64(binary.LittleEndian.Uint16(v[4:]))<<32
	case 7:
		return uint64(binary.LittleEndian.Uint32(v)) | uint64(binary.LittleEndian.Uint16(v[4:]))<<32 | uint64(v[6])<<48
	case 8:
		return binary.LittleEndian.Uint64(v)
	}

	panic("not reached")
}

func mux64(a, b uint64) uint64 {
	lo, hi := mullu(a, b)
	return lo ^ hi
}

// xor-mul-xor mixer
func mix(v, p uint64) uint64 {
	v *= p
	return v ^ rot64(v, s0)
}

func Sum64(data []byte, seed uint64) uint64 {
	var a uint64 = seed
	var b uint64 = uint64(len(data))

	if len(data) > 32 {
		c := rot64(uint64(len(data)), s1) + seed
		d := uint64(len(data)) ^ rot64(seed, s1)

		for len(data) >= 32 {
			v := data

			w0 := fetch64(v, 0)
			w1 := fetch64(v, 1)
			w2 := fetch64(v, 2)
			w3 := fetch64(v, 3)

			d02 := w0 ^ rot64(w2+d, s1)
			c13 := w1 ^ rot64(w3+c, s1)
			c += a ^ rot64(w0, s0)
			d -= b ^ rot64(w1, s2)
			a ^= p1 * (d02 + w3)
			b ^= p0 * (c13 + w2)
			data = data[8*4:]
		}

		a ^= p6 * (rot64(c, s1) + d)
		b ^= p5 * (c + rot64(d, s1))
	}

	v := data

	switch len(v) {
	default:
		b += mux64(fetch64(v, 0), p4)
		v = v[8:]
		fallthrough
	case 24:
		fallthrough
	case 23:
		fallthrough
	case 22:
		fallthrough
	case 21:
		fallthrough
	case 20:
		fallthrough
	case 19:
		fallthrough
	case 18:
		fallthrough
	case 17:
		a += mux64(fetch64(v, 0), p3)
		v = v[8:]
		fallthrough
	case 16:
		fallthrough
	case 15:
		fallthrough
	case 14:
		fallthrough
	case 13:
		fallthrough
	case 12:
		fallthrough
	case 11:
		fallthrough
	case 10:
		fallthrough
	case 9:
		b += mux64(fetch64(v, 0), p2)
		v = v[8:]
		fallthrough
	case 8:
		fallthrough
	case 7:
		fallthrough
	case 6:
		fallthrough
	case 5:
		fallthrough
	case 4:
		fallthrough
	case 3:
		fallthrough
	case 2:
		fallthrough
	case 1:
		a += mux64(fetch_tail(v), p1)
	case 0:
	}

	return mux64(rot64(a+b, s1), p4) + mix(a^b, p0)
}

// From https://golang.org/src/runtime/softfloat64.go
// 64x64 -> 128 multiply.
// adapted from hacker's delight.
func mullu(u, v uint64) (lo, hi uint64) {
	const (
		s    = 32
		mask = 1<<s - 1
	)
	u0 := u & mask
	u1 := u >> s
	v0 := v & mask
	v1 := v >> s
	w0 := u0 * v0
	t := u1*v0 + w0>>s
	w1 := t & mask
	w2 := t >> s
	w1 += u0 * v1
	return u * v, u1*v1 + w2 + w1>>s
}
