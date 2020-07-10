package sesh

//Parser はMKSのシリアライズに関するインターフェースを定義する
type Parser interface {
	//Describe はbufで表現されるバイト列をsysに基づき
	//復元に必要な情報を付加し、閾値秘密分散処理をして返す
	Describe(buf []byte, sys *MasterKeySys) [][]byte
	//Parse はchildkeyを利用してMasterKeyを復元し返す
	//エラーハンドリングは行数が多くなるため必ずしも実装する必要はない
	Parse(childKey [][]byte) []byte
}
