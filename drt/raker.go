package drt

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"./container"
	"./core"
	"./crypt"
	"./fracc"
	"./metacc"
	"./sesh"
)

//Raker は復号に関わるインターフェースを管理する
type Raker struct {
	OriginHash         crypt.Hash        //元データのハッシュ
	FragmentHash       crypt.Hash        //断片データのハッシュ
	Padding            crypt.Cipher      //パディング
	StreamCipher       crypt.Cipher      //ストリーム暗号
	Scrambler          crypt.Cipher      //一体化
	FragmentDownloader fracc.Downloader  //断片データ送信
	MetadataDownloader metacc.Downloader //メタデータ送信
	MetadataCodec      metacc.Codec      //メタデータのシリアライズ
	MKSParser          sesh.Parser       //メタデータの閾値秘密分散表現
}

//Rake は各種復号に関する関数を呼び出す手続きを簡略化する
//致命的なエラーが有ればエラーハンドルする(予定)
//この関数を気に入らないなら、別に作れ(そのために、他の関数を公開してある)
func (r *Raker) Rake(metadataKey string) (origin *core.Origin, err error) {
	/*defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()*/
	list := r.DownloadMetadata(metadataKey)
	metadata := r.DecodeMetadata(list)
	buf, bufs := r.CreateBuffers(metadata)
	table := r.LoadFragmentTable(bufs, metadata)
	r.FragmentDownloader.Download(table, r.FragmentHash)
	buf = r.Decrypt(buf, metadata)
	origin, err = r.CreateOrigin(buf, metadata)
	return
}

//DownloadMetadata はメタデータをダウンロードしlistに格納する
//メタデータキーを解析したうえで、メタデータのダウンロードをMetadataDownloaderに委譲する
func (r *Raker) DownloadMetadata(metadataKey string) []core.Part {
	var keyPairs map[string]string
	err := json.Unmarshal([]byte(metadataKey), &keyPairs)
	if err != nil {
		fmt.Println(metadataKey)
		panic(err)
	}
	var list []core.Part
	var accessKeys []string
	for key, value := range keyPairs {
		var p core.Part
		p.Dest = key
		list = append(list, p)
		accessKeys = append(accessKeys, value)
	}
	list = r.MetadataDownloader.Download(list, accessKeys)
	return list
}

//DecodeMetadata はメタデータを解析する
//閾値秘密分散処理をされたメタデータを結合をMKSPaserに委譲した上で
//メタデータ解析処理をMetadetaCodecに委譲する
func (r *Raker) DecodeMetadata(list []core.Part) *core.Metadata {
	masterKey := make([][]byte, len(list))
	for i, v := range list {
		masterKey[i] = v.Buffer
	}
	buf := r.MKSParser.Parse(masterKey)
	return r.MetadataCodec.Read(buf)
}

//CreateBuffers はダペンデータ割り当て用のバッファを作成する
//bufsに新しいバッファは与えていけない
func (r *Raker) CreateBuffers(m *core.Metadata) ([]byte, [][]byte) {
	buf := make([]byte, m.FragmentSize*int(m.Division))
	bufs := container.Divide(buf, m.Division)
	bufs = crypt.Shuffle(bufs, m.Order)
	return buf, bufs
}

//LoadFragmentTable はメタデータから断片データのテーブルを作成する
//必要な情報はメタデータから読み出し、バッファはCreateBuffersで作成したバッファを割り当てる
func (r *Raker) LoadFragmentTable(bufs [][]byte, m *core.Metadata) [][]core.Fragment {
	dn, rn := m.Division, m.Replication
	table := make([][]core.Fragment, rn)
	for ri := uint8(0); ri < rn; ri++ {
		table[ri] = make([]core.Fragment, dn)
		for di := uint8(0); di < dn; di++ {
			f := &table[ri][di]
			f.Buffer = bufs[di]
			f.Hash = m.FragmentHash[di]
			f.DestID = m.FragmentDest[ri][di]
			f.Dest = m.Destinations[f.DestID]
			f.Order = m.FragmentOrder[ri][di]
			f.Prefix = m.FragmentPrefix
		}
	}
	return table
}

//Decrypt バッファを復号する
//並び替えはCreateBuffersの時に実施してあるため、この段階では行わない
func (r *Raker) Decrypt(buf []byte, m *core.Metadata) []byte {
	buf = r.Scrambler.Decrypt(buf, m)
	buf = r.StreamCipher.Decrypt(buf, m)
	buf = r.Padding.Decrypt(buf, m)
	return buf
}

//CreateOrigin 元データを復元する
//ハッシュが一致しない場合はpanic()する
func (r *Raker) CreateOrigin(buf []byte, m *core.Metadata) (origin *core.Origin, err error) {
	if !reflect.DeepEqual(m.OriginHash, r.OriginHash.Calc(buf)) {
		err = errors.New("Origin Hash does not match")
		origin = nil
		return
	}
	origin = core.NewOrigin(m.Addition, buf)
	return
}
