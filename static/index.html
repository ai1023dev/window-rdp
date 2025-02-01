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

// 키 변환 함수
func parseKey(key string) string {
	switch key {
	case "A":
		return "a"
	case "B":
		return "b"
	case "C":
		return "c"
	case "D":
		return "d"
	case "E":
		return "e"
	case "F":
		return "f"
	case "G":
		return "g"
	case "H":
		return "h"
	case "I":
		return "i"
	case "J":
		return "j"
	case "K":
		return "k"
	case "L":
		return "l"
	case "M":
		return "m"
	case "N":
		return "n"
	case "O":
		return "o"
	case "P":
		return "p"
	case "Q":
		return "q"
	case "R":
		return "r"
	case "S":
		return "s"
	case "T":
		return "t"
	case "U":
		return "u"
	case "V":
		return "v"
	case "W":
		return "w"
	case "X":
		return "x"
	case "Y":
		return "y"
	case "Z":
		return "z"
	case "ENTER":
		return "enter"
	case "SPACE":
		return "space"
	case "SHIFT":
		return "shift"
	case "CTRL":
		return "ctrl"
	case "ALT":
		return "alt"
	case "TAB":
		return "tab"
	default:
		return ""
	}
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
				robotgo.MouseClick(param1, false)
				log.Printf("마우스 버튼 클릭: %s\n", param1)
			} else {
				log.Println("잘못된 마우스 버튼:", param1)
			}

		case "keydown":
			key := parseKey(param1)
			if key != "" {
				robotgo.KeyToggle(key, "down")
				log.Printf("키 누름: %s\n", key)
			} else {
				log.Println("잘못된 키 입력:", param1)
			}

		case "keyup":
			key := parseKey(param1)
			if key != "" {
				robotgo.KeyToggle(key, "up")
				log.Printf("키 놓음: %s\n", key)
			} else {
				log.Println("잘못된 키 입력:", param1)
			}

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

		document.addEventListener('keydown', function(event) {
			const key = event.key.toUpperCase();
			ws.send('keydown:' + key);
		});

		document.addEventListener('keyup', function(event) {
			const key = event.key.toUpperCase();
			ws.send('keyup:' + key);
		});
	</script>
</body>
</html>`
		fmt.Fprint(w, html)
	})

	log.Println("서버가 8080 포트에서 실행 중...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
