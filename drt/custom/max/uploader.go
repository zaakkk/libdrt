package max

import (
	"../../core"
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

//uploadAndClose minとの整合性を取りながらchanをクローズして終了を宣言する
func (u *FragmentUploader) uploadAndClose(f *core.Fragment, done chan struct{}) {
	u.upload(f)
	close(done)
}

//Upload は最低限の断片データの送信手続きを規定する
func (u *FragmentUploader) Upload(table [][]core.Fragment) {
	rn, dn := len(table), len(table[0])
	done := make([][]chan struct{}, rn)
	for ri := 0; ri < rn; ri++ {
		done[ri] = make([]chan struct{}, dn)
	}
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			f := &table[ri][di]
			d := make(chan struct{})
			go u.uploadAndClose(f, d)
			done[ri][di] = d
		}
	}
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			<-done[ri][di]
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

//uploadAndClose minとの整合性を取りながらchanをクローズして終了を宣言する
func (u *MetadataUploader) uploadAndClose(p *core.Part, accessKey chan string) {
	accessKey <- u.upload(p)
	close(accessKey)
}

//Upload は最低限のメタデータの送信手続きを規定する
func (u *MetadataUploader) Upload(list []core.Part) []string {
	size := len(list)
	accessKeyList := make([]string, size)
	channels := make([]chan string, size)
	for i := range list {
		c := make(chan string)
		go u.uploadAndClose(&list[i], c)
		channels[i] = c
	}
	for i := range list {
		accessKeyList[i] = <-channels[i]
	}
	return accessKeyList
}
