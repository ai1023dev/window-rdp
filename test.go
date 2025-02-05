package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	mux := http.NewServeMux()
	mux.HandleFunc("/offer", handleOffer)

	http.ListenAndServe(":8080", mux)
}

func handleOffer(w http.ResponseWriter, r *http.Request) {
	// WebRTC offer 받기
	var offer webrtc.SessionDescription
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode offer: %v", err), http.StatusBadRequest)
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating peer connection: %v", err), http.StatusInternalServerError)
		return
	}

	// 트랙 생성 (NewTrackLocalStaticRTP 사용)
	track, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating track: %v", err), http.StatusInternalServerError)
		return
	}
	_, err = peerConnection.AddTrack(track)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding track: %v", err), http.StatusInternalServerError)
		return
	}

	// 오디오 캡처 시작
	go captureAudio(track)

	// Offer를 PeerConnection에 설정
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		http.Error(w, fmt.Sprintf("Error setting remote description: %v", err), http.StatusInternalServerError)
		return
	}

	// Answer 생성
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating answer: %v", err), http.StatusInternalServerError)
		return
	}

	// 생성된 Answer를 클라이언트로 반환
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		http.Error(w, fmt.Sprintf("Error setting local description: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding answer: %v", err), http.StatusInternalServerError)
	}
}

func captureAudio(track *webrtc.TrackLocalStaticRTP) {
	buffer := make([]int16, 1024) // 오디오 데이터를 저장할 버퍼
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, 1024, &buffer)
	if err != nil {
		fmt.Println("Error opening audio stream:", err)
		return
	}
	defer stream.Close()

	for {
		err := stream.Read()
		if err != nil {
			fmt.Println("Error reading audio stream:", err)
			return
		}

		// RTP 패킷 생성 후 전송
		pkt := &rtp.Packet{
			Header:  rtp.Header{Timestamp: uint32(time.Now().UnixNano()), PayloadType: 111},
			Payload: int16ToBytes(buffer), // 변환 후 Payload에 저장
		}

		if err := track.WriteRTP(pkt); err != nil {
			fmt.Println("Error writing RTP packet:", err)
		}
	}
}

func int16ToBytes(samples []int16) []byte {
	buf := make([]byte, len(samples)*2)
	for i, sample := range samples {
		buf[i*2] = byte(sample & 0xFF)
		buf[i*2+1] = byte((sample >> 8) & 0xFF)
	}
	return buf
}
