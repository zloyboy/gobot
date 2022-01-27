package telegram

import (
	"fmt"
	"strconv"
)

func (b *Bot) readStatFromDb() {
	b.stat.setCount(b.dbase.CountUsers(), b.dbase.CountIll())
	for i := 0; i < 6; i++ {
		if b.stat.setAgeCnt(i, b.dbase.CountAgeGroup(age_mid[i])) {
			b.stat.setIllCnt(i, b.dbase.CountAgeGroupIll(age_mid[i]))
		}
	}
}

var ages = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "60 ++"}
var age_mid = [6]int{15, 25, 35, 45, 55, 65}
var age_idx = map[string]int{ages[0]: 0, ages[1]: 1, ages[2]: 2, ages[3]: 3, ages[4]: 4, ages[5]: 5}

func inAges(age string) bool {
	switch age {
	case
		ages[0], ages[1], ages[2], ages[3], ages[4], ages[5]:
		return true
	}
	return false
}

type Static struct {
	cntAll, cntIll, cntNo int
	ages_stat             [6][2]int
}

func Stat() *Static {
	return &Static{
		cntAll:    0,
		cntIll:    0,
		cntNo:     0,
		ages_stat: [6][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}},
	}
}

func (s *Static) setCount(all int, ill int) {
	s.cntAll = all
	s.cntIll = ill
	s.cntNo = s.cntAll - s.cntIll
}

func (s *Static) setAgeCnt(idx int, cntAgeGroupAll int) bool {
	s.ages_stat[idx][0] = cntAgeGroupAll
	return 0 < s.ages_stat[idx][0]
}

func (s *Static) setIllCnt(idx int, cntAgeGroupIll int) {
	s.ages_stat[idx][1] = cntAgeGroupIll
}

func (s *Static) MakeStatic() string {
	var perYes, perNo float32
	if s.cntAll == 0 {
		perYes, perNo = 0, 0
	} else {
		perYes = float32(s.cntIll) / float32(s.cntAll) * 100
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
		"\n" + fmt.Sprintf("%.2f", perYes) + "%" + " переболело: " + strconv.Itoa(s.cntIll) + " из " + strconv.Itoa(s.cntAll) +
		"\n" + fmt.Sprintf("%.2f", perNo) + "%" + " не болело: " + strconv.Itoa(s.cntNo) + " из " + strconv.Itoa(s.cntAll) +
		"\nЗаболеваемость по возрастным группам:" + outAge
}

func (s *Static) RefreshStatic(idx int, ill int) {
	s.cntAll++
	s.cntIll += ill
	s.cntNo = s.cntAll - s.cntIll
	s.ages_stat[idx][0]++
	s.ages_stat[idx][1] += ill
}
