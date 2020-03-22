package sesh

//TableSize は(k,n)閾値秘密分散のテーブルサイズを計算する
func TableSize(k uint8, n uint8) uint8 {
	return combination(n-k+1, n)
}

//MasterKeySys は閾値秘密分散の割り当てテーブルを管理する
type MasterKeySys struct {
	k        uint8    //閾値
	n        uint8    //サーバー数
	row      [][]bool //生のテーブル
	size     uint8    //テーブルの大きさ
	property uint8    //各々の送信先が持つファイルの数
}

//Threshold は閾値を返す
func (t *MasterKeySys) Threshold() uint8 {
	return t.k
}

//CoOwners は鍵の所有者の数を返す
func (t *MasterKeySys) CoOwners() uint8 {
	return t.n
}

//Replication は複製数を返す
func (t *MasterKeySys) Replication() uint8 {
	return t.n - t.k + 1
}

//Size はテーブルの大きさ(鍵の分割数)を返す
func (t *MasterKeySys) Size() uint8 {
	return t.size
}

//Property は鍵の所有者が保有する鍵の数を返す
func (t *MasterKeySys) Property() uint8 {
	if t.property == 0 {
		t.property = TableSize(t.k-1, t.n-1)
	}
	return t.property
}

//At は対応する場所に鍵が存在するか否かを判定する
func (t *MasterKeySys) At(di uint8, ci uint8) bool {
	return t.row[di][ci]
}

//NewMKS は安全にMKSを作成する
func NewMKS(k uint8, n uint8) *MasterKeySys {
	if k > n {
		panic("master key system creation error; k > n")
	}
	if n > 10 {
		panic("master key system creation error; n > 10")
	}
	if k > 10 {
		panic("master key system creation error; k > 10")
	}
	t := new(MasterKeySys)
	t.k = k
	t.n = n
	t.size = TableSize(t.k, t.n)
	t.row = make([][]bool, t.size)
	for i := uint8(0); i < t.n; i++ {
		t.row[i] = make([]bool, t.n)
	}
	t.createMasterKeySys(t.k, t.n, 0)
	return t
}

//-----特別読む必要はない------
//読むのであれば論文を読んでから
func (t *MasterKeySys) createMasterKeySys(k uint8, n uint8, index uint8) {
	if k == 1 {
		for i := uint8(0); i < n; i++ {
			t.row[index][i] = true
		}
	} else if k == n {
		for i := uint8(0); i < n; i++ {
			t.row[index+i][i] = true
		}
	} else {
		d0 := TableSize(k-1, n-1)
		d1 := TableSize(k, n-1)
		t.createMasterKeySys(k-1, n-1, index)
		t.createMasterKeySys(k, n-1, index+d0)
		for i := uint8(0); i < d1; i++ {
			t.row[index+d0+i][n-1] = true
		}
	}
}
