package util

import (
	"fmt"
	"math/rand"
	"time"
)

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

func RdBool() bool {
	return rand.Intn(2) == 1
}

// RdRange rand int in range
func RdRange(n, m int64) int64 {
	if n == m {
		return n
	}
	if m < n {
		n, m = m, n
	}
	return n + rand.Int63n(m-n)
}

// RdInt64 rand int64
func RdInt64() int64 {
	if rand.Intn(2) == 1 {
		return rand.Int63()
	}
	return -rand.Int63() - 1
}

// RdDate rand date
func RdDate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 1, 0, time.UTC).Unix()
	max := time.Date(2100, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// RdString rand string with given length. TODO: support rand multi-byte utf8
func RdString(length int) string {
	res := ""
	for i := 0; i < length; i++ {
		charCode := RdRange(33, 127)
		// char '\' and '"' should be escaped
		if charCode == 92 || charCode == 34 {
			charCode++
			// res = fmt.Sprintf("%s%s", res, "\\")
		}
		res = fmt.Sprintf("%s%s", res, string(rune(charCode)))
	}
	return res
}

// RdTimestamp return same format as RdDate except rand range
// TIMESTAMP has a range of '1970-01-01 00:00:01' UTC to '2038-01-19 03:14:07'
func RdTimestamp() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
