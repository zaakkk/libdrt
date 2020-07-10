package metacc

import "github.com/zaakkk/libdrt/drt/core"

//Downloader はメタデータをダウンロードするためのインターフェースを定義する
type Downloader interface {
	//Download はaccessKeysの中のアクセスキーを利用し、メタデータをダウンロードする
	//すべてのメタデータをダウンロードする必要はなく、listの順番は関係ない
	//メタデータが足りないかどうかは別の関数が判断するため、この関数が判断しなくて良い
	Download(list []core.Part, accessKeys []string) []core.Part
}
