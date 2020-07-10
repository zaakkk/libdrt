package basic

import (
	"encoding/json"

	"github.com/zaakkk/libdrt/drt/core"
)

//MetadataJSONCodec はメタデータをjsonを使ってシリアライズする
type MetadataJSONCodec struct {
}

//Write はjson形式でメタデータを表現する
func (p *MetadataJSONCodec) Write(m *core.Metadata) []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes
}

//Read はjson形式のメタデータを読み取る
func (p *MetadataJSONCodec) Read(buf []byte) *core.Metadata {
	//fmt.Println(string(buf))
	m := new(core.Metadata)
	err := json.Unmarshal(buf, m)
	if err != nil {
		panic(err)
	}
	return m
}
