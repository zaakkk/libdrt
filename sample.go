package main

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"./drt"
	"./drt/core"
	"./drt/crypt"
	"./drt/custom/min"
)

type mail struct {
	from     string
	username string
	password string
	to       string
	sub      string
	msg      string
}

func (m mail) body() string {
	return "To: " + m.to + "\r\n" +
		"Subject: " + m.sub + "\r\n\r\n" +
		m.msg + "\r\n"
}

func yahooMailSend(m mail) error {
	smtpSvr := "smtp.mail.yahoo.co.jp:587"
	auth := smtp.PlainAuth("", m.username, m.password, "smtp.mail.yahoo.co.jp")
	if err := smtp.SendMail(smtpSvr, auth, m.from, []string{m.to}, []byte(m.body())); err != nil {
		return err
	}
	return nil
}

//断片データ送信関数
func storeFragment(f *core.Fragment) {

	//メール設定
	m := mail{
		from:     "*********@yahoo.co.jp",
		username: "*********@yahoo.co.jp",
		password: "*********",
		to:       f.Dest,
		sub:      f.Prefix,
		msg:      string(f.Buffer),
	}
	if err := yahooMailSend(m); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	filename := f.Dest + f.Prefix + strconv.Itoa(int(f.Order))
	err := ioutil.WriteFile(f.Dest+"/"+filename, f.Buffer, 0666)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second * 10)
}

//メタデータ送信関数
func storeMetadata(p *core.Part) string {
	accessKey := strconv.Itoa(int(crypt.CreateRandomByte(255)))

	//メール設定
	m := mail{
		from:     "*********@yahoo.co.jp",
		username: "*********@yahoo.co.jp",
		password: "*********",
		to:       p.Dest,
		sub:      p.Dest + accessKey,
		msg:      string(p.Buffer),
	}
	if err := yahooMailSend(m); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	filename := p.Dest + accessKey
	err := ioutil.WriteFile(p.Dest+"/"+filename, p.Buffer, 0666)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second * 10)

	return accessKey
}

//断片データ受信関数
func readFragment(f *core.Fragment) bool {
	filename := f.Dest + f.Prefix + strconv.Itoa(int(f.Order))
	bytes, err := ioutil.ReadFile(f.Dest + "/" + filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	copy(f.Buffer, bytes)
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
	filename := p.Dest + accessKey
	//fmt.Printf("%s \n", filename)
	bytes, err := ioutil.ReadFile(p.Dest + "/" + filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	p.Buffer = make([]byte, len(bytes))
	copy(p.Buffer, bytes)
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
	fu := min.NewFragmentUploader(storeFragment)
	mu := min.NewMetadataUploader(storeMetadata)
	d := api.BuildDistributer(fu, mu)

	//複合を担う構造体を扱うときは、断片データ受信インターフェースとメタデータ受信インターフェースが
	//デフォルト値が用意されていないため自分で用意する必要がある
	//ただ、自分で1から作ることは困難なである、
	//しかし、簡単にインターフェースを満たすためにmin.New~を使ってインターフェースを満たす構造体を用意することができる
	fd := min.NewFragmentDownloader(readFragment)
	md := min.NewMetadataDownloader(readMetadata)
	r := api.BuildRaker(fd, md)

	return d, r
}

//ファイルの読み込み
func readFile() *core.Origin {
	filename := "dummy.txt"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		panic("failed to open dummy file")
	}

	//drt.core.Originは元データの情報(ファイル名とバッファ)を管理する
	return core.NewOrigin([]byte(filename), bytes)
}

func main() {
	//暗号化と復号化を担う構造体の作成
	//d: drt/Distributer
	//r: drt/Raker
	d, r := setup()

	//暗号化パラメータの作成
	//NewSettingでSetting(drt/Setting)を作成する
	//Settingは暗号化パラメータを保持する
	//Settingが扱うパラメータは分割数、一体化数、ランダマイズ鍵長、断片データ名の接頭辞、断片データの閾値秘密分散表、メタデータの閾値秘密分散表
	//Settingには最低限のパラメータが最初から設定されているため、セッターで特別設定する必要はない
	//Setting.ToParameter()で安全に暗号化鍵を生成する
	//param(drt/Parameter)
	//宛先メールアドレス
	fragmentDest := []string{"17aj148@ms.dendai.ac.jp", "drt0000000@gmail.com", "taisei.y_is_here@au.com"}
	metadataDest := []string{"17aj148@ms.dendai.ac.jp", "drt0000000@gmail.com", "taisei.y_is_here@au.com"}

	param := drt.NewSetting(fragmentDest, 2, metadataDest, 2).SetDivision(4).SetPrefix(12).SetScramble(1).ToParameter()

	//ファイルの読み込み
	origin := readFile()

	//暗号化
	//key(string)にはメタデータにアクセスするための情報が記述されている
	//何らかの処理に失敗した場合はerrに何かが入っている
	key, err := d.Distribute(origin, param)
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}
	fmt.Printf("\nsecceed to distribute : key is %s\n\n", key)

	//復号
	//復号結果がrecoveredに入っている
	//何らかの処理に失敗した場合はerrに何かが入っている
	recovered, err := r.Rake(key)
	if err != nil {
		fmt.Printf("error: %#v\n", err)
		panic(err)
	}
	fmt.Println("secceed to rake. content is ...")
	fmt.Println(string(recovered.Buffer))

	fmt.Println()
}
