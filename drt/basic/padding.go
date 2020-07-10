package basic

import (
	"github.com/zaakkk/libdrt/drt/core"
)

//MinPadding は最低限の実装を持つ
type MinPadding struct {
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

//Encrypt は16と分割数の公倍数になるようにパディングをする
//パフォーマンスはあまりよくない
func (s *MinPadding) Encrypt(buf []byte, m *core.Metadata) []byte {
	gcm := lcm(16, int(m.Division))
	padSize := (gcm - m.DataSize%gcm) % gcm
	for i := 0; i < padSize; i++ {
		buf = append(buf, 0xFF)
	}
	return buf
}

//Decrypt は元ファイルの大きさに戻す
func (s *MinPadding) Decrypt(buf []byte, m *core.Metadata) []byte {
	return buf[:m.DataSize]
}
