package drt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zaakkk/libdrt/drt/container"
	"github.com/zaakkk/libdrt/drt/core"
	"github.com/zaakkk/libdrt/drt/crypt"
	"github.com/zaakkk/libdrt/drt/fracc"
	"github.com/zaakkk/libdrt/drt/metacc"
	"github.com/zaakkk/libdrt/drt/sesh"
)

//Parameter は暗号化パラメータを表す
//デバッグ以外の目的で直接生成すべきではない
//Setting.ToParameter()を介して生成すべきである
type Parameter struct {
	Division        uint8
	Scramble        []byte
	Randomize       []byte
	Prefix          string
	Order           []byte
	FragmentHandler *sesh.MasterKeySysHandler
	MetadataHandler *sesh.MasterKeySysHandler
}

//Distributer は暗号化に関わるインターフェースを管理する
type Distributer struct {
	OriginHash       crypt.Hash      //元データのハッシュ
	FragmentHash     crypt.Hash      //断片データのハッシュ
	Padding          crypt.Cipher    //パディング
	StreamCipher     crypt.Cipher    //ストリーム暗号
	Scrambler        crypt.Cipher    //一体化
	FragmentUploader fracc.Uploader  //断片データ送信
	MetadataUploader metacc.Uploader //メタデータ送信
	MetadataCodec    metacc.Codec    //メタデータのシリアライズ
	MKSParser        sesh.Parser     //メタデータの閾値秘密分散表現
}

//Distribute は各種暗号化に関する関数を呼び出す手続きを簡略化する
//致命的なエラーが有ればエラーハンドルする(予定)
//この関数を気に入らないなら、別に作れ(そのために、他の関数を公開してある)
func (d *Distributer) Distribute(orgn *core.Origin, prm *Parameter) (metadataKey string, err error) {
	/*defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()*/
	var metadata core.Metadata
	d.StoreParams(prm.Division, prm.Scramble, prm.Randomize, prm.Prefix, prm.Order, prm.FragmentHandler, &metadata)
	d.StoreOriginInfo(orgn.Buffer, orgn.Addition, &metadata)

	startEncrypt := time.Now()
	bufs := d.Encrypt(orgn.Buffer, &metadata)
	endEncrypt := time.Now()
	fmt.Printf("Encrypt(DRT): %f\n", (endEncrypt.Sub(startEncrypt)).Seconds())

	table := d.CreateFragmentTable(bufs, prm.FragmentHandler, &metadata)
	d.StoreFragmentTable(table, &metadata)
	d.FragmentUploader.Upload(table)
	list := d.CreatePartList(prm.MetadataHandler, &metadata)
	metadataKey = d.UploadMetadata(list)
	return
}

//StoreParams は暗号化に関わる情報を保存
func (d *Distributer) StoreParams(dn uint8, sk []byte, rk []byte, prfx string, order []byte, hndr *sesh.MasterKeySysHandler, m *core.Metadata) {
	m.Division = dn
	m.Replication = hndr.GetAssignment().Replication()
	m.Scramble = sk
	m.Randomize = rk
	m.Order = order
	m.Destinations = hndr.GetDest()
	m.FragmentPrefix = prfx
}

//StoreOriginInfo は元データに関わることを保存
func (d *Distributer) StoreOriginInfo(buf []byte, addition []byte, m *core.Metadata) {
	m.Addition = addition
	m.DataSize = len(buf)
	m.OriginHash = d.OriginHash.Calc(buf)
}

//Encrypt はbufの暗号化を行いbufを返す
func (d *Distributer) Encrypt(buf []byte, m *core.Metadata) [][]byte {
	buf = d.Padding.Encrypt(buf, m)
	buf = d.StreamCipher.Encrypt(buf, m)
	buf = d.Scrambler.Encrypt(buf, m)
	bufs := container.Divide(buf, m.Division)
	return crypt.Shuffle(bufs, m.Order)
}

//CreateFragmentTable はFragmentの多次元配列(複製数x分割数 !!複製数が先!!)を作成する
//Fragmentにコンストラクタを設けて、この関数の可読性をあげるべき
func (d *Distributer) CreateFragmentTable(bufs [][]byte, hndr *sesh.MasterKeySysHandler, m *core.Metadata) [][]core.Fragment {
	if len(bufs)%int(m.Division) != 0 {
		panic("unable to divide buffer")
	}
	table := make([][]core.Fragment, m.Replication)
	var dn, rn = m.Division, m.Replication
	for ri := uint8(0); ri < rn; ri++ {
		table[ri] = make([]core.Fragment, m.Division)
		for di := uint8(0); di < dn; di++ {
			f := &table[ri][di]
			f.Buffer = bufs[di]
			f.Hash = d.FragmentHash.Calc(f.Buffer)
			f.DestID = hndr.AssignDest(di, ri)
			f.Dest = hndr.GetDest()[f.DestID]
			f.Prefix = m.FragmentPrefix
		}
	}
	d.SpecifyOrder(table, m)
	return table
}

//SpecifyOrder は各断片データの送信先ごとでの位置を特定する
func (d *Distributer) SpecifyOrder(table [][]core.Fragment, m *core.Metadata) {
	var dn, rn = m.Division, m.Replication
	orders := make([]uint8, len(m.Destinations))
	for ri := uint8(0); ri < rn; ri++ {
		for di := uint8(0); di < dn; di++ {
			f := &table[ri][di]
			order := orders[f.DestID]
			f.Order = order
			orders[f.DestID]++
		}
	}
}

//StoreFragmentTable は断片データの情報のうち、必要なものをメタデータに保管する
func (d *Distributer) StoreFragmentTable(table [][]core.Fragment, m *core.Metadata) {
	var dn, rn = m.Division, m.Replication
	m.FragmentHash = make([][]byte, dn)
	for di := uint8(0); di < dn; di++ {
		m.FragmentHash[di] = table[0][di].Hash
	}
	m.FragmentDest = make([][]byte, rn)
	m.FragmentOrder = make([][]byte, rn)
	for ri := uint8(0); ri < rn; ri++ {
		m.FragmentDest[ri] = make([]byte, dn)
		m.FragmentOrder[ri] = make([]byte, dn)
		for di := uint8(0); di < dn; di++ {
			m.FragmentDest[ri][di] = table[ri][di].DestID
			m.FragmentOrder[ri][di] = table[ri][di].Order
		}
	}
	m.FragmentSize = len(table[0][0].Buffer)
}

//CreatePartList はメタデータをシリアライズした上で、閾値秘密分散処理をする
func (d *Distributer) CreatePartList(hndr *sesh.MasterKeySysHandler, m *core.Metadata) []core.Part {
	binary := d.MetadataCodec.Write(m)
	childKeys := d.MKSParser.Describe(binary, hndr.GetAssignment())
	table := make([]core.Part, len(hndr.GetDest()))
	for i, v := range childKeys {
		table[i].Buffer = v
		table[i].Dest = hndr.GetDest()[i]
	}
	return table
}

//UploadMetadata はMetadataUploaderにuploadを委譲したうえで
//upload()の戻り値であるアクセスキー(メタデータの保存場所などを表現した文字列)を
//json形式にフォーマットして返す
func (d *Distributer) UploadMetadata(table []core.Part) string {
	accessKeys := d.MetadataUploader.Upload(table)
	metadataKey := make(map[string]string)
	for i := range table {
		metadataKey[table[i].Dest] = accessKeys[i]
	}
	bytes, err := json.Marshal(metadataKey)
	if err != nil {
		panic("jsonParser is broken")
	}
	return string(bytes)
}
