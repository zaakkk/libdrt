package sesh

//MasterKeySysHandler はMasterKeySysを扱いやすくするためのクラス
type MasterKeySysHandler struct {
	dest       []string      //送信先
	assignment *MasterKeySys //マスターキー本体
}

//NewMKSHandler はマスターキーを安全な方法で生成する
func NewMKSHandler(dest []string, threshold uint8) *MasterKeySysHandler {
	h := new(MasterKeySysHandler)
	h.dest = dest
	h.assignment = NewMKS(threshold, uint8(len(dest)))
	return h
}

//GetDest は送信先のゲッター
func (h *MasterKeySysHandler) GetDest() []string {
	return h.dest
}

//GetAssignment は割り当て(MKS)のゲッター
func (h *MasterKeySysHandler) GetAssignment() *MasterKeySys {
	return h.assignment
}

//AssignDest は断片データをMKSに対応する送信先を返す
func (h *MasterKeySysHandler) AssignDest(divisionIndex uint8, replicationIndex uint8) uint8 {
	rn := h.assignment.Replication()
	if replicationIndex > rn {
		panic("MSK is encounted illegal argment; replicaitonIndex")
	}
	count := uint8(0)
	di := divisionIndex % h.assignment.size
	for i := uint8(0); i < h.assignment.n; i++ {
		if h.assignment.row[di][i] {
			if count == replicationIndex {
				return i
			}
			count++
		}
	}
	panic("MSK handler is broken")
}
