package container

//Divide はバッファを分割して返す
//余りが有っても扱うことができる
func Divide(buf []byte, dn uint8) [][]byte {
	size := len(buf) / int(dn)
	bufs := make([][]byte, dn)
	for i := uint8(0); i < dn; i++ {
		begin := size * int(i)
		end := begin + size
		if i == dn-1 {
			bufs[int(i)] = buf[begin:len(buf)]
		} else {
			bufs[int(i)] = buf[begin:end]
		}
	}
	return bufs
}

//Bundle はバッファを分割して返す
//余りが有っても扱うことができる
func Bundle(bufs [][]byte) []byte {
	dn := uint8(len(bufs))
	size := int(dn) * len(bufs[0])
	buf := make([]byte, size)
	for i := uint8(0); i < dn; i++ {
		begin := size * int(i)
		end := begin + size
		if i == dn-1 {
			bufs[int(i)] = buf[begin:len(buf)]
		} else {
			bufs[int(i)] = buf[begin:end]
		}
	}
	return buf
}
