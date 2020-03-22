package basic

import (
	"encoding/json"
	"fmt"

	"../sesh"
)

//JSONMSK はjson形式MSKを扱う
type JSONMSK struct {
}

//child は子鍵を表現する
type child struct {
	Threshold     uint8          //閾値
	Server        uint8          //サーバー数
	MasterKeySize int            //マスターキーの大きさ
	Property      uint8          //サーバーに割り当てられた鍵の大きさ
	Keys          map[int][]byte //鍵
}

//newChil 子鍵を作成する
func newChild(buf []byte, sys *sesh.MasterKeySys, srvIndex uint8) child {
	var c child
	c.Threshold = sys.Threshold()
	c.Server = sys.CoOwners()
	c.MasterKeySize = len(buf)
	c.Property = sys.Property()
	c.Keys = make(map[int][]byte)
	size := c.MasterKeySize / int(sys.Size())
	for i := uint8(0); i < sys.Size(); i++ {
		if !sys.At(i, srvIndex) {
			continue
		}
		begin := size * int(i)
		end := begin + size
		if i != c.Server-1 {
			c.Keys[int(i)] = buf[begin:end]
		} else {
			c.Keys[int(i)] = buf[begin:len(buf)]
		}
	}
	return c
}

//Describe はbufで表現されるバイト列をsysに基づき
//復元に必要な情報を付加し、閾値秘密分散処理をして返す
func (p *JSONMSK) Describe(buf []byte, sys *sesh.MasterKeySys) [][]byte {
	children := make([]child, sys.CoOwners())
	for i := range children {
		children[i] = newChild(buf, sys, uint8(i))
	}
	childKey := make([][]byte, sys.CoOwners())
	for i, v := range children {
		str, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			panic("JSONMSK is broken")
		}
		childKey[i] = str
	}
	return childKey
}

//Parse はchildkeyを利用してMasterKeyを復元し返す
func (p *JSONMSK) Parse(childKey [][]byte) []byte {
	children := make([]child, len(childKey))
	for i, v := range childKey {
		err := json.Unmarshal(v, &children[i])
		if err != nil {
			fmt.Println(err)
		}
	}

	buf := make([]byte, children[0].MasterKeySize)
	size := sesh.TableSize(children[0].Threshold, children[0].Server)
	keySize := children[0].MasterKeySize / int(size)
	for di := uint8(0); di < size; di++ {
		failed := true
		for _, child := range children {
			childKey, ok := child.Keys[int(di)]
			if ok {
				begin := keySize * int(di)
				end := begin + keySize
				if di != size-1 {
					copy(buf[begin:end], childKey)
				} else {
					copy(buf[begin:len(buf)], childKey)
				}
				failed = false
				break
			}
		}
		if failed {
			panic("insufficiency of child key")
		}
	}
	return buf
}
