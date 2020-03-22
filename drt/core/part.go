package core

//Part は閾値秘密分散したファイルを表す
type Part struct {
	Buffer []byte
	Dest   string
}
