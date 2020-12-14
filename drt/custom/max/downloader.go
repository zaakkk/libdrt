package max

import (
	"reflect"

	"github.com/zaakkk/libdrt/drt/core"
	"github.com/zaakkk/libdrt/drt/crypt"
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

///*
//Download は最低限の断片データの送信手続きを規定する
//default
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

//*/

/*
//Download は最低限の断片データの送信手続きを規定する
//並列実行数制御不可能版
//TODO: WaitGroupを使った処理が出来ない
func (d *FragmentDownloader) Download(table [][]core.Fragment, hash crypt.Hash) {
	dn := len(table[0])
	success := make([]chan bool, dn)

	var wg sync.WaitGroup
	log.Println("A:%d", dn)
	wg.Add(dn)
	for di := 0; di < dn; di++ {

		c := make(chan bool)
		go func(table2 [][]core.Fragment, success2 chan bool, di2 int, hash2 crypt.Hash) {
			defer wg.Done()
			d.downloadAndClose(table2, success2, di2, hash2)
			log.Println("WOW")
		}(table, c, di, hash)
		success[di] = c
	}
	log.Println("B")
	wg.Wait()
	log.Println("C")
	for di := 0; di < dn; di++ {
		if !<-success[di] {
			panic("failed to collect enough fragment")
		}
	}
	log.Println("D")
}
*/
/*
// FDownloadConcurrency は最大同時並列実行数
const FDownloadConcurrency = 10 //最大同時並列実行数

//Download は最低限の断片データの送信手続きを規定する
//並列実行数変化可能版
func (d *FragmentDownloader) Download(table [][]core.Fragment, hash crypt.Hash) {
	dn := len(table[0])
	success := make([]chan bool, dn)

	var wg sync.WaitGroup
	sem := make(chan struct{}, FDownloadConcurrency) // concurrency数のバッファ

	for di := 0; di < dn; di++ {
		sem <- struct{}{}
		wg.Add(1)

		c := make(chan bool)
		go func(table2 [][]core.Fragment, success2 chan bool, di2 int, hash2 crypt.Hash) {
			defer wg.Done()
			defer func() { <-sem }() //処理が終わったらチャネル解放
			d.downloadAndClose(table2, success2, di2, hash2)
		}(table, c, di, hash)
		success[di] = c
	}
	for di := 0; di < dn; di++ {
		if !<-success[di] {
			panic("failed to collect enough fragment")
		}
	}
}

*/

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

///*

//Download は最低限のメタデータの送信手続きを規定する
//default
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

//*/

/*
//Download は最低限のメタデータの送信手続きを規定する
//並列実行数変更不可能版
func (d *MetadataDownloader) Download(list []core.Part, accessKeys []string) []core.Part {
	size := len(list)
	done := make([]chan struct{}, size)

	var wg sync.WaitGroup
	log.Println("E")
	wg.Add(size)
	log.Println("F")
	for i := range list {
		c := make(chan struct{})
		go func(p *core.Part, accessKey string, done chan struct{}) {
			defer wg.Done()
			d.downloadAndClose(p, accessKey, done)
		}(&list[i], accessKeys[i], c)
		done[i] = c
	}
	log.Println("G")
	wg.Wait()
	log.Println("H")
	for i := range list {
		<-done[i]
	}
	log.Println("I")
	return list
}
*/

/*
// MDownloadConcurrency は最大同時並列実行数
const MDownloadConcurrency = 3

//Download は最低限のメタデータの送信手続きを規定する
//並列実行数変更可能版
func (d *MetadataDownloader) Download(list []core.Part, accessKeys []string) []core.Part {
	done := make([]chan struct{}, len(list))

	var wg sync.WaitGroup
	sem := make(chan struct{}, MDownloadConcurrency) // concurrency数のバッファ

	for i := range list {
		c := make(chan struct{})
		wg.Add(1)
		go func(p *core.Part, accessKey string, done chan struct{}) {
			defer wg.Done()
			defer func() { <-sem }() //処理が終わったらチャネル解放
			d.downloadAndClose(p, accessKey, done)
		}(&list[i], accessKeys[i], c)
		done[i] = c
	}
	for i := range list {
		<-done[i]
	}
	return list
}

*/
