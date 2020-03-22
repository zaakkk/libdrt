drt
===============

# 使い方

詳細はsample.goを参照 まず、これを読む前にsample.goを見ろ

sample.goでは、mapを利用した実装をサンプルとして提供している

```
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

```

# 概要
このライブラリが提供するインターフェースを実装するだけでDRTを扱うことができる。一部機能についてはデフォルトの処理が実装してある。実装されていない機能についても実装が容易になるように機能を提供している。

また、多岐にわたるインターフェースの設定が容易になるようなビルダーメソッドも提供している。

# 構造体
このライブラリで利用する構造体の大部分はdrt または drt/coreで定義されている。

以下は、主な構造体の解説

### drt.Distributer
Distributerは暗号化処理に必要になるインターフェースを保持し、各種暗号化に関する関数を呼び出す手続きを簡略化する。暗号化処理を担うメンバメソッドはditribute()であり、柔軟性のためにこの関数が利用するメンバメソッドは公開されている。Ditribute()を利用して実装できない処理は他のメンバメソッドを利用して実現する必要がある。Distributerが特に要求するインターフェースはfracc.Uploaderとmetacc.Uplaoderであり、**この二つは必ず自分で実装しなければならない。**

### drt.Raker
Rakerは復号処理に必要に必要になるインターフェースを保持し、各種復号に関する関数を呼び出す手続きを簡略化する。復号処理を担うメンバメソッドはRake()であり、柔軟性のためにこの関数が利用するメンバメソッドは公開されている。Rake()を利用して実装できない処理は他のメンバメソッドを利用して実現する必要がある。
Rakerが特に要求するインターフェースはfracc.Downloaderとmetacc.Downlaoderである。**この二つは必ず自分で実装しなければならない。**

### drt.Parameter
Parameterは暗号化処理において必要になる情報を保持する。Distributer.Distribute()の第二引数である。一体化鍵、ストリーム暗号化鍵などは**Ditribute()の前に**用意する必要がある。後述のSetting構造体を用いることで、この手順を簡略化することができる。

### drt.Setting
Settingは分割数、一体化数、ランダマイズ鍵長、断片データ名の接頭辞、断片データの閾値秘密分散表、メタデータの閾値秘密分散表を扱う。**閾値秘密分散表はコンストラクタで与えるべきである**(drt.NewSetting())。メンバは公開されていないが、すべてセッターが用意されている。toParameter()でメンバが設定されたParameterを得ることができる。

### drt.core.Origin
元データに関連する二つの情報、つまり追加データ(ファイル名等)とバッファを管理する。メンバはAddition(追加データ)とBufferである。コンストラクタ(drt.core.NewOrigin)もあるが、第一引数と第二引数を間違わないように注意する必要がある。第二引数がバッファである。

### drt.core.Fragment
断片データに関連する情報を扱う。**FragmentのBufferに新しいスライスを割り当ててはならない**。f.prefx は別の暗号化の時に作成した断片データと名前空間が衝突しないように使う。f.Destは断片データの送信先である。f.Orderは送信先ごとの並び順であり、断片データのファイル名等に使う。

### drt.core.Part(仮称)
閾値秘密分散処理をしたメタデータを扱う

### drt.API
このライブラリが提供するインターフェースを一元管理する。変更したいインターフェースはこの構造体に与えるべきである。このクラスのメンバメソッドにDistributerとRakerのビルダーメソッドがある。基本的にこのメンバメソッドを利用して、


# インターフェース
このライブラリが提供するインテーフェースとは以下の8種類である
- crypt.Hash
- crypt.Cipher
- metacc.Codec
- sesh.Parser
- fracc.Uploader
- fracc.Downloader
- metacc.Uploader
- metacc.Downloader

これらのインターフェースは以下の処理で用いる
- 元データのハッシュ(crypt.Hash)
- 断片データのハッシュ(crypt.Hash)
- パディング(crypt.Cipher)
- ストリーム暗号(crypt.Cipher)
- 一体化(crypt.Cipher)
- メタデータのシリアライズ(metacc.Codec)
- メタデータの閾値秘密分散表現(sesh.Parser)
- 断片データの送信(fracc.Uploader)
- 断片データの受信(fracc.Downloader)
- メタデータの送信(metacc.Uploader)
- メタデータの受信(metacc.Downloader)

これらの実装方法について以下で述べる
### crypto.Hash

```
//Hash はハッシュ関数や誤り検知符号などのインターフェースを定義する
type Hash interface {
	//Calc はbufの中身のハッシュを計算して返す
	Calc(buf []byte) []byte
}
```

引数のbufにハッシュを計算すべきデータが入っている。これに対応するハッシュを計算して戻り値として返す必要がある。  
デフォルトでは、元データ用と断片データともにSHA512が設定される。SHA256も実装されているので必要に応じて置き換えることができる。

### crypto.Cipher

```
//Cipher は暗号のインターフェースを定義する
type Cipher interface {
	//bufを暗号化して返す
	Encrypt(buf []byte, m *core.Metadata) []byte
	//bufを復号して返す
	Decrypt(buf []byte, m *core.Metadata) []byte
}
```

Encrypt、Decryptともに第一引数に暗号化すべきデータが入っている。第二引数のｍに暗号化鍵が設定されているため、これを利用して暗号化・復号をすべきである。決してこの場所で暗号化鍵を作成してはならない。  
パディングのデフォルト値は、分割数と16の倍数になるように0xFFを詰める  
ストリーム暗号のデフォルト値は、RC4である  
一体化はのデフォルト値は、{8, 16, 32, 64}のワード長と{加算, 減算, 排他的論理和}の演算の複合であるが、対象データが8バイトの整数倍でないと処理できない

### metacc.Codec

```
//Codec はメタデータのシリアライズに関するインターフェースを定義する
type Codec interface {
	//Write はメタデータを受け取りシリアライズして返す
	Write(m *core.Metadata) []byte
	//Read はバイト列を受け取り、デシリアライズしてメタデータを構築して返す
	Read(buf []byte) *core.Metadata
}
```

Writeは引数のメタデータをシリアライズして戻り値として返す必要がある。stringに変換できる必要はない。  
Readは引数のバイト列を読み取り、メタデータを構築する必要がある  
デフォルトの実装は、jsonである。(これを別に実装するのは正気じゃない)  

### sesh.Parser

```
//Parser はMKSのシリアライズに関するインターフェースを定義する
type Codec interface {
	//Write はメタデータを受け取りシリアライズして返す
	Write(m *core.Metadata) []byte
	//Read はバイト列を受け取り、デシリアライズしてメタデータを構築して返す
	Read(buf []byte) *core.Metadata
}
```

説明文準備中(これを別に実装するのは正気じゃない)  
[MasterKeySystem(pdf)](https://ipsj.ixsq.nii.ac.jp/ej/?action=repository_action_common_download&item_id=45132&item_no=1&attribute_id=1&file_no=1)を参照  
デフォルトの実装は、jsonである。

### fracc.Uploader

```
//Uploader は断片データ送信のインターフェースを定義する
//実装例はcustomパッケージにある
type Uploader interface {
	//Uploadはtableの中の断片データを送信しなければならない
	Upload(table [][]core.Fragment)
}
```

tableの中の断片データをすべて送信しなけらばならない。すでにシャッフリングされているため送信順序等は無視してよい。  
**デフォルトの実装は用意されていない**  が一から作るのは難しいため、実装を簡略化する構造体を用意してある。詳細はsample.goを参照。

### fracc.Downloader

```
//Downloader は断片データ受信のインターフェースを定義する
//実装例はcustomパッケージにある
type Downloader interface {
	//Downloadはtableの中の断片データを複製された物の中から一つずつ受信しなければならない
	Download(table [][]core.Fragment, hash crypt.Hash)
}

```

tableの中の断片データを複製された物の中から一つずつ受信しなければならない。すでにシャッフリングされているため受信順序等は無視してよい。  
受信した断片データの検証には第二引数のhashを利用すべきである。検証すべきハッシュ値はFragment.Hashに格納されている。  
**デフォルトの実装は用意されていない**  が一から作るのは難しいため、実装を簡略化する構造体を用意してある。詳細はsample.goを参照。

### metacc.Uploader

```
//Uploader はメタデータの送信に関するインターフェースを規定する
//実装例はCustomにある
type Uploader interface {
	//listにあるメタデータを送信し、アクセスするのに必要な情報(アクセスキー)の配列を返さなければならない
	//ただし、アクセスキーはlistの順で返さなければならない
	Upload(list []core.Part) []string
}
```
listにあるメタデータを送信し、アクセスするのに必要な情報(アクセスキー)の配列を返さなければならない。  
ただし、アクセスキーはlistの順で返さなければならない

### metacc.Downloader

```
//Downloader はメタデータをダウンロードするためのインターフェースを定義する
type Downloader interface {
	//Download はaccessKeysの中のアクセスキーを利用し、メタデータをダウンロードする
	//すべてのメタデータをダウンロードする必要はなく、listの順番は関係ない
	//メタデータが足りないかどうかは別の関数が判断するため、この関数が判断しなくて良い
	Download(list []core.Part, accessKeys []string) []core.Part
}
```

メタデータを必要な数だけダウンロードする。必要な数が分からなければすべてダウンロードすればよい。すべてのメタデータをダウンロードする必要はなく、listの順番は関係ない。メタデータが足りないかどうかは別の関数が判断するため、この関数が判断しなくて良い

# Custom
ネットワーク関連の基礎となるインターフェースはdrt.customのサブパッケージとして提供する。今は最低限のminしか用意できていない。

### min

```
//NewFragmentUploader は送信方法を規定した関数を引数に取るコンストラクタ
func NewFragmentUploader(fn func(f *core.Fragment)) *FragmentUploader {}
```

core.FragmentUploaderインターフェースを満たす構造体のコンストラクタ。  
引数のfn関数は、引数の断片データを送信すればよい  
送信とはf.Bufferを外部に送信することである。  
このときの方法は様々であるが、f.prefix string, f.Dest string, f.Order byteなどを利用して送信先を決定すべきである。必ずしもすべてを利用する必要はない。f.prefx は別の暗号化の時に作成した断片データと名前空間が衝突しないように使う。f.Destは断片データの送信先である。f.Orderは送信先ごとの並び順であり、断片データのファイル名等に使う。

```
//NewFragmentDownloader は送信方法を規定した関数を引数に取るコンストラクタ
func NewFragmentDownloader(fn func(f *core.Fragment) bool) *FragmentDownloader {}
```

core.FragmentDownloaderインターフェースを満たす構造体のコンストラクタ。  
引数のfn関数は、引数の断片データを受信する。受信に成功したか否かを返す必要がある。成功したときはtrueである。  
受信とはf.Bufferにダウンロードしたデータをコピーすることである。
すでにf.Bufferには十分な大きさのスライスを割り当ててある。  
間違えてもf.Bufferに新しいスライスを与えてはならない。

```
//NewMetadataUploader は送信方法を規定した関数を引数に取るコンストラクタ
func NewMetadataUploader(fn func(p *core.Part) string) *MetadataUploader {}

```

core.MetadataUploaderインターフェースを満たす構造体のコンストラクタ。
引数のfn関数は、引数のメタデータを送信すればよい。このとき、メタデータにアクセスするための情報(例えば、ファイル名)を作成し返却する必要がある。


```
//NewMetadataDownloader は送信方法を規定した関数を引数に取るコンストラクタ
func NewMetadataDownloader(fn func(p *core.Part, accessKeys string)) *MetadataDownloader {}
```

引数のfn関数は、引数のaccessKeyとp.Destに基づき、メタデータを受信して内容をp.Bufferに与える必要がある。ただし、p.Bufferにnilであるため、スライスは自分で用意する必要がある。

# Contribution

### Custom
ネットワーク関連の基本構造体はcustomのサブパッケージとして提供する予定だが、まだ最低限のminしか用意できていない。そのため、他の方式の実装を他の人にお願いしたい。

### 用語

MKS 閾値秘密分散の方式  
アクセスキー 閾値秘密分散したメタデータにアクセスするための文字列  
メタデータキー アクセスキーと送信先情報を束ねたJSON文字列  

### 命名規則

fracc = FRAgment ACCess  
metacc = METadata ACCess 
sesh = SEcure SHering 

table = 断片データの表  
list = Partの表

m = Metadataの変数
buf = byte[]の変数 
Buffer = メンバとしてのbuf

ループインデックスをi,j,kと続けることは禁止する。例えば、ri,diなど続けるべきである。
