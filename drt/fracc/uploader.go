package fracc

import "../core"

//Uploader は断片データ送信のインターフェースを定義する
//実装例はcustomパッケージにある
type Uploader interface {
	//Uploadはtableの中の断片データをすべて送信しなければならない
	Upload(table [][]core.Fragment)
}