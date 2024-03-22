package shrimp

import (
	"fmt"
	"github.com/heilkit/tg"
	"github.com/heilkit/tg/scheduler"
	"math/rand/v2"
	"net"
	"net/http"
	"time"
)

type Balancer struct {
	Bots []*tg.Bot
}

func NewBalancer(tokens ...string) (*Balancer, error) {
	bots := make([]*tg.Bot, len(tokens))
	for i, token := range tokens {
		bot, err := tg.NewBot(tg.Settings{
			URL:     telegramapi(),
			Token:   token,
			Local:   tg.LocalMoving(),
			OnError: tg.OnErrorLog(),
			Client: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						Timeout:   300 * time.Second,
						KeepAlive: 300 * time.Second,
					}).DialContext,
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          100,
					IdleConnTimeout:       330 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 2 * time.Second,
				},
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       0,
			},
			Scheduler: scheduler.ExtraConservative(),
			Retries:   5,
		})
		if err != nil {
			return nil, fmt.Errorf("creating balancer %s %v", token, err)
		}
		bots[i] = bot
	}
	return &Balancer{Bots: bots}, nil
}

func (b *Balancer) Rand() *tg.Bot {
	return b.Bots[rand.IntN(len(b.Bots))]
}
