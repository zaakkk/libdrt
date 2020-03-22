package core

//Metadata は閾値秘密分散前のメタデータを表現する
type Metadata struct {
	Addition       []byte   //追加データ；ファイル名とか
	DataSize       int      //元データの大きさ
	OriginHash     []byte   //元データのハッシュ
	Division       uint8    //分割数
	Replication    uint8    //複製数
	Scramble       []byte   //一体化鍵
	Randomize      []byte   //ランダマイズ鍵
	Order          []byte   //シャッフルの順番
	Destinations   []string //断片データ送信先
	FragmentSize   int      //断片データの大きさ
	FragmentHash   [][]byte //断片データのハッシュ；複製してもハッシュ値は変わらないため分割数分しかハッシュは計算しない
	FragmentDest   [][]byte //断片データの送信先ID；メタデータを圧縮するためにIDを用いる
	FragmentOrder  [][]byte //断片データの送信先ごとでの並び順；接頭辞と合わせてファイル名等を表現する
	FragmentPrefix string   //断片データの接頭辞；セッション、ディレクトリ名、ファイル名などを表現する
}
