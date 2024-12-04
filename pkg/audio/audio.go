package audio

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func PlayAudio(data []byte) {
	st, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer st.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(st, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func SaveAudio(data []byte, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	_ = n
}