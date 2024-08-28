package news

import (
	"testing"
)

func TestWeatherPlugin_Do(t *testing.T) {
	p := NewPlugin(SetRssSource(""), SetTopN(10))
	t.Log(p.Do("北京"))
}
