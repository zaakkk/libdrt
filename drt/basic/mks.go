package basic

import (
	"encoding/json"
	"fmt"

	"github.com/zaakkk/libdrt/drt/container"
	"github.com/zaakkk/libdrt/drt/sesh"
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
	size := sys.Size()
	divided := container.Divide(buf, size)
	for i := uint8(0); i < size; i++ {
		if sys.At(i, srvIndex) {
			c.Keys[int(i)] = divided[i]
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
	size := sesh.TableSize(children[0].Threshold, children[0].Server)
	divided := make([][]byte, size)
	for di := uint8(0); di < size; di++ {
		failed := true
		for _, child := range children {
			childKey, ok := child.Keys[int(di)]
			if ok {
				divided[di] = childKey
				failed = false
				break
			}
		}
		if failed {
			panic("insufficiency of child key")
		}
	}
	return container.Bundle(divided)
}
