package ability

import "testing"

func Test_hotday_textFunc(t *testing.T) {
	h := newHotday()
	textFunc := h.TextFunc("综合")
	t.Log(textFunc)
}
