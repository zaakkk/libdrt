package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"./drt"
	"./drt/core"
	"./drt/crypt"
	"./drt/custom/min"
)

//断片データ保存用マップ
var fragmentStorage map[string][]byte = map[string][]byte{}

//メタデータ保存用マップ
var metadataStorage map[string][]byte = map[string][]byte{}

//断片データ送信関数
func storeFragment(f *core.Fragment) {
	filename := f.Dest + f.Prefix + strconv.Itoa(int(f.Order))
	fragmentStorage[filename] = make([]byte, len(f.Buffer))
	copy(fragmentStorage[filename], f.Buffer)
}

//メタデータ送信関数
func storeMetadata(p *core.Part) string {
	accessKey := strconv.Itoa(int(crypt.CreateRandomByte(255)))
	filename := p.Dest + accessKey
	fragmentStorage[filename] = make([]byte, len(p.Buffer))
	copy(fragmentStorage[filename], p.Buffer)
	return accessKey
}

//断片データ受信関数
func readFragment(f *core.Fragment) bool {
	filename := f.Dest + f.Prefix + strconv.Itoa(int(f.Order))
	//!!!!!f.Bufferは書き換えるな!!!!!!
	if _, ok := fragmentStorage[filename]; ok {
		copy(f.Buffer, fragmentStorage[filename])
		return true
	}
	return false
}

//メタデータ受信関数
func readMetadata(p *core.Part, accessKey string) {
	filename := p.Dest + accessKey
	//fmt.Printf("%s \n", filename)
	if _, ok := fragmentStorage[filename]; ok {
		p.Buffer = make([]byte, len(fragmentStorage[filename]))
		copy(p.Buffer, fragmentStorage[filename])
	}
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
	fragmentDest := []string{"1", "2", "3"}
	metadataDest := []string{"A", "B", "C"}
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
