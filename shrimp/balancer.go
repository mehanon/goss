package shrimp

import (
	"fmt"
	"github.com/heilkit/tg"
	"github.com/heilkit/tg/scheduler"
	"math/rand/v2"
)

type Balancer struct {
	Bots []*tg.Bot
}

func NewBalancer(tokens ...string) (*Balancer, error) {
	bots := make([]*tg.Bot, len(tokens))
	for i, token := range tokens {
		bot, err := tg.NewBot(tg.Settings{
			URL:       telegramapi(),
			Token:     token,
			Local:     tg.LocalMoving(),
			OnError:   tg.OnErrorLog(),
			Scheduler: scheduler.ExtraConservative(),
			Retries:   10,
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
