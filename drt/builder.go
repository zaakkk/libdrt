package drt

import (
	"github.com/zaakkk/libdrt/drt/basic"
	"github.com/zaakkk/libdrt/drt/crypt"
	"github.com/zaakkk/libdrt/drt/fracc"
	"github.com/zaakkk/libdrt/drt/metacc"
	"github.com/zaakkk/libdrt/drt/sesh"
)

//API は各種インターフェースを管理する
type API struct {
	OriginHash    crypt.Hash   //元データのハッシュ
	FragmentHash  crypt.Hash   //断片データのハッシュ
	Padding       crypt.Cipher //パディング
	StreamCipher  crypt.Cipher //ストリーム暗号
	Scrambler     crypt.Cipher //一体化
	MetadataCodec metacc.Codec //メタデータのシリアライズ
	MKSParser     sesh.Parser  //メタデータの閾値秘密分散表現
}

//SetDefault デフォルトを設定する 設定されている場合は変更しない
func (i *API) SetDefault() *API {
	if i.OriginHash == nil {
		i.OriginHash = new(basic.SHA512)
	}

	if i.FragmentHash == nil {
		i.FragmentHash = new(basic.SHA512)
	}

	if i.Padding == nil {
		i.Padding = new(basic.MinPadding)
	}

	if i.StreamCipher == nil {
		i.StreamCipher = new(basic.RC4)
	}

	if i.Scrambler == nil {
		i.Scrambler = new(basic.Scrambler)
	}

	if i.MetadataCodec == nil {
		i.MetadataCodec = new(basic.MetadataJSONCodec)
	}

	if i.MKSParser == nil {
		i.MKSParser = new(basic.JSONMSK)
	}

	return i
}

//BuildDistributer はDistributerの作成手続きを簡略化する
func (i *API) BuildDistributer(fu fracc.Uploader, mu metacc.Uploader) *Distributer {
	d := new(Distributer)
	d.FragmentUploader = fu
	d.MetadataUploader = mu

	d.OriginHash = i.OriginHash
	d.FragmentHash = i.FragmentHash
	d.Padding = i.Padding
	d.StreamCipher = i.StreamCipher
	d.Scrambler = i.Scrambler
	d.MetadataCodec = i.MetadataCodec
	d.MKSParser = i.MKSParser
	return d
}

//BuildRaker はRakerの作成手続きを簡略化する
func (i *API) BuildRaker(fd fracc.Downloader, md metacc.Downloader) *Raker {
	r := new(Raker)
	r.FragmentDownloader = fd
	r.MetadataDownloader = md

	r.OriginHash = i.OriginHash
	r.FragmentHash = i.FragmentHash
	r.Padding = i.Padding
	r.StreamCipher = i.StreamCipher
	r.Scrambler = i.Scrambler
	r.MetadataCodec = i.MetadataCodec
	r.MKSParser = i.MKSParser
	return r
}
