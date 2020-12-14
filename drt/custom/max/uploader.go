package max

import (
	"log"
	"strconv"
	"sync"
	"syscall/js"

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

//uploadAndClose minとの整合性を取りながらchanをクローズして終了を宣言する
func (u *FragmentUploader) uploadAndClose(f *core.Fragment, done chan struct{}) {
	u.upload(f)
	close(done)
}

/*
//Upload は最低限の断片データの送信手続きを規定する
//default
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

*/

/*
//Upload は最低限の断片データの送信手続きを規定する
//並行数変更不可能版
func (u *FragmentUploader) Upload(table [][]core.Fragment) {
	rn, dn := len(table), len(table[0])
	done := make([][]chan struct{}, rn)

	var wg sync.WaitGroup

	for ri := 0; ri < rn; ri++ {
		done[ri] = make([]chan struct{}, dn)
	}
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			wg.Add(1)

			f := &table[ri][di]
			d := make(chan struct{})
			go func(fragment *core.Fragment, done chan struct{}) {
				defer wg.Done()
				u.uploadAndClose(fragment, done)
			}(f, d)

			done[ri][di] = d
		}
	}
	wg.Wait()
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			<-done[ri][di]
		}
	}
	document := js.Global().Get("document")
	sendInterval := document.Call("getElementById", "sendInterval").Get("value").String()
	si, err := time.ParseDuration(sendInterval)
	if err != nil {
		log.Println(err)
	}

	startSleep := time.Now()
	time.Sleep(si)
	endSleep := time.Now()
	fmt.Printf("SleepF: %f\n", (endSleep.Sub(startSleep)).Seconds())
}
*/

///*
// FUploadConcurrency は最大同時並列実行数
//const FUploadConcurrency = 10

//Upload は最低限の断片データの送信手続きを規定する
//並行数変更可能版
func (u *FragmentUploader) Upload(table [][]core.Fragment) {
	rn, dn := len(table), len(table[0])
	done := make([][]chan struct{}, rn)

	var wg sync.WaitGroup

	document := js.Global().Get("document")
	FUploadConcurrency := document.Call("getElementById", "FUploadConcurrency").Get("value").String()
	Fi, err := strconv.ParseUint(FUploadConcurrency, 10, 8)
	if err != nil {
		log.Println(err)
	}

	sem := make(chan struct{}, Fi)

	for ri := 0; ri < rn; ri++ {
		done[ri] = make([]chan struct{}, dn)
	}
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			f := &table[ri][di]
			d := make(chan struct{})

			sem <- struct{}{}

			wg.Add(1)
			go func(fragment *core.Fragment, done chan struct{}) {
				defer wg.Done()
				defer func() { <-sem }()
				u.uploadAndClose(fragment, done)
			}(f, d)

			done[ri][di] = d
		}
		wg.Wait()
	}
	for ri := 0; ri < rn; ri++ {
		for di := 0; di < dn; di++ {
			<-done[ri][di]
		}
	}
}

//*/

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
//default
///*
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

//*/

/*
//Upload は最低限のメタデータの送信手続きを規定する
//並行数制限なし
func (u *MetadataUploader) Upload(list []core.Part) []string {
	size := len(list)
	accessKeyList := make([]string, size)
	channels := make([]chan string, size)

	var wg sync.WaitGroup
	wg.Add(size)
	for i := range list {
		c := make(chan string)
		go func(p *core.Part, accessKey chan string) {
			defer wg.Done()
			u.uploadAndClose(p, accessKey)

			//待機時間
			//time.Sleep(1 * time.Second)
		}(&list[i], c)
		channels[i] = c
	}
	for i := range list {
		accessKeyList[i] = <-channels[i]
	}

	wg.Wait()
	document := js.Global().Get("document")
	sendInterval := document.Call("getElementById", "sendInterval").Get("value").String()
	si, err := time.ParseDuration(sendInterval)
	if err != nil {
		log.Println(err)
	}

	startSleep := time.Now()
	time.Sleep(si)
	endSleep := time.Now()
	fmt.Printf("SleepM: %f\n", (endSleep.Sub(startSleep)).Seconds())
	return accessKeyList
}

*/

/*
// MUploadConcurrency は最大同時並列実行数
const MUploadConcurrency = 3

//Upload は最低限のメタデータの送信手続きを規定する
//並列実行数変更可能版
func (u *MetadataUploader) Upload(list []core.Part) []string {
	size := len(list)
	accessKeyList := make([]string, size)
	channels := make([]chan string, size)

	var wg sync.WaitGroup
	sem := make(chan struct{}, MUploadConcurrency) //concurrency数のバッファ

	for i := range list {
		sem <- struct{}{}

		wg.Add(1)

		c := make(chan string)
		channels[i] = c

		go func(p *core.Part, accessKey chan string) {
			defer wg.Done()
			defer func() { <-sem }()
			u.uploadAndClose(p, accessKey)
		}(&list[i], c)

	}

	wg.Wait()

	for i := range list {
		accessKeyList[i] = <-channels[i]
	}
	return accessKeyList
}

*/
