package main

import "fmt"

/*
Фасад  используется для облегченного использования\доступа к сложной подсистеме
*/

type audioPlayer struct {
	//some audio data
}

func (a *audioPlayer) playAudio() {
	fmt.Println("pretend i am playing audio")
}

type videoPlayer struct {
	//some video data
}

func (v videoPlayer) playVideo() {
	fmt.Println("pretend i am playing video")
}

type screenManager struct {
	//something
}

func (s screenManager) screenManage() {
	fmt.Println("i make sure you see what you are supposed to")
}

type multimediaFacade struct { //facade for all of the above
	audioPlayer   *audioPlayer
	videoPlayer   *videoPlayer
	screenManager *screenManager
}

func newMultimediaFacade() *multimediaFacade {
	return &multimediaFacade{
		audioPlayer:   &audioPlayer{},
		videoPlayer:   &videoPlayer{},
		screenManager: &screenManager{},
	}
}

func (m *multimediaFacade) playMovie() { //simplified methods for clients
	m.audioPlayer.playAudio()
	m.videoPlayer.playVideo()
	m.screenManager.screenManage()
}

func main() {
	mMediaSystem := newMultimediaFacade()

	mMediaSystem.playMovie()
}
