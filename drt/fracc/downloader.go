package fracc

import (
	"../core"
	"../crypt"
)

//Downloader は断片データ受信のインターフェースを定義する
//実装例はcustomパッケージにある
type Downloader interface {
	//Downloadはtableの中の断片データを複製された物の中から一つずつ受信しなければならない
	Download(table [][]core.Fragment, hash crypt.Hash)
}
