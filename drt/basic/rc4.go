package basic

import (
	"crypto/rc4"

	"../core"
)

//RC4 はRC4を実装する
//他のシステムのRC4と互換性があるかは若干疑わしい
type RC4 struct {
}

//Encrypt 暗号化
func (c *RC4) Encrypt(buf []byte, m *core.Metadata) []byte {
	ci, err := rc4.NewCipher(m.Scramble)
	if err != nil {
		panic(err)
	}
	ci.XORKeyStream(buf, buf)
	return buf
}

//Decrypt 復号
//ストリーム暗号なので暗号化と同じ処理
func (c *RC4) Decrypt(buf []byte, m *core.Metadata) []byte {
	return c.Encrypt(buf, m)
}
