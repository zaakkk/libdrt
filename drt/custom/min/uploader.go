package min

import (
	"github.com/zaakkk/libdrt/drt/core"
)

//FragmentUploader は最低限の断片データの送信手続きを規定する
//具体的な送信手続きはuploadに委譲する
type FragmentUploader struct {
	upload func(f *core.Fragment)
}

//NewFragmentUploader は送信方法を規定した関数を引数に取るコンストラクタ
func NewFragmentUploader(fn func(f *core.Fragment)) *FragmentUploader {
	u := new(FragmentUploader)
	u.upload = fn
	return u
}

//Upload は最低限の断片データの送信手続きを規定する
func (u *FragmentUploader) Upload(table [][]core.Fragment) {
	rn, dn := len(table), len(table[0])
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			f := &table[ri][di]
			u.upload(f)

		}
	}
}

//MetadataUploader は最低限のメタデータの送信手続きを規定する
//具体的な送信手続きはuploadに委譲する
type MetadataUploader struct {
	upload func(p *core.Part) string
}

//NewMetadataUploader は送信方法を規定した関数を引数に取るコンストラクタ
func NewMetadataUploader(fn func(p *core.Part) string) *MetadataUploader {
	u := new(MetadataUploader)
	u.upload = fn
	return u
}

//Upload は最低限のメタデータの送信手続きを規定する
func (u *MetadataUploader) Upload(list []core.Part) []string {
	accessKeyList := make([]string, len(list))
	for i, v := range list {
		accessKey := u.upload(&v)
		accessKeyList[i] = accessKey

	}
	return accessKeyList
}
