package shrimp

import (
	"fmt"
	"log"
	"os"
)

var limit = fmt.Sprintf("%d", 1<<30)

func cyberdropdl() string {
	if env := os.Getenv("CYBERDROP_DL"); env != "" {
		return env
	}
	return "cyberdrop-dl"
}

func telegramapi() string {
	if env := os.Getenv("LOCAL_TG_API"); env != "" {
		return env
	}
	log.Print("WARN: no LOCAL_TG_API")
	return ""
}
