package domain

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"os"
)

type Video struct {
	Uuid   string `json:"uuid"`
	Path   string `json: "path"`
	Status string `json: "status"`
}

func (video *Video) Unmarshal(payload []byte) Video {
	err := json.Unmarshal(payload, &video)

	if err != nil {
		panic(err)
	}
	return *video
}

func (video *Video) Download(bucketName string, storagePath string) (Video, error) {
	ctx := context.TODO()

	//client, err := storage.NewClient(ctx)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	httpClient := &http.Client{Transport: transCfg}

	client, err := storage.NewClient(ctx, option.WithEndpoint("https://fake:4443/storage/v1/"), option.WithHTTPClient(httpClient))


	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
		return *video, err
	}


	bkt := client.Bucket(bucketName)

	obj := bkt.Object(video.Path)

	r, err := obj.NewReader(ctx)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
		return *video, err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
		return *video, err
	}

	f, err := os.Create(storagePath+"/"+video.Uuid+".mp4");

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
		return *video, err
	}

	_, err = f.Write(body)

	if err != nil {
		video.Status = "error"
		fmt.Println(err.Error())
		return *video, err
	}

	defer f.Close()

	fmt.Println("Video ", video.Uuid, "has been stored")

	return *video, nil
}
