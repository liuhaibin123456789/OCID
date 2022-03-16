package tool

import (
	"math/rand"
	"time"
)

//RandNum 生成随机数字
func RandNum() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(100000) + 10000
}
