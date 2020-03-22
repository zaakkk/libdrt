package sesh

//combination は組み合わせの数を計算する
//nCkなので注意
func combination(k uint8, n uint8) uint8 {
	lk := uint64(k)
	ln := uint64(n)
	v := uint64(1)
	for i := uint64(0); i < lk; i++ {
		v *= (ln - i)
	}
	for i := lk; i > 1; i-- {
		v /= i
	}
	return uint8(v)
}
