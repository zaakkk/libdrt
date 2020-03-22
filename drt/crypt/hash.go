package crypt

//Hash はハッシュ関数や誤り検知符号などのインターフェースを定義する
type Hash interface {
	//Calc はbufの中身のハッシュを計算して返す
	Calc(buf []byte) []byte
}
