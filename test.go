package main

import (
	"time"

	"github.com/go-vgo/robotgo"
)

// 외부에서 호출 가능하도록 Test 함수의 첫 글자를 대문자로 변경
func main() {
	time.Sleep(5 * time.Second)
	robotgo.Scroll(0, 1)
}
