package core

//Origin は元データを表現する
type Origin struct {
	Addition []byte //追加データ；ファイル名やを表す
	Buffer   []byte //データ領域；パディング用に若干多くキャパシティを設けることが望ましい
}

//NewOrigin はOriginのコンストラクタ
func NewOrigin(add []byte, buf []byte) *Origin {
	o := new(Origin)
	o.Addition = add
	o.Buffer = buf
	return o
}
