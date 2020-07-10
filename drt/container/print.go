package container

import "fmt"

//PrintBuf バッファの出力
func PrintBuf(buf []byte) {
	fmt.Printf("\n-----------------printing buf------------------\n")
	for i, v := range buf {
		fmt.Printf("%02X ", v)
		if (i+1)%16 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n------------------------------------------------\n\n")
}
