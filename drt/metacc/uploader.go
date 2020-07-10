package metacc

import "github.com/zaakkk/libdrt/drt/core"

//Uploader はメタデータの送信に関するインターフェースを規定する
//実装例はCustomにある
type Uploader interface {
	//listにあるメタデータを送信し、アクセスするのに必要な情報(アクセスキー)の配列を返さなければならない
	//ただし、アクセスキーはlistの順で返さなければならない
	Upload(list []core.Part) []string
}
