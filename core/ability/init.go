package ability

import (
	"strings"
	"wechat-gptbot/core/ability/proto"
)

var keywords []string

type Ability interface {
	TextFunc(text string) string
	help() *proto.AbilityHelpInfo
}

var AbilityMap map[string]Ability

func init() {

	kfcThurday := newKfc()
	gaso := newGasoline()
	h := newHotday()
	car := newLimitCar()
	s := newSuoxie()
	t := newTiangou()
	rate := newExchangeRate()
	w := newWeather()
	xz := newStarTalk()
	f := newFilm()
	AbilityMap = map[string]Ability{
		kfcThurday.help().Keyword: kfcThurday,
		gaso.help().Keyword:       gaso,
		h.help().Keyword:          h,
		car.help().Keyword:        car,
		s.help().Keyword:          s,
		t.help().Keyword:          t,
		rate.help().Keyword:       rate,
		w.help().Keyword:          w,
		xz.help().Keyword:         xz,
		f.help().Keyword:          f,
		qiuqian{}.help().Keyword:  newQiuqian(),
	}
	for k, _ := range AbilityMap {
		keywords = append(keywords, k)
	}
}

func Contains(keyword string) Ability {
	for key, v := range AbilityMap {
		if strings.Contains(keyword, key) {
			return v
			break
		}
	}
	return nil
}
