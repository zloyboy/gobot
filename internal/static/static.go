package static

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/zloyboy/gobot/internal/database"
	"github.com/zloyboy/gobot/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	chart "github.com/wcharczuk/go-chart"
)

const (
	idxAll = 0
	idxIll = 1
	idxVac = 2

	idxBirth = user.Idx_birth
)

type Static struct {
	mx                                                     sync.RWMutex
	cntAll, cntIll, cntVac                                 int
	age_stat                                               [6][3]int
	vacOpn, orgOpn                                         [3]int
	dbase                                                  *database.Dbase
	chartAll, chartIll, chartVac, chartVacOpn, chartOrgOpn tgbotapi.FileBytes
}

func Stat(db *database.Dbase) *Static {
	return &Static{
		cntAll:   0,
		cntIll:   0,
		cntVac:   0,
		age_stat: [6][3]int{},
		vacOpn:   [3]int{},
		orgOpn:   [3]int{},
		dbase:    db,
		chartAll: tgbotapi.FileBytes{Bytes: nil},
		chartIll: tgbotapi.FileBytes{Bytes: nil},
		chartVac: tgbotapi.FileBytes{Bytes: nil},
	}
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

	s.cntAll++
	s.cntIll += haveIll
	s.cntVac += haveVac

	s.age_stat[ageGrp][0]++
	s.age_stat[ageGrp][1] += haveIll
	s.age_stat[ageGrp][2] += haveVac

	s.makeChartAll()
	s.makeChartIll()
	s.makeChartVac()
	s.makeChartVacOpn()
	s.makeChartOrgOpn()

	s.mx.Unlock()
}

func (s *Static) makeChart(title string, values []chart.Value, bWidth, bSpace int) chart.BarChart {
	return chart.BarChart{
		Title:      title,
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 30,
				Left:   10,
				Right:  20,
			},
		},
		Width:      500,
		Height:     200,
		BarWidth:   bWidth,
		BarSpacing: bSpace,
		XAxis:      chart.StyleShow(),
		YAxis: chart.YAxis{
			Style:          chart.StyleShow(),
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
	var perIll, perVac float64
	if s.cntAll == 0 {
		perIll, perVac = 0, 0
	} else {
		perIll = float64(s.cntIll) / float64(s.cntAll)
		perVac = float64(s.cntVac) / float64(s.cntAll)
	}

	labelIll := fmt.Sprintf("Переболело %.2f %%\n%d из %d", perIll*100, s.cntIll, s.cntAll)
	labelVac := fmt.Sprintf("Вакцинировано %.2f %%\n%d из %d", perVac*100, s.cntVac, s.cntAll)

	values := []chart.Value{{Label: labelIll, Value: perIll}, {Label: labelVac, Value: perVac}}
	response := s.makeChart(fmt.Sprintf("Краткая статистика по COVID-19, опрошено %d", s.cntAll), values, 200, 10)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartAll = tgbotapi.FileBytes{Name: "common.png", Bytes: buffer.Bytes()}
}

var age = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "после 60"}

func (s *Static) makeChartIll() {
	values := make([]chart.Value, 0, 6)

	for i := 0; i < 6; i++ {
		var ageIll float64
		label := age[i]
		if 0 < s.age_stat[i][idxIll] {
			label += "\n" + fmt.Sprintf("%d из %d", s.age_stat[i][idxIll], s.age_stat[i][idxAll])
			ageIll = float64(s.age_stat[i][idxIll]) / float64(s.age_stat[i][idxAll])
		}
		values = append(values, chart.Value{Label: label, Value: ageIll})
	}

	response := s.makeChart("Заболеваемость по возрастам", values, 50, 30)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartIll = tgbotapi.FileBytes{Name: "ill.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartVac() {
	values := make([]chart.Value, 0, 6)

	for i := 0; i < 6; i++ {
		var ageVac float64
		label := age[i]
		if 0 < s.age_stat[i][idxVac] {
			label += "\n" + fmt.Sprintf("%d из %d", s.age_stat[i][idxVac], s.age_stat[i][idxAll])
			ageVac = float64(s.age_stat[i][idxVac]) / float64(s.age_stat[i][idxAll])
		}
		values = append(values, chart.Value{Label: label, Value: ageVac})
	}

	response := s.makeChart("Вакцинация по возрастам", values, 50, 30)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartVac = tgbotapi.FileBytes{Name: "vac.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartVacOpn() {
	var perYes, perNo, perBad float64
	if s.cntAll == 0 {
		perYes, perNo, perBad = 0, 0, 0
	} else {
		perYes = float64(s.vacOpn[0]) / float64(s.cntAll)
		perNo = float64(s.vacOpn[1]) / float64(s.cntAll)
		perBad = float64(s.vacOpn[2]) / float64(s.cntAll)
	}

	labelYes := fmt.Sprintf("Помогают: %.2f %%\n%d из %d", perYes*100, s.vacOpn[0], s.cntAll)
	labelNo := fmt.Sprintf("Бесп-зны: %.2f %%\n%d из %d", perNo*100, s.vacOpn[1], s.cntAll)
	labelBad := fmt.Sprintf("Опасны: %.2f %%\n%d из %d", perBad*100, s.vacOpn[2], s.cntAll)

	values := []chart.Value{{Label: labelYes, Value: perYes}, {Label: labelNo, Value: perNo}, {Label: labelBad, Value: perBad}}
	response := s.makeChart("Мнение о полезности вакцин", values, 120, 50)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartVacOpn = tgbotapi.FileBytes{Name: "vacopn.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartOrgOpn() {
	var perNat, perHum, perUnk float64
	if s.cntAll == 0 {
		perNat, perHum, perUnk = 0, 0, 0
	} else {
		perNat = float64(s.orgOpn[0]) / float64(s.cntAll)
		perHum = float64(s.orgOpn[1]) / float64(s.cntAll)
		perUnk = float64(s.orgOpn[2]) / float64(s.cntAll)
	}

	labelNat := fmt.Sprintf("Природа: %.2f %%\n%d из %d", perNat*100, s.orgOpn[0], s.cntAll)
	labelHum := fmt.Sprintf("Люди: %.2f %%\n%d из %d", perHum*100, s.orgOpn[1], s.cntAll)
	labelUnk := fmt.Sprintf("Не знаю: %.2f %%\n%d из %d", perUnk*100, s.orgOpn[2], s.cntAll)

	values := []chart.Value{{Label: labelNat, Value: perNat}, {Label: labelHum, Value: perHum}, {Label: labelUnk, Value: perUnk}}
	response := s.makeChart("Мнение о происхождении вируса", values, 120, 50)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartOrgOpn = tgbotapi.FileBytes{Name: "orgopn.png", Bytes: buffer.Bytes()}
}
