package util

import "math/rand"

func RandDecimal(m, d int) string {
	ms := randNum(m - d)
	ds := randNum(d)
	var i int
	for i = range ms {
		if ms[i] != byte('0') {
			break
		}
	}
	ms = ms[i:]
	l := len(ms) + len(ds) + 1
	flag := rand.Intn(2)
	//check for 0.0... avoid -0.0
	zeroFlag := true
	for i := range ms {
		if ms[i] != byte('0') {
			zeroFlag = false
		}
	}
	for i := range ds {
		if ds[i] != byte('0') {
			zeroFlag = false
		}
	}
	if zeroFlag {
		flag = 0
	}
	vs := make([]byte, 0, l+flag)
	if flag == 1 {
		vs = append(vs, '-')
	}
	vs = append(vs, ms...)
	if len(ds) == 0 {
		return string(vs)
	}
	vs = append(vs, '.')
	vs = append(vs, ds...)
	return string(vs)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

const numberBytes = "0123456789"

func randNum(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Int63()%int64(len(numberBytes))]
	}
	return b
}
