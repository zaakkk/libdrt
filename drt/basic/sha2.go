package basic

import (
	"crypto/sha256"
	"crypto/sha512"
)

//SHA256 はハッシュの計算を標準ライブラリに丸投げ
type SHA256 struct {
}

//Calc はハッシュの計算
func (h *SHA256) Calc(buf []byte) []byte {
	hash := sha256.Sum256([]byte(buf))
	return hash[:]
}

//SHA512 is SHA512
//faster than SHA256
type SHA512 struct {
}

//Calc はハッシュの計算
func (h *SHA512) Calc(buf []byte) []byte {
	hash := sha512.Sum512([]byte(buf))
	return hash[:]
}
