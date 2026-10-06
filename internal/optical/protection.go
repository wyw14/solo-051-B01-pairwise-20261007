package optical

import "github.com/wyw14/cry-161/internal/model"

func (b *Broadcaster) PublishThen(sample model.ThicknessSample, action func() error) error {
	b.Publish(sample)
	if action == nil {
		return nil
	}
	return action()
}
