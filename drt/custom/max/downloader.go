package max

import (
	"reflect"

	"../../core"
	"../../crypt"
)

//FragmentDownloader は最低限の断片データの送信手続きを規定する
//具体的な送信手続きはdownloadに委譲する
type FragmentDownloader struct {
	download func(f *core.Fragment) bool
}

//NewFragmentDownloader は送信方法を規定した関数を引数に取るコンストラクタ
func NewFragmentDownloader(fn func(f *core.Fragment) bool) *FragmentDownloader {
	d := new(FragmentDownloader)
	d.download = fn
	return d
}

func (d *FragmentDownloader) downloadAndClose(table [][]core.Fragment, success chan bool, di int, hash crypt.Hash) {
	defer close(success)
	rn := len(table)
	for ri := 0; ri < rn; ri++ {
		f := &table[ri][di]
		result := d.download(f)
		if !result {
			continue
		}
		if reflect.DeepEqual(f.Hash, hash.Calc(f.Buffer)) {
			success <- true
			return
		}
	}
	success <- false
}

//Download は最低限の断片データの送信手続きを規定する
func (d *FragmentDownloader) Download(table [][]core.Fragment, hash crypt.Hash) {
	dn := len(table[0])
	success := make([]chan bool, dn)

	for di := 0; di < dn; di++ {
		c := make(chan bool)
		go d.downloadAndClose(table, c, di, hash)
		success[di] = c
	}
	for di := 0; di < dn; di++ {
		if !<-success[di] {
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

func (d *MetadataDownloader) downloadAndClose(p *core.Part, accessKeys string, done chan struct{}) {
	d.download(p, accessKeys)
	close(done)
}

//Download は最低限のメタデータの送信手続きを規定する
func (d *MetadataDownloader) Download(list []core.Part, accessKeys []string) []core.Part {
	done := make([]chan struct{}, len(list))
	for i := range list {
		c := make(chan struct{})
		go d.downloadAndClose(&list[i], accessKeys[i], c)
		done[i] = c
	}
	for i := range list {
		<-done[i]
	}
	return list
}
