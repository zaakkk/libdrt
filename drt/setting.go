package drt

import (
	"encoding/base64"

	"github.com/zaakkk/libdrt/drt/crypt"
	"github.com/zaakkk/libdrt/drt/sesh"
)

//Setting は暗号強度や送信先情報を保持する
//ToPrameter()を利用してParameterに変換することができる
type Setting struct {
	division        uint8
	scramble        uint8
	randomize       uint8
	prefix          uint8
	fragmentHandler *sesh.MasterKeySysHandler
	metadataHandler *sesh.MasterKeySysHandler
}

//NewSetting 新しくSettingを生成した上で、最低限の暗号化強度を設定する
func NewSetting(destf []string, tf uint8, destm []string, tm uint8) *Setting {
	s := new(Setting)
	s.division = 60
	s.scramble = 6
	s.randomize = 32
	s.prefix = 36
	s.fragmentHandler = sesh.NewMKSHandler(destf, tf)
	s.metadataHandler = sesh.NewMKSHandler(destm, tm)
	return s
}

//SetDivision は分割数のセッター
func (s *Setting) SetDivision(division uint8) *Setting {
	s.division = division
	return s
}

//SetScramble は一体化数のセッター
func (s *Setting) SetScramble(scramble uint8) *Setting {
	s.scramble = scramble
	return s
}

//SetRandomize はストリーム暗号鍵長のセッター
func (s *Setting) SetRandomize(randomize uint8) *Setting {
	s.randomize = randomize
	return s
}

//SetPrefix は断片データの接頭辞長のセッター
func (s *Setting) SetPrefix(prefix uint8) *Setting {
	s.prefix = prefix
	return s
}

//ToParameter はSettingをParameterに変換する
//この際、安全な方法で各種暗号化鍵を生成すし設定する
func (s *Setting) ToParameter() *Parameter {
	p := new(Parameter)
	p.Division = s.division
	p.Scramble = crypt.CreateRandomBytes(int(s.scramble))
	p.Randomize = crypt.CreateRandomBytes(int(s.randomize))
	p.Prefix = base64.URLEncoding.EncodeToString(crypt.CreateRandomBytes(int(s.prefix)))
	p.Order = crypt.RandomOrder(s.division)
	p.FragmentHandler = s.fragmentHandler
	p.MetadataHandler = s.metadataHandler
	return p
}
