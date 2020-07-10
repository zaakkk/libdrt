package min

import (
	"fmt"
	"reflect"

	"github.com/zaakkk/libdrt/drt/core"
	"github.com/zaakkk/libdrt/drt/crypt"
)

//FragmentDownloader は最低限の断片データの送信手続きを規定する
//具体的な送信手続きはdownloadに委譲する
type FragmentDownloader struct {
	download func(*core.Fragment) bool
}

//NewFragmentDownloader は送信方法を規定した関数を引数に取るコンストラクタ
func NewFragmentDownloader(fn func(f *core.Fragment) bool) *FragmentDownloader {
	d := new(FragmentDownloader)
	d.download = fn
	return d
}

//Download は最低限の断片データの送信手続きを規定する
func (d *FragmentDownloader) Download(table [][]core.Fragment, hash crypt.Hash) {
	rn, dn := len(table), len(table[0])
	for di := 0; di < dn; di++ {
		failed := true
		for ri := 0; ri < rn; ri++ {
			f := &table[ri][di]
			if d.download(f) {
				if reflect.DeepEqual(f.Hash, hash.Calc(f.Buffer)) {
					failed = false
					break
				}
				fmt.Println("hash error")
			}
		}
		if failed {
			panic("failed to collect enough fragment")
		}
	}
}

//MetadataDownloader は最低限のメタデータの送信手続きを規定する
//具体的な送信手続きはdownloadに委譲する
type MetadataDownloader struct {
	download func(p *core.Part, accessKey string)
}

//NewMetadataDownloader は送信方法を規定した関数を引数に取るコンストラクタ
func NewMetadataDownloader(fn func(p *core.Part, accessKeys string)) *MetadataDownloader {
	d := new(MetadataDownloader)
	d.download = fn
	return d
}

//Download は最低限のメタデータの送信手続きを規定する
func (d *MetadataDownloader) Download(list []core.Part, accessKeys []string) []core.Part {
	for i := range list {
		d.download(&list[i], accessKeys[i])
	}
	return list
}
