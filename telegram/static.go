package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zloyboy/gobot/user"
)

func (b *Bot) readStatFromDb() {
	b.stat.cntAll, b.stat.cntIll, b.stat.cntVac, b.stat.age_stat = b.dbase.ReadCountAge()
}

type Static struct {
	cntAll, cntIll, cntVac int
	age_stat               [6][3]int
}

func Stat() *Static {
	return &Static{
		cntAll:   0,
		cntIll:   0,
		cntVac:   0,
		age_stat: [6][3]int{},
	}
}

func (s *Static) MakeStatic() string {
	var ages = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "60 ++"}
	var perIll, perVac float32
	if s.cntAll == 0 {
		perIll, perVac = 0, 0
	} else {
		perIll = float32(s.cntIll) / float32(s.cntAll) * 100
		perVac = float32(s.cntVac) / float32(s.cntAll) * 100
	}
	ageIll := [6]float32{}
	ageVac := [6]float32{}
	statIll, statVac := "", ""
	for i := 0; i < 6; i++ {
		if 0 < s.age_stat[i][0] {
			ageIll[i] = float32(s.age_stat[i][1]) / float32(s.age_stat[i][0]) * 100
			ageVac[i] = float32(s.age_stat[i][2]) / float32(s.age_stat[i][0]) * 100
		}
		statIll += "\n" + ages[i] + " - " + fmt.Sprintf("%.2f", ageIll[i]) + "% - " + strconv.Itoa(s.age_stat[i][1]) + " из " + strconv.Itoa(s.age_stat[i][0])
		statVac += "\n" + ages[i] + " - " + fmt.Sprintf("%.2f", ageVac[i]) + "% - " + strconv.Itoa(s.age_stat[i][2]) + " из " + strconv.Itoa(s.age_stat[i][0])
	}
	return "Краткая статистика по COVID-19\nОпрошено: " + strconv.Itoa(s.cntAll) +
		"\n--------------------" +
		"\n" + fmt.Sprintf("%.2f", perIll) + "%" + " переболел: " + strconv.Itoa(s.cntIll) + " из " + strconv.Itoa(s.cntAll) + statIll +
		"\n--------------------" +
		"\n" + fmt.Sprintf("%.2f", perVac) + "%" + " вакцинировано: " + strconv.Itoa(s.cntVac) + " из " + strconv.Itoa(s.cntAll) + statVac +
		"\n--------------------"
}

func (s *Static) RefreshStatic(usr user.UserData) {
	ageGrp := user.GetAgeGroup(time.Now().Year() - usr.Base[idx_birth])
	haveIll := 0
	if 0 < usr.CountIll {
		haveIll = 1
	}
	haveVac := 0
	if 0 < usr.CountVac {
		haveVac = 1
	}

	s.cntAll++
	s.cntIll += haveIll
	s.cntVac += haveVac

	s.age_stat[ageGrp][0]++
	s.age_stat[ageGrp][1] += haveIll
	s.age_stat[ageGrp][2] += haveVac
}
