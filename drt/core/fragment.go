package core

//Fragment は断片データを表現する
type Fragment struct {
	Buffer []byte //担当するデータ領域
	Hash   []byte //ハッシュ値
	DestID uint8  //送信先ID
	Dest   string //送信先
	Prefix string //接頭辞
	Order  uint8  //送信先ごとの並び順
}
