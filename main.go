package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
)

// HTTP 요청 핸들러
func handleEvent(w http.ResponseWriter, r *http.Request) {
	// 쿼리 파라미터에서 action, param1, param2 추출
	action := r.URL.Query().Get("action")
	param1 := r.URL.Query().Get("param1")
	param2 := r.URL.Query().Get("param2")

	// 현재 모니터 해상도 가져오기
	screenBounds := screenshot.GetDisplayBounds(0)
	screenWidth := screenBounds.Dx()
	screenHeight := screenBounds.Dy()

	switch action {
	case "mousemove":
		// 비율 값(0~1)을 실제 좌표로 변환
		xRatio, err1 := strconv.ParseFloat(param1, 64)
		yRatio, err2 := strconv.ParseFloat(param2, 64)
		if err1 != nil || err2 != nil {
			log.Println("좌표 변환 오류:", param1, param2)
			return
		}

		x := int(xRatio * float64(screenWidth))
		y := int(yRatio * float64(screenHeight))

		robotgo.Move(x, y)
		log.Printf("마우스 이동: %d, %d\n", x, y)

	case "mousedown":
		if param1 == "left" || param1 == "right" {
			robotgo.MouseDown(param1) // 마우스 버튼 누름
			log.Printf("마우스 버튼 누름: %s\n", param1)
		} else {
			log.Println("잘못된 마우스 버튼:", param1)
		}

	case "mouseup":
		if param1 == "left" || param1 == "right" {
			robotgo.MouseUp(param1) // 마우스 버튼 놓음
			log.Printf("마우스 버튼 놓음: %s\n", param1)
		} else {
			log.Println("잘못된 마우스 버튼:", param1)
		}

	case "scroll":
		if param1 == "up" {
			robotgo.Scroll(0, -1)
			log.Printf("스크롤 위")
		} else if param1 == "down" {
			robotgo.Scroll(0, 1)
			log.Printf("스크롤 아래")
		} else {
			log.Println("잘못된 스크롤 방향:", param1)
		}

	case "keydown":
		key := strings.ToLower(param1) // 대문자를 자동으로 소문자로 변환
		switch key {
		case "arrowup":
			key = "up"
		case "arrowdown":
			key = "down"
		case "arrowleft":
			key = "left"
		case "arrowright":
			key = "right"
		}
		robotgo.KeyToggle(key, "down")
		log.Printf("키 누름: %s\n", key)

	case "keyup":
		key := strings.ToLower(param1)
		switch key {
		case "arrowup":
			key = "up"
		case "arrowdown":
			key = "down"
		case "arrowleft":
			key = "left"
		case "arrowright":
			key = "right"
		}
		robotgo.KeyToggle(key, "up")
		log.Printf("키 놓음: %s\n", key)

	default:
		log.Printf("알 수 없는 액션: %s\n", action)
	}
}

// 서버 실행
func main() {
	http.HandleFunc("/event", handleEvent)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>원격 제어 웹</title>
		</head>
		<body>
			<script>
				function sendEvent(action, param1, param2) {
					let url = '/event?action=' + action + '&param1=' + param1;
					if (param2 !== undefined) {
						url += '&param2=' + param2;
					}
					fetch(url).catch(error => console.error('Error:', error));
				}

				document.addEventListener('mousemove', function(event) {
					const xRatio = event.clientX / window.innerWidth;
					const yRatio = event.clientY / window.innerHeight;
					sendEvent('mousemove', xRatio.toFixed(5), yRatio.toFixed(5));
				});

				document.addEventListener('mousedown', function(event) {
					const button = event.button === 0 ? 'left' : event.button === 2 ? 'right' : '';
					if (button) sendEvent('mousedown', button);
					event.preventDefault();
				});

				document.addEventListener('mouseup', function(event) {
					const button = event.button === 0 ? 'left' : event.button === 2 ? 'right' : '';
					if (button) sendEvent('mouseup', button);
					event.preventDefault();
				});

				document.addEventListener('wheel', function(event) {
					const direction = event.deltaY > 0 ? 'up' : 'down';
					sendEvent('scroll', direction);
					event.preventDefault();
				});

				document.addEventListener('keydown', function(event) {
					sendEvent('keydown', event.key.toLowerCase());
					event.preventDefault();
				});

				document.addEventListener('keyup', function(event) {
					sendEvent('keyup', event.key.toLowerCase());
					event.preventDefault();
				});
			</script>
		</body>
		</html>
		`
		fmt.Fprint(w, html)
	})

	log.Println("서버가 8080 포트에서 실행 중...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
