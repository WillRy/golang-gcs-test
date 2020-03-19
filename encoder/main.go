package main
import (
	v "encoder/domain"
	// log "github.com/sirupsen/logrus"
)

func main()  {

	var video v.Video

	data := []byte("{\"uuid\":\"abc123\", \"path\":\"convite.txt\",\"status\":\"pending\"}")

	video.Unmarshal(data)
	// log.Info(video.Path)

	video.Download("codeeducationtest","/tmp")

}
