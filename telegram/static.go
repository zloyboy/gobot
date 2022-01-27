package telegram

import (
	"fmt"
	"strconv"
)

func (b *Bot) readStatFromDb() {
	b.stat.setCount(b.dbase.CountUsers(), b.dbase.CountRes())
	for i := 0; i < 6; i++ {
		if b.stat.setAgeCnt(i, b.dbase.CountAge(i*10+15)) {
			b.stat.setIllCnt(i, b.dbase.CountAgeRes(i*10+15))
		}
	}
}

var ages = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "60 ++"}
var ageGroup = map[string]int{ages[0]: 15, ages[1]: 25, ages[2]: 35, ages[3]: 45, ages[4]: 55, ages[5]: 65}

func inAges(age string) bool {
	switch age {
	case
		ages[0], ages[1], ages[2], ages[3], ages[4], ages[5]:
		return true
	}
	return false
}

type Static struct {
	cntAll, cntYes, cntNo int
	ages_stat             [6][2]int
}

func Stat() *Static {
	return &Static{
		cntAll:    0,
		cntYes:    0,
		cntNo:     0,
		ages_stat: [6][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}},
	}
}

func (s *Static) setCount(all int, ill int) {
	s.cntAll = all
	s.cntYes = ill
	s.cntNo = s.cntAll - s.cntYes
}

func (s *Static) setAgeCnt(idx int, cntAll int) bool {
	s.ages_stat[idx][0] = cntAll
	return 0 < s.ages_stat[idx][0]
}

func (s *Static) setIllCnt(idx int, cntIll int) {
	s.ages_stat[idx][1] = cntIll
}

func (s *Static) MakeStatic() string {
	var perYes, perNo float32
	if s.cntAll == 0 {
		perYes, perNo = 0, 0
	} else {
		perYes = float32(s.cntYes) / float32(s.cntAll) * 100
		perNo = float32(s.cntNo) / float32(s.cntAll) * 100.0
	}
	perAge := [6]float32{0, 0, 0, 0, 0, 0}
	var outAge = ""
	for i := 0; i < 6; i++ {
		if 0 < s.ages_stat[i][0] {
			perAge[i] = float32(s.ages_stat[i][1]) / float32(s.ages_stat[i][0]) * 100
		}
		outAge += "\n" + ages[i] + " - " + fmt.Sprintf("%.2f", perAge[i]) + "% - " + strconv.Itoa(s.ages_stat[i][1]) + " из " + strconv.Itoa(s.ages_stat[i][0])
	}
	return "Независимая статистика по COVID-19\nОпрошено: " + strconv.Itoa(s.cntAll) +
		"\n" + fmt.Sprintf("%.2f", perYes) + "%" + " переболело: " + strconv.Itoa(s.cntYes) + " из " + strconv.Itoa(s.cntAll) +
		"\n" + fmt.Sprintf("%.2f", perNo) + "%" + " не болело: " + strconv.Itoa(s.cntNo) + " из " + strconv.Itoa(s.cntAll) +
		"\nЗаболеваемость по возрастным группам:" + outAge
}
