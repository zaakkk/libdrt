package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/zaakkk/libdrt/drt"
	"github.com/zaakkk/libdrt/drt/core"
	"github.com/zaakkk/libdrt/drt/crypt"
	"github.com/zaakkk/libdrt/drt/custom/max"

	"github.com/zaakkk/libdrt/drtMail/coreMail"
	"github.com/zaakkk/libdrt/drtMail/recieve"
	"github.com/zaakkk/libdrt/drtMail/send"
)

var keyGlobal string

//断片データ送信関数
//送信元:宛先
func storeFragment(f *core.Fragment) {
	//startSendF := time.Now()

	addressAndPass := strings.Split(f.Dest, "::")
	from := addressAndPass[0]
	fromPass := addressAndPass[1]
	to := addressAndPass[2]

	//メール設定

	m := coreMail.MailStruct{
		From:     from,
		Username: from,
		Password: fromPass,
		To:       to,
		Sub:      to + f.Prefix + strconv.Itoa(int(f.Order)),
		Msg:      f.Buffer,
	}

	if err := send.GMailSend(m); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//APIサーバに直接送る
	//if err := send.GMailSendToAPI(m); err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	/*
		document := js.Global().Get("document")
		sendInterval := document.Call("getElementById", "sendInterval").Get("value").String()
		si, err := time.ParseDuration(sendInterval)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(si)
	*/

	//endSendF := time.Now()
	//fmt.Printf("storeF: %f\n", (endSendF.Sub(startSendF)).Seconds())
	//fmt.Printf("%f\n", (endSendF.Sub(startSendF)).Seconds())

}

//メタデータ送信関数
func storeMetadata(p *core.Part) string {
	//startSendM := time.Now()

	accessKey := strconv.Itoa(int(crypt.CreateRandomByte(255)))
	addressAndPass := strings.Split(p.Dest, "::")
	from := addressAndPass[0]
	fromPass := addressAndPass[1]
	to := addressAndPass[2]

	//メール設定
	m := coreMail.MailStruct{
		From:     from,
		Username: from,
		Password: fromPass,
		To:       to,
		Sub:      to + accessKey,
		Msg:      p.Buffer,
	}
	//webサーバ経由(おそらく)
	if err := send.GMailSend(m); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//APIサーバに直接送る(使用不可)
	//if err := send.GMailSendToAPI(m); err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}

	/*
		document := js.Global().Get("document")
		sendInterval := document.Call("getElementById", "sendInterval").Get("value").String()
		si, err := time.ParseDuration(sendInterval)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(si)
	*/

	//endSendM := time.Now()
	//fmt.Printf("storeM: %f\n", (endSendM.Sub(startSendM)).Seconds())
	//fmt.Printf("%f\n", (endSendM.Sub(startSendM)).Seconds())

	return accessKey
}

//断片データ受信関数
func readFragment(f *core.Fragment) bool {
	//startRecieveF := time.Now()

	addressAndPass := strings.Split(f.Dest, "::")
	to := addressAndPass[2]
	toPass := addressAndPass[3]

	sub := to + f.Prefix + strconv.Itoa(int(f.Order))
	data, err := recieve.GMailRecieve(to, toPass, sub)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//例外処理必要
	copy(f.Buffer, data)

	//endRecieveF := time.Now()
	//fmt.Printf("readF: %f\n", (endRecieveF.Sub(startRecieveF)).Seconds())
	//fmt.Printf("%f\n", (endRecieveF.Sub(startRecieveF)).Seconds())

	return true

	//!!!!!f.Bufferは書き換えるな!!!!!!
	//if _, ok := storage[filename]; ok {
	//	copy(f.Buffer, storage[filename])
	//	return true
	//}

	//return false
}

//メタデータ受信関数
func readMetadata(p *core.Part, accessKey string) {

	//startRecieveM := time.Now()

	addressAndPass := strings.Split(p.Dest, "::")
	to := addressAndPass[2]
	toPass := addressAndPass[3]

	sub := to + accessKey
	data, err := recieve.GMailRecieve(to, toPass, sub)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	p.Buffer = make([]byte, len(data))
	copy(p.Buffer, data)

	//endRecieveM := time.Now()
	//fmt.Printf("readM: %f\n", (endRecieveM.Sub(startRecieveM)).Seconds())
	//fmt.Printf("%f\n", (endRecieveM.Sub(startRecieveM)).Seconds())

}

//DistributerとRakerの準備
func setup() (*drt.Distributer, *drt.Raker) {
	var api drt.API
	//api.OriginHash = new(myCheackSum) //apiを変更したい場合はここで行う
	api.SetDefault()

	//暗号化を担う構造体を扱うときは、断片データ送信インターフェースとメタデータ送信インターフェースが
	//デフォルト値が用意されていないため自分で用意する必要がある
	//ただ、自分で1から作ることは困難なである、
	//しかし、簡単にインターフェースを満たすためにmin.New~を使ってインターフェースを満たす構造体を用意することができる
	//fu := min.NewFragmentUploader(storeFragment)
	//mu := min.NewMetadataUploader(storeMetadata)
	fu := max.NewFragmentUploader(storeFragment)
	mu := max.NewMetadataUploader(storeMetadata)
	d := api.BuildDistributer(fu, mu)

	//複合を担う構造体を扱うときは、断片データ受信インターフェースとメタデータ受信インターフェースが
	//デフォルト値が用意されていないため自分で用意する必要がある
	//ただ、自分で1から作ることは困難なである、
	//しかし、簡単にインターフェースを満たすためにmin.New~を使ってインターフェースを満たす構造体を用意することができる
	//fd := min.NewFragmentDownloader(readFragment)
	//md := min.NewMetadataDownloader(readMetadata)
	fd := max.NewFragmentDownloader(readFragment)
	md := max.NewMetadataDownloader(readMetadata)
	r := api.BuildRaker(fd, md)

	return d, r
}

//テキストの読み込み
func readText(text string) *core.Origin {

	filename := "dummy"
	//drt.core.Originは元データの情報(ファイル名とバッファ)を管理する
	return core.NewOrigin([]byte(filename), []byte(text))
}

func distribute(d *drt.Distributer) {
	//送信元メール::送信元パスワード::宛先アドレス::宛先パスワード
	fragmentDest := []string{
		"example1@yahoo.co.jp::password1::example2@yahoo.co.jp::password2",
		"example3@yahoo.co.jp::password3::example4@yahoo.co.jp::password4",
		"example5@yahoo.co.jp::password5::example6@yahoo.co.jp::password6",
	}

	metadataDest := []string{
		"example1@yahoo.co.jp::password1::example2@yahoo.co.jp::password2",
		"example3@yahoo.co.jp::password3::example4@yahoo.co.jp::password4",
		"example5@yahoo.co.jp::password5::example6@yahoo.co.jp::password6",
	}

	//テキストボックスに入力した文章を取り出す
	document := js.Global().Get("document")
	text := document.Call("getElementById", "text").Get("value").String()
	divisionNumber := document.Call("getElementById", "divisionNumber").Get("value").String()
	scrambleNumber := document.Call("getElementById", "scrambleNumber").Get("value").String()
	//fmt.Println("text: " + text)

	//fmt.Println("dNum: " + divisionNumber)
	//fmt.Println("sNum: " + scrambleNumber)
	dn, err := strconv.ParseUint(divisionNumber, 10, 8)
	if err != nil {
		log.Println(err)
	}
	sn, err := strconv.ParseUint(scrambleNumber, 10, 8)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(dn + sn)

	param := drt.NewSetting(fragmentDest, 2, metadataDest, 2).SetDivision(uint8(dn)).SetPrefix(12).SetScramble(uint8(sn)).ToParameter()
	//param := drt.NewSetting(fragmentDest, 2, metadataDest, 2).SetDivision(4).SetPrefix(12).SetScramble(1).ToParameter()

	origin := readText(text)

	//暗号化
	//key(string)にはメタデータにアクセスするための情報が記述されている
	//何らかの処理に失敗した場合はerrに何かが入っている
	startDistribute := time.Now()
	key, err := d.Distribute(origin, param)
	endDistribute := time.Now()
	keyGlobal = key
	fmt.Printf("Distribute: %f\n", (endDistribute.Sub(startDistribute)).Seconds())
	//fmt.Printf("%f\n", (endDistribute.Sub(startDistribute)).Seconds())
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}
	//fmt.Printf("\nsecceed to distribute : key is %s\n\n", key)
}

func rake(r *drt.Raker) {
	//復号
	//復号結果がrecoveredに入っている
	//何らかの処理に失敗した場合はerrに何かが入っている
	startRake := time.Now()
	//recovered, err := r.Rake(key)
	_, err := r.Rake(keyGlobal)
	endRake := time.Now()
	fmt.Printf("Rake: %f\n", (endRake.Sub(startRake)).Seconds())
	//fmt.Printf("%f\n", (endRake.Sub(startRake)).Seconds())

	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}
	//fmt.Println("secceed to rake. content is ...")
	//fmt.Println(string(recovered.Buffer))
}

func main() {
	//暗号化と復号化を担う構造体の作成
	//d: drt/Distributer
	//r: drt/Raker
	startAll := time.Now()

	//d, r := setup()
	d, _ := setup()

	//暗号化パラメータの作成
	//NewSettingでSetting(drt/Setting)を作成する
	//Settingは暗号化パラメータを保持する
	//Settingが扱うパラメータは分割数、一体化数、ランダマイズ鍵長、断片データ名の接頭辞、断片データの閾値秘密分散表、メタデータの閾値秘密分散表
	//Settingには最低限のパラメータが最初から設定されているため、セッターで特別設定する必要はない
	//Setting.ToParameter()で安全に暗号化鍵を生成する
	//param(drt/Parameter)
	distribute(d)
	//time.Sleep(time.Second * 30)
	//rake(r)
	//time.Sleep(time.Second * 30)

	fmt.Println()
	endAll := time.Now()
	fmt.Printf("All: %f\n", (endAll.Sub(startAll)).Seconds())
	//fmt.Printf("%f\n", (endAll.Sub(startAll)).Seconds())
}
