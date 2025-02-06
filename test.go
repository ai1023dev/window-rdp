package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const SCREEN_CAPTURE = "desktop" // 전체 화면을 캡처

var ffmpegArgs = []string{
	"-f", "gdigrab", // Windows에서 화면 캡처를 위한 gdigrab 사용
	"-i", SCREEN_CAPTURE, // 화면 캡처 입력
	"-c:v", "libx264", // 비디오 코덱 = libx264
	"-crf", "23", // 비디오 품질 (23이 기본)
	"-pix_fmt", "yuv420p", // 호환성 있는 픽셀 포맷
	"-hls_time", "2", // HLS 세그먼트 길이
	"-hls_list_size", "5", // 최대 리스트 크기 = 5
	"-hls_delete_threshold", "1", // 세그먼트 최대 개수
	"-hls_flags", "delete_segments", // 오래된 세그먼트 삭제
	"-f", "hls", // 출력 포맷 = HLS
	"public/video.m3u8", // HLS 출력 파일을 public 폴더에 저장
}

func main() {
	// 기존 HLS 파일 삭제: *.ts, *.m3u8 파일 찾기
	fmt.Println("Initialize video")

	// *.ts 파일 삭제
	tsFiles, err := filepath.Glob("public/*.ts")
	if err != nil {
		log.Fatal("Error finding .ts files: ", err)
	}
	for _, file := range tsFiles {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Error removing .ts file %s: %v\n", file, err)
		} else {
			fmt.Printf("Removed file: %s\n", file)
		}
	}

	// *.m3u8 파일 삭제
	m3u8Files, err := filepath.Glob("public/*.m3u8")
	if err != nil {
		log.Fatal("Error finding .m3u8 files: ", err)
	}
	for _, file := range m3u8Files {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Error removing .m3u8 file %s: %v\n", file, err)
		} else {
			fmt.Printf("Removed file: %s\n", file)
		}
	}

	// 스트리밍 시작
	fmt.Println("Start streaming")
	ffmpegCmd := exec.Command("ffmpeg", ffmpegArgs...)
	ffmpegCmd.Stdout = os.Stdout
	ffmpegCmd.Stderr = os.Stderr

	err = ffmpegCmd.Start()
	if err != nil {
		log.Fatal("Error starting ffmpeg: ", err)
	}
	fmt.Println("Running ffmpeg command: ", strings.Join(ffmpegArgs, " "))

	// 정적 파일 서빙
	fs := http.FileServer(http.Dir("public")) // 'public' 폴더를 서빙
	http.Handle("/public/", http.StripPrefix("/public", fs))

	// HTTP 서버 시작
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html") // index.html 파일 서빙
	})

	// 서버 시작
	fmt.Println("Server started on port 3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
