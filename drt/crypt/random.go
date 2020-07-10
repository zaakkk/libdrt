package crypt

import (
	"crypto/rand"
	"math/big"
)

//CreateRandomByte は乱数を1byteを0~maxの間で生成する
func CreateRandomByte(max int) byte {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic("failed to create random value")
	}
	if n.Cmp(big.NewInt(int64(0))) == 0 {
		return 0
	}
	return n.Bytes()[0]
}

//CreateRandomBytes は乱数列を生成する
func CreateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(255))
		if err != nil {
			panic("failed to create random value")
		}
		if n.Cmp(big.NewInt(int64(0))) == 0 {
			bytes[i] = 0
		} else {
			bytes[i] = n.Bytes()[0]
		}
	}
	return bytes
}

//RandomOrder は安全な方法でbyte列を並び替える
func RandomOrder(size uint8) []byte {
	order := make([]byte, size)
	for i := 0; i < len(order); i++ {
		order[i] = byte(i)
	}
	for i := len(order) - 1; i > 0; i-- {
		index := CreateRandomByte(i)
		temp := order[i]
		order[i] = order[index]
		order[index] = temp
	}
	return order
}

//Shuffle バッファをシャッフルする
func Shuffle(bufs [][]byte, order []byte) [][]byte {
	if len(bufs) != len(order) {
		panic("Shuffle encounted illegal argument")
	}
	size := len(order)
	for i := 0; i < size; i++ {
		index := order[i]
		temp := bufs[index]
		bufs[index] = bufs[i]
		bufs[i] = temp
	}
	return bufs
}

//Sort バッファをソートする
func Sort(bufs [][]byte, order []byte) [][]byte {
	if len(bufs) != len(order) {
		panic("Shuffle encounted illegal argument")
	}
	size := len(order)
	for i := size - 1; i >= 0; i-- {
		index := order[i]
		temp := bufs[index]
		bufs[index] = bufs[i]
		bufs[i] = temp
	}
	return bufs
}
