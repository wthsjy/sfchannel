package sfchannel

import (
	"testing"
	"time"
)

func TestSafeChannel(t *testing.T) {
	ch := NewChannel[string](10)
	go func() {
		for {
			time.Sleep(time.Second)
			err := ch.Produce(time.Now().String())
			t.Log("err:", err)
			success, err := ch.ProduceOrDiscard(time.Now().String())
			t.Log("success:", success, ", err:", err)
		}
	}()
	go func() {
		time.Sleep(time.Second * 12)
		ch.Close()
	}()
	go func() {
		for {
			select {
			case val, ok := <-ch.Consume():
				time.Sleep(time.Second * 2)
				t.Log("val:", val, ", ok:", ok)
				if !ok {
					return
				}
			}
		}

	}()

	select {}

}
