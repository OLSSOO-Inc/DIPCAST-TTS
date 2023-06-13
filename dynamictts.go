// 리드코리아 REST API를 통한 TTS변환 /usr/vt/rest/vtspeech --voice hyeryun --text  "안녕하세요,얼쑤팩토리입니다"  --lang Korean, --aformat mp3, --mp3rate 512, --ip 127.0.0.1, --port 7000, --srate 8000

package dynamictts

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os/exec"
)

// 화자와 TTS문자열 받아서 처리 Readspeaker
type ConfigReadspeaker struct {
	Speaker string
	Speak   string
	rsBin   string
}

// 화자와 속도, TTS문자열 받아서 처리 Naver tts-premium
type ConfigNavercpv struct {
	Speaker  string
	Speak    string
	Speed    string
	Apikeyid string
	Apikey   string
	baseUrl  string
}

// 파일 형식 지정
type Speech struct {
	bytes.Buffer
}

// Readspeaker
func SpeakReadspeaker(t ConfigReadspeaker) (*Speech, error) {

	// REST Command 실행
	args := []string{"--voice", t.Speaker, "--text", t.Speak, "--lang", "Korean", "--aformat", "mp3", "--mp3rate", "512", "--ip", "127.0.0.1", "--port", "7000", "--srate", "8000"}
	cmd := exec.Command(t.rsBin, args...)

	output, _ := cmd.CombinedOutput()

	speech := &Speech{}
	if _, err := io.Copy(&speech.Buffer, bytes.NewReader(output)); err != nil {
		return nil, err
	}
	cmd.Run()
	return speech, nil

}

// Naver tts-premium
func SpeakNavercpv(t ConfigNavercpv) (*Speech, error) {

	client := &http.Client{}
	data := url.Values{
		"speaker": {t.Speaker},
		"text":    {t.Speak},
		//"volume" : {t.Speak},
		"speed": {t.Speed},
		//"pitch" : {t.Speak},
		//"emotion" : {t.Speak},
		//"emotion-strength" : {t.Speak},
		//"alpha" : {t.Speak},
		//"end-pitch" : {t.Speak},

	}

	req, _ := http.NewRequest("POST", t.baseUrl, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-NCP-APIGW-API-KEY-ID", t.Apikeyid)
	req.Header.Add("X-NCP-APIGW-API-KEY", t.Apikey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	speech := &Speech{}
	if _, err := io.Copy(&speech.Buffer, res.Body); err != nil {
		return nil, err
	}

	return speech, nil
}
