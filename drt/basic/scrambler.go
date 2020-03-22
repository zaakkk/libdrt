package basic

import (
	"reflect"
	"unsafe"

	"../core"
)

//Scrambler は基本的な一体化処理を行うが
//バッファの長さが8で割り切れない場合は処理できず、panic()する
type Scrambler struct {
}

func (s *Scrambler) canProcess(buf []byte) {
	if len(buf)%8 > 0 {
		panic("Scrambler ecnounted illegal argument; Buffer's length must be able be divided by 8")
	}
}

//Encrypt は暗号化
func (s *Scrambler) Encrypt(buf []byte, m *core.Metadata) []byte {
	s.canProcess(buf)
	ptr := unsafe.Pointer(&buf)
	sl := *(*reflect.SliceHeader)(ptr)
	p := sl.Data
	l := uintptr(sl.Len)
	for i := 0; i < len(m.Scramble); i++ {
		operattion := m.Scramble[i] % 12
		switch operattion {
		case 0:
			encrypt8A(p, l)
		case 1:
			encrypt16A(p, l)
		case 2:
			encrypt32A(p, l)
		case 3:
			encrypt64A(p, l)
		case 4:
			encrypt8S(p, l)
		case 5:
			encrypt16S(p, l)
		case 6:
			encrypt32S(p, l)
		case 7:
			encrypt64S(p, l)
		case 8:
			encrypt8X(p, l)
		case 9:
			encrypt16X(p, l)
		case 10:
			encrypt32X(p, l)
		case 11:
			encrypt64X(p, l)
		}
	}
	return buf
}

//Decrypt は復号
func (s *Scrambler) Decrypt(buf []byte, m *core.Metadata) []byte {
	s.canProcess(buf)
	ptr := unsafe.Pointer(&buf)
	sl := *(*reflect.SliceHeader)(ptr)
	p := sl.Data
	l := uintptr(sl.Len)
	for i := len(m.Scramble) - 1; i >= 0; i-- {
		operattion := m.Scramble[i] % 12
		switch operattion {
		case 0:
			decrypt8A(p, l)
		case 1:
			decrypt16A(p, l)
		case 2:
			decrypt32A(p, l)
		case 3:
			decrypt64A(p, l)
		case 4:
			decrypt8S(p, l)
		case 5:
			decrypt16S(p, l)
		case 6:
			decrypt32S(p, l)
		case 7:
			decrypt64S(p, l)
		case 8:
			decrypt8X(p, l)
		case 9:
			decrypt16X(p, l)
		case 10:
			decrypt32X(p, l)
		case 11:
			decrypt64X(p, l)
		}
	}
	return buf
}

func toSlice64(p uintptr, len uintptr) []uint64 {
	const size = 8
	var sh reflect.SliceHeader
	sh.Data = p
	sh.Len = int(len) / size
	sh.Cap = sh.Len
	return *(*[]uint64)(unsafe.Pointer(&sh))
}

func encrypt64A(p uintptr, l uintptr) {
	buf := toSlice64(p, l)
	var r1, r2 uint64
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 += r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt64A(p uintptr, l uintptr) {
	buf := toSlice64(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] -= buf[i-1]
	}
	buf[0] -= buf[size-1]
}

func encrypt64S(p uintptr, l uintptr) {
	buf := toSlice64(p, l)
	var r1, r2 uint64
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 -= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt64S(p uintptr, l uintptr) {
	buf := toSlice64(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] += buf[i-1]
	}
	buf[0] += buf[size-1]
}

func encrypt64X(p uintptr, l uintptr) {
	buf := toSlice64(p, l)
	var r1, r2 uint64
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 ^= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt64X(p uintptr, l uintptr) {
	buf := toSlice64(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] ^= buf[i-1]
	}
	buf[0] ^= buf[size-1]
}

func toSlice32(p uintptr, len uintptr) []uint32 {
	const size = 4
	var sh reflect.SliceHeader
	sh.Data = p
	sh.Len = int(len) / size
	sh.Cap = sh.Len
	return *(*[]uint32)(unsafe.Pointer(&sh))
}

func encrypt32A(p uintptr, l uintptr) {
	buf := toSlice32(p, l)
	var r1, r2 uint32
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 += r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt32A(p uintptr, l uintptr) {
	buf := toSlice32(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] -= buf[i-1]
	}
	buf[0] -= buf[size-1]
}

func encrypt32S(p uintptr, l uintptr) {
	buf := toSlice32(p, l)
	var r1, r2 uint32
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 -= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt32S(p uintptr, l uintptr) {
	buf := toSlice32(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] += buf[i-1]
	}
	buf[0] += buf[size-1]
}

func encrypt32X(p uintptr, l uintptr) {
	buf := toSlice32(p, l)
	var r1, r2 uint32
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 ^= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt32X(p uintptr, l uintptr) {
	buf := toSlice32(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] ^= buf[i-1]
	}
	buf[0] ^= buf[size-1]
}

func toSlice16(p uintptr, len uintptr) []uint16 {
	const size = 2
	var sh reflect.SliceHeader
	sh.Data = p
	sh.Len = int(len) / size
	sh.Cap = sh.Len
	return *(*[]uint16)(unsafe.Pointer(&sh))
}

func encrypt16A(p uintptr, l uintptr) {
	buf := toSlice16(p, l)
	var r1, r2 uint16
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 += r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt16A(p uintptr, l uintptr) {
	buf := toSlice16(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] -= buf[i-1]
	}
	buf[0] -= buf[size-1]
}

func encrypt16S(p uintptr, l uintptr) {
	buf := toSlice16(p, l)
	var r1, r2 uint16
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 -= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt16S(p uintptr, l uintptr) {
	buf := toSlice16(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] += buf[i-1]
	}
	buf[0] += buf[size-1]
}

func encrypt16X(p uintptr, l uintptr) {
	buf := toSlice16(p, l)
	var r1, r2 uint16
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 ^= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt16X(p uintptr, l uintptr) {
	buf := toSlice16(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] ^= buf[i-1]
	}
	buf[0] ^= buf[size-1]
}

func toSlice8(p uintptr, len uintptr) []uint8 {
	const size = 1
	var sh reflect.SliceHeader
	sh.Data = p
	sh.Len = int(len) / size
	sh.Cap = sh.Len
	return *(*[]uint8)(unsafe.Pointer(&sh))
}

func encrypt8A(p uintptr, l uintptr) {
	buf := toSlice8(p, l)
	var r1, r2 uint8
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 += r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt8A(p uintptr, l uintptr) {
	buf := toSlice8(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] -= buf[i-1]
	}
	buf[0] -= buf[size-1]
}

func encrypt8S(p uintptr, l uintptr) {
	buf := toSlice8(p, l)
	var r1, r2 uint8
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 -= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt8S(p uintptr, l uintptr) {
	buf := toSlice8(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] += buf[i-1]
	}
	buf[0] += buf[size-1]
}

func encrypt8X(p uintptr, l uintptr) {
	buf := toSlice8(p, l)
	var r1, r2 uint8
	size := len(buf)
	r2 = buf[size-1]
	for i := 0; i < size; i++ {
		r1 = buf[i]
		r1 -= r2
		buf[i] = r1
		r2 = r1
	}
}

func decrypt8X(p uintptr, l uintptr) {
	buf := toSlice8(p, l)
	size := len(buf)
	for i := size - 1; i > 0; i-- {
		buf[i] += buf[i-1]
	}
	buf[0] += buf[size-1]
}
