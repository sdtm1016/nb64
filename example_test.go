package nb64

import (
	"fmt"
	"log"
)

// deceive go-vet
func Singles() {}
func Doubles() {}
func Threes()  {}
func Fours()   {}
func Mixing()  {}

func ExampleSingles() {
	src := "We are all good kids."
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 21 25 r5UGHllQYdmxAn378iDXpyc1w
}

func ExampleThrees() {
	src := "我们都是好孩子。"
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 24 26 .GIRE7sJD9GYvFl9FtpFtQDAC.
}

func ExampleMixing_a() {
	src := "We are all good kids.我们都是好孩子。"
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 45 51 r5UGHllQYdmxAn378iDXpyc1w.GIRE7sJD9GYvFl9FtpFtQDAC.
}

func ExampleMixing_b() {
	src := "We 我们 are 是 all 都 good 好 kids 孩子.。"
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 50 68 r5UA.GIRE7s.QYeWVA.GYv.QYdmxA.JD9.QZ9-_Jg.Fl9.Qa9OTng.FtpFtQ.XA.DAC.
}

func ExampleDoubles() {
	src := "¢"
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 2 5 .ACi.
}

func ExampleFours() {
	src := "𐄁"
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 4 5 .QEB.
}

func ExampleMixing_c() {
	src := "We 我𐄁们 are ¢ 是 all 都 ¢ good 好 kids 孩𐄁子.。"
	dst, err := URLEncoding.Encode([]byte(src))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(src), len(dst), string(dst))
	// Output: 64 88 r5UA.GIRQEBE7s.QYeWVA.ACi.QA.GYv.QYdmxA.JD9.QA.ACi.QZ9-_Jg.Fl9.Qa9OTng.FtpQEBFtQ.XA.DAC.
}
