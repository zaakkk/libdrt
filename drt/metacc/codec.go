package metacc

import "../core"

//Codec はメタデータのシリアライズに関するインターフェースを定義する
type Codec interface {
	//Write はメタデータを受け取りシリアライズして返す
	Write(m *core.Metadata) []byte
	//Read はバイト列を受け取り、デシリアライズしてメタデータを構築して返す
	Read(buf []byte) *core.Metadata
}
