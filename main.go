package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
)

// WebSocket 업그레이더 설정
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 화면 캡처 함수
func captureScreen() ([]byte, error) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 웹소켓 핸들러
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 업그레이드 오류:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket 연결 성공")

	// 현재 모니터 해상도 가져오기
	screenBounds := screenshot.GetDisplayBounds(0)
	screenWidth := screenBounds.Dx()
	screenHeight := screenBounds.Dy()

	// 스크린 전송 루틴
	go func() {
		for {
			imgData, err := captureScreen()
			if err != nil {
				log.Println("스크린 캡처 오류:", err)
				break
			}

			err = conn.WriteMessage(websocket.BinaryMessage, imgData)
			if err != nil {
				log.Println("이미지 전송 오류:", err)
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 메시지 수신 루프
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("메시지 수신 오류:", err)
			break
		}

		// ":" 기준으로 메시지 분리
		parts := strings.Split(string(message), ":")
		if len(parts) < 2 {
			log.Println("잘못된 메시지 형식:", string(message))
			continue
		}

		action := parts[0]
		param1 := parts[1]
		param2 := ""
		if len(parts) > 2 {
			param2 = parts[2]
		}

		switch action {
		case "mousemove":
			// 비율 값(0~1)을 실제 좌표로 변환
			xRatio, err1 := strconv.ParseFloat(param1, 64)
			yRatio, err2 := strconv.ParseFloat(param2, 64)
			if err1 != nil || err2 != nil {
				log.Println("좌표 변환 오류:", param1, param2)
				continue
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
			scrollAmount, err := strconv.Atoi(param1)
			if err != nil {
				log.Println("스크롤 값 변환 오류:", param1)
				continue
			}
			robotgo.ScrollSmooth(0, scrollAmount)
			log.Printf("스크롤 이동: %d\n", scrollAmount)

		case "keydown":
			key := strings.ToLower(param1) // 대문자를 자동으로 소문자로 변환
			robotgo.KeyToggle(key, "down")
			log.Printf("키 누름: %s\n", key)

		case "keyup":
			key := strings.ToLower(param1)
			robotgo.KeyToggle(key, "up")
			log.Printf("키 놓음: %s\n", key)

		default:
			log.Printf("알 수 없는 액션: %s\n", action)
		}
	}
}

// 서버 실행
func main() {
	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
	<title>원격 제어 웹</title>
</head>
<body>
	<h1>실시간 화면 공유 및 원격 제어</h1>
	<canvas id="screen" width="1280" height="720"></canvas>
	<script>
		const canvas = document.getElementById('screen');
		const ctx = canvas.getContext('2d');
		const ws = new WebSocket('ws://' + window.location.host + '/ws');

		ws.binaryType = 'arraybuffer';

		ws.onmessage = function(event) {
			const img = new Image();
			const blob = new Blob([event.data]);
			img.src = URL.createObjectURL(blob);
			img.onload = function() {
				ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
			};
		};

		canvas.addEventListener('mousemove', function(event) {
			const rect = canvas.getBoundingClientRect();
			const xRatio = (event.clientX - rect.left) / canvas.width;
			const yRatio = (event.clientY - rect.top) / canvas.height;
			ws.send('mousemove:' + xRatio.toFixed(5) + ':' + yRatio.toFixed(5));
		});

		canvas.addEventListener('mousedown', function(event) {
			const button = event.button === 0 ? 'left' : event.button === 2 ? 'right' : '';
			if (button) ws.send('mousedown:' + button);
		});

		canvas.addEventListener('mouseup', function(event) {
			const button = event.button === 0 ? 'left' : event.button === 2 ? 'right' : '';
			if (button) ws.send('mouseup:' + button);
		});

		canvas.addEventListener('wheel', function(event) {
			const scrollAmount = event.deltaY > 0 ? -5 : 5;
			ws.send('scroll:' + scrollAmount);
		});

		document.addEventListener('keydown', function(event) {
			ws.send('keydown:' + event.key.toLowerCase());
		});

		document.addEventListener('keyup', function(event) {
			ws.send('keyup:' + event.key.toLowerCase());
		});
	</script>
</body>
</html>`
		fmt.Fprint(w, html)
	})

	log.Println("서버가 8080 포트에서 실행 중...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
