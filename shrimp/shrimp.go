package shrimp

import (
	"errors"
	"fmt"
	"github.com/heilkit/tg"
	"github.com/heilkit/tg/tgmedia"
	"github.com/heilkit/tg/video"
	_ "gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

type Shrimp struct {
	API    *Balancer
	Config *Config
	Chat   *tg.Chat
}

func FromFile(filename string) (*Shrimp, error) {
	config, err := ConfigFromFile(filename)
	if err != nil {
		return nil, err
	}

	balancer, err := NewBalancer(config.Tokens...)
	if err != nil {
		return nil, err
	}

	return &Shrimp{
		API:    balancer,
		Chat:   &tg.Chat{ID: config.ChatId},
		Config: config,
	}, nil
}

func (shr *Shrimp) SyncConfig() error {
	if shr.Config.Source == "" {
		return nil
	}
	return errors.New("todo")
}

func (shr *Shrimp) Download(url string) error {
	cmd := exec.Command(cyberdropdl(),
		"--maximum-image-size", limit,
		"--maximum-video-size", limit,
		"--maximum-other-size", limit,
		"--download",
		url,
	)
	cmd.Dir = shr.Config.Dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func (shr *Shrimp) Upload(tag string) error {
	topic := &tg.SendOptions{}
	topic.ThreadID, _ = shr.Config.TagTopics[tag]
	if topic.ThreadID == 0 {
		thread, err := shr.API.Rand().CreateTopic(shr.Chat, &tg.Topic{Name: tag})
		if err != nil {
			return err
		}
		shr.Config.TagTopics[tag] = thread.ThreadID
		topic.ThreadID = thread.ThreadID
		if err := shr.Config.Update(); err != nil {
			return err
		}
	}

	return filepath.Walk(path.Join(shr.Config.Dir, "Downloads", "Cyberdrop-DL Downloads"),
		func(filepath string, info fs.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			defer func(name string) { _ = os.Remove(filepath) }(filepath)

			media := tgmedia.FromDisk(filepath)
			filename := path.Base(filepath)
			switch tmp := media.(type) {
			case *tg.Photo:
				FileName := fmt.Sprintf("@%s_%s", tag, filename)
				media = &tg.Document{File: tg.FromDisk(filepath), FileName: FileName}

			case *tg.Video:
				tmp.FileName = fmt.Sprintf("@%s_%s", tag, filename)
				media = tmp.With(video.ThumbnailAt(0.5))

			case *tg.Document:
				tmp.FileName = fmt.Sprintf("@%s_%s", tag, filename)
				media = tmp
			}

			if _, err := shr.API.Rand().Send(shr.Chat, media, topic); err != nil {
				log.Println(err)
			}

			return nil
		},
	)
}

func (shr *Shrimp) Loop() error {
	for _, profile := range shr.Config.Profiles {
		for _, url := range profile.URLs {
			downloadStart := time.Now()
			if err := shr.Download(url); err != nil {
				return err
			}
			log.Printf("DOWNLOAD %s (%s) in %s", profile.Tag, url, time.Since(downloadStart).String())

			uploadStart := time.Now()
			if err := shr.Upload(profile.Tag); err != nil {
				return err
			}
			log.Printf("UPLOAD %s (%s) in %s", profile.Tag, url, time.Since(uploadStart).String())
		}
	}
	return nil
}
