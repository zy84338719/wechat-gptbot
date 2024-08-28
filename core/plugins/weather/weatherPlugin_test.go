package weather

import (
	"testing"
)

func TestWeatherPlugin_Do(t *testing.T) {
	plugin := NewPlugin()
	t.Log(plugin.Do("成都"))
}
