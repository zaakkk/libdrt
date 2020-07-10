package crypt

import "github.com/zaakkk/libdrt/drt/core"

//Cipher は暗号のインターフェースを定義する
type Cipher interface {
	//bufを暗号化して返す
	Encrypt(buf []byte, m *core.Metadata) []byte
	//bufを復号して返す
	Decrypt(buf []byte, m *core.Metadata) []byte
}
