package static

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/zloyboy/gobot/internal/database"
	"github.com/zloyboy/gobot/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	chart "github.com/zloyboy/chart"
)

const (
	idxAll = 0
	idxIll = 1
	idxVac = 2

	idxBirth  = user.Idx_birth
	idxGender = user.Idx_gender
)

const (
	idx_birth = iota
	idx_gender
	idx_educat
	idx_vacopn
	idx_orgopn
)

type Static struct {
	mx                                                     sync.RWMutex
	cntAll, cntIll, cntVac                                 [2]int
	age_stat                                               [6][2][3]int
	vacOpn, orgOpn                                         [3]int
	dbase                                                  *database.Dbase
	chartAll, chartIll, chartVac, chartVacOpn, chartOrgOpn tgbotapi.FileBytes
}

func Stat(db *database.Dbase) *Static {
	return &Static{
		cntAll:   [2]int{0, 0},
		cntIll:   [2]int{0, 0},
		cntVac:   [2]int{0, 0},
		age_stat: [6][2][3]int{},
		vacOpn:   [3]int{},
		orgOpn:   [3]int{},
		dbase:    db,
		chartAll: tgbotapi.FileBytes{Bytes: nil},
		chartIll: tgbotapi.FileBytes{Bytes: nil},
		chartVac: tgbotapi.FileBytes{Bytes: nil},
	}
}

func (s *Static) CntAll() int {
	return s.cntAll[0] + s.cntAll[1]
}

func (s *Static) CntIll() int {
	return s.cntIll[0] + s.cntIll[1]
}

func (s *Static) PercentIll() float64 {
	return float64(s.CntIll()) / float64(s.CntAll()) * 100
}

func (s *Static) CntVac() int {
	return s.cntVac[0] + s.cntVac[1]
}

func (s *Static) PercentVac() float64 {
	return float64(s.CntVac()) / float64(s.CntAll()) * 100
}

func (s *Static) CntAgeAll(idx int) int {
	return s.age_stat[idx][0][idxAll] + s.age_stat[idx][1][idxAll]
}

func (s *Static) CntAgeIll(idx int) int {
	return s.age_stat[idx][0][idxIll] + s.age_stat[idx][1][idxIll]
}

func (s *Static) PercentAgeIll(idx int) int {
	return int(float64(s.CntAgeIll(idx)) / float64(s.CntAgeAll(idx)) * 100)
}

func (s *Static) CntAgeVac(idx int) int {
	return s.age_stat[idx][0][idxVac] + s.age_stat[idx][1][idxVac]
}

func (s *Static) PercentAgeVac(idx int) int {
	return int(float64(s.CntAgeVac(idx)) / float64(s.CntAgeAll(idx)) * 100)
}

func (s *Static) GetChartAll() tgbotapi.FileBytes {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.chartAll
}

func (s *Static) GetChartIll() tgbotapi.FileBytes {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.chartIll
}

func (s *Static) GetChartVac() tgbotapi.FileBytes {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.chartVac
}

func (s *Static) GetChartVacOpn() tgbotapi.FileBytes {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.chartVacOpn
}

func (s *Static) GetChartOrgOpn() tgbotapi.FileBytes {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.chartOrgOpn
}

func (s *Static) ReadStatFromDb() {
	s.cntAll, s.cntIll, s.cntVac, s.age_stat = s.dbase.ReadCountAge()
	s.makeChartAll()
	s.makeChartIll()
	s.makeChartVac()
	s.vacOpn, s.orgOpn = s.dbase.ReadOpinion()
	s.makeChartVacOpn()
	s.makeChartOrgOpn()
}

func (s *Static) RefreshStatic(usr user.UserData) {
	gen := usr.Base[idxGender]
	if gen != 0 && gen != 1 {
		log.Println("wrong gender", gen)
		return
	}

	ageGrp := user.GetAgeGroup(time.Now().Year() - usr.Base[idxBirth])
	haveIll := 0
	if 0 < usr.CountIll {
		haveIll = 1
	}
	haveVac := 0
	if 0 < usr.CountVac {
		haveVac = 1
	}

	s.mx.Lock()

	s.cntAll[gen]++
	s.cntIll[gen] += haveIll
	s.cntVac[gen] += haveVac

	s.age_stat[ageGrp][gen][0]++
	s.age_stat[ageGrp][gen][1] += haveIll
	s.age_stat[ageGrp][gen][2] += haveVac

	vac := usr.Base[idx_vacopn]
	if 0 <= vac && vac <= 2 {
		s.vacOpn[vac]++
	}
	org := usr.Base[idx_orgopn]
	if 0 <= org && org <= 2 {
		s.orgOpn[org]++
	}

	s.makeChartAll()
	s.makeChartIll()
	s.makeChartVac()
	s.makeChartVacOpn()
	s.makeChartOrgOpn()

	s.mx.Unlock()
}

func (s *Static) makeChart(title string, values []chart.Value, bWidth, bSpace int) chart.BarChart {
	return chart.BarChart{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    30,
				Bottom: 30,
				Left:   5,
				Right:  5,
			},
		},
		Width:      500,
		Height:     200,
		BarWidth:   bWidth,
		BarSpacing: bSpace,
		YAxis: chart.YAxis{
			ValueFormatter: chart.PercentValueFormatter,
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: 1,
			},
		},
		Bars: values,
	}
}

func (s *Static) makeChart2(title string, values []chart.DoubleValue, bWidth, bSpace, top, subtop int, first bool) chart.BarChart2 {
	return chart.BarChart2{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    35,
				Bottom: 45,
				Left:   5,
				Right:  5,
			},
		},
		Width:       560,
		Height:      300,
		BarWidth:    bWidth,
		BarSpacing:  bSpace,
		LabelFirst:  first,
		LabelTop:    top,
		SubLabelTop: subtop,
		YAxis: chart.YAxis{
			ValueFormatter: chart.PercentValueFormatter,
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: 1,
			},
		},
		Bars: values,
	}
}

func (s *Static) makeChartAll() {
	var perIll, perVac [2]float64
	if s.cntAll[0] == 0 {
		perIll[0], perVac[0] = 0, 0
	} else {
		perIll[0] = float64(s.cntIll[0]) / float64(s.cntAll[0])
		perVac[0] = float64(s.cntVac[0]) / float64(s.cntAll[0])
	}
	if s.cntAll[1] == 0 {
		perIll[1], perVac[1] = 0, 0
	} else {
		perIll[1] = float64(s.cntIll[1]) / float64(s.cntAll[1])
		perVac[1] = float64(s.cntVac[1]) / float64(s.cntAll[1])
	}

	var labelHead [2]string
	labelHead[0] = fmt.Sprintf("Переболело %.2f%% (%d из %d)", s.PercentIll(), s.CntIll(), s.CntAll())
	labelHead[1] = fmt.Sprintf("Вакцинировано %.2f%% (%d из %d)", s.PercentVac(), s.CntVac(), s.CntAll())

	var labelIll, labelVac [2]string
	labelIll[0] = fmt.Sprintf("женщин %.2f%%\n(%d из %d)", perIll[0]*100, s.cntIll[0], s.cntAll[0])
	labelIll[1] = fmt.Sprintf("мужчин %.2f%%\n(%d из %d)", perIll[1]*100, s.cntIll[1], s.cntAll[1])
	labelVac[0] = fmt.Sprintf("женщин %.2f%%\n(%d из %d)", perVac[0]*100, s.cntVac[0], s.cntAll[0])
	labelVac[1] = fmt.Sprintf("мужчин %.2f%%\n(%d из %d)", perVac[1]*100, s.cntVac[1], s.cntAll[1])

	values := []chart.DoubleValue{
		{Label: labelHead[0], Lab: labelIll, Val: perIll},
		{Label: labelHead[1], Lab: labelVac, Val: perVac},
	}

	response := s.makeChart2(fmt.Sprintf("Краткая статистика по COVID-19, опрошено %d", s.CntAll()), values, 230, 10, 30, 20, true)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartAll = tgbotapi.FileBytes{Name: "common.png", Bytes: buffer.Bytes()}
}

var age = [6]string{"до 20 лет", "20-29 лет", "30-39 лет", "40-49 лет", "50-59 лет", "после 60"}

func (s *Static) makeChartIll() {
	values := make([]chart.DoubleValue, 0, 6)

	for i := 0; i < 6; i++ {
		var ageIll [2]float64
		labelAge := fmt.Sprintf("%s\n%d%% (%d/%d)", age[i], s.PercentAgeIll(i), s.CntAgeIll(i), s.CntAgeAll(i))
		labelIll := [2]string{"ж\n0%", "м\n0%"}
		if 0 < s.age_stat[i][0][idxIll] {
			ageIll[0] = float64(s.age_stat[i][0][idxIll]) / float64(s.age_stat[i][0][idxAll])
			labelIll[0] = fmt.Sprintf("ж\n%d%%", int(ageIll[0]*100))
		}
		if 0 < s.age_stat[i][1][idxIll] {
			ageIll[1] = float64(s.age_stat[i][1][idxIll]) / float64(s.age_stat[i][1][idxAll])
			labelIll[1] = fmt.Sprintf("м\n%d%%", int(ageIll[1]*100))
		}
		values = append(values, chart.DoubleValue{Label: labelAge, Lab: labelIll, Val: ageIll})
	}

	response := s.makeChart2("Заболеваемость по возрастам", values, 75, 10, 15, 35, false)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartIll = tgbotapi.FileBytes{Name: "ill.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartVac() {
	values := make([]chart.DoubleValue, 0, 6)

	for i := 0; i < 6; i++ {
		var ageVac [2]float64
		labelAge := fmt.Sprintf("%s\n%d%% (%d/%d)", age[i], s.PercentAgeVac(i), s.CntAgeVac(i), s.CntAgeAll(i))
		labelVac := [2]string{"ж\n0%", "м\n0%"}
		if 0 < s.age_stat[i][0][idxVac] {
			ageVac[0] = float64(s.age_stat[i][0][idxVac]) / float64(s.age_stat[i][0][idxAll])
			labelVac[0] = fmt.Sprintf("ж\n%d%%", int(ageVac[0]*100))
		}
		if 0 < s.age_stat[i][1][idxVac] {
			ageVac[1] = float64(s.age_stat[i][1][idxVac]) / float64(s.age_stat[i][1][idxAll])
			labelVac[1] = fmt.Sprintf("м\n%d%%", int(ageVac[1]*100))
		}
		values = append(values, chart.DoubleValue{Label: labelAge, Lab: labelVac, Val: ageVac})
	}

	response := s.makeChart2("Вакцинация по возрастам", values, 75, 10, 15, 35, false)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartVac = tgbotapi.FileBytes{Name: "vac.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartVacOpn() {
	var perYes, perNo, perBad float64
	if s.CntAll() == 0 {
		perYes, perNo, perBad = 0, 0, 0
	} else {
		perYes = float64(s.vacOpn[0]) / float64(s.CntAll())
		perNo = float64(s.vacOpn[1]) / float64(s.CntAll())
		perBad = float64(s.vacOpn[2]) / float64(s.CntAll())
	}

	labelYes := fmt.Sprintf("Помогают: %.2f %%\n(%d из %d)", perYes*100, s.vacOpn[0], s.CntAll())
	labelNo := fmt.Sprintf("Бесп-зны: %.2f %%\n(%d из %d)", perNo*100, s.vacOpn[1], s.CntAll())
	labelBad := fmt.Sprintf("Опасны: %.2f %%\n(%d из %d)", perBad*100, s.vacOpn[2], s.CntAll())

	values := []chart.Value{{Label: labelYes, Value: perYes}, {Label: labelNo, Value: perNo}, {Label: labelBad, Value: perBad}}
	response := s.makeChart("Мнение о полезности вакцин", values, 120, 50)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartVacOpn = tgbotapi.FileBytes{Name: "vacopn.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartOrgOpn() {
	var perNat, perHum, perUnk float64
	if s.CntAll() == 0 {
		perNat, perHum, perUnk = 0, 0, 0
	} else {
		perNat = float64(s.orgOpn[0]) / float64(s.CntAll())
		perHum = float64(s.orgOpn[1]) / float64(s.CntAll())
		perUnk = float64(s.orgOpn[2]) / float64(s.CntAll())
	}

	labelNat := fmt.Sprintf("Природа: %.2f %%\n(%d из %d)", perNat*100, s.orgOpn[0], s.CntAll())
	labelHum := fmt.Sprintf("Люди: %.2f %%\n(%d из %d)", perHum*100, s.orgOpn[1], s.CntAll())
	labelUnk := fmt.Sprintf("Не знаю: %.2f %%\n(%d из %d)", perUnk*100, s.orgOpn[2], s.CntAll())

	values := []chart.Value{{Label: labelNat, Value: perNat}, {Label: labelHum, Value: perHum}, {Label: labelUnk, Value: perUnk}}
	response := s.makeChart("Мнение о происхождении вируса", values, 120, 50)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartOrgOpn = tgbotapi.FileBytes{Name: "orgopn.png", Bytes: buffer.Bytes()}
}
