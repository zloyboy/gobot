package telegram

import (
	"fmt"
	"strconv"
)

func (b *Bot) readStatFromDb() {
	b.stat.setCount(b.dbase.CountUsers(), b.dbase.CountIll())
}

type Static struct {
	cntAll, cntIll, cntNo int
}

func Stat() *Static {
	return &Static{
		cntAll: 0,
		cntIll: 0,
		cntNo:  0,
	}
}

func (s *Static) setCount(all int, ill int) {
	s.cntAll = all
	s.cntIll = ill
	s.cntNo = s.cntAll - s.cntIll
}

func (s *Static) MakeStatic() string {
	var perYes, perNo float32
	if s.cntAll == 0 {
		perYes, perNo = 0, 0
	} else {
		perYes = float32(s.cntIll) / float32(s.cntAll) * 100
		perNo = float32(s.cntNo) / float32(s.cntAll) * 100.0
	}
	/*perAge := [6]float32{0, 0, 0, 0, 0, 0}
	var outAge = ""
	for i := 0; i < 6; i++ {
		if 0 < s.ages_stat[i][0] {
			perAge[i] = float32(s.ages_stat[i][1]) / float32(s.ages_stat[i][0]) * 100
		}
		outAge += "\n" + ages[i] + " - " + fmt.Sprintf("%.2f", perAge[i]) + "% - " + strconv.Itoa(s.ages_stat[i][1]) + " из " + strconv.Itoa(s.ages_stat[i][0])
	}*/
	return "Независимая статистика по COVID-19\nОпрошено: " + strconv.Itoa(s.cntAll) +
		"\n" + fmt.Sprintf("%.2f", perYes) + "%" + " переболело: " + strconv.Itoa(s.cntIll) + " из " + strconv.Itoa(s.cntAll) +
		"\n" + fmt.Sprintf("%.2f", perNo) + "%" + " не болело: " + strconv.Itoa(s.cntNo) + " из " + strconv.Itoa(s.cntAll) +
		"\n--------------------"
}

func (s *Static) RefreshStatic(ill int) {
	s.cntAll++
	s.cntIll += ill
	s.cntNo = s.cntAll - s.cntIll
}
