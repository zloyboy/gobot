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
	mx                           sync.RWMutex
	cntAll, cntIll, cntVac       int
	age_stat                     [6][3]int
	dbase                        *database.Dbase
	chartAll, chartIll, chartVac tgbotapi.FileBytes
}

func Stat(db *database.Dbase) *Static {
	return &Static{
		cntAll:   0,
		cntIll:   0,
		cntVac:   0,
		age_stat: [6][3]int{},
		dbase:    db,
		chartAll: tgbotapi.FileBytes{Bytes: nil},
		chartIll: tgbotapi.FileBytes{Bytes: nil},
		chartVac: tgbotapi.FileBytes{Bytes: nil},
	}
}

func (s *Static) ReadStatFromDb() {
	s.cntAll, s.cntIll, s.cntVac, s.age_stat = s.dbase.ReadCountAge()
	s.makeChartAll()
	s.makeChartIll()
	s.makeChartVac()
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

	s.mx.Unlock()
}

func (s *Static) makeAgeChart(title string, values []chart.Value) chart.BarChart {
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
		BarWidth:   50,
		BarSpacing: 30,
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

	response := s.makeAgeChart("Заболеваемость по возрастам", values)
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

	response := s.makeAgeChart("Вакцинация по возрастам", values)
	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartVac = tgbotapi.FileBytes{Name: "vac.png", Bytes: buffer.Bytes()}
}

func (s *Static) makeChartAll() {
	var perIll, perVac float64
	if s.cntAll == 0 {
		perIll, perVac = 0, 0
	} else {
		perIll = float64(s.cntIll) / float64(s.cntAll)
		perVac = float64(s.cntVac) / float64(s.cntAll)
	}

	values := make([]chart.Value, 0, 2)
	label := fmt.Sprintf("Переболело %.2f %%\n%d из %d", perIll*100, s.cntIll, s.cntAll)
	values = append(values, chart.Value{Label: label, Value: perIll})
	label = fmt.Sprintf("Вакцинировано %.2f %%\n%d из %d", perVac*100, s.cntVac, s.cntAll)
	values = append(values, chart.Value{Label: label, Value: perVac})

	response := chart.BarChart{
		Title:      fmt.Sprintf("Краткая статистика по COVID-19, опрошено %d", s.cntAll),
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
		BarWidth:   200,
		BarSpacing: 30,
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

	buffer := &bytes.Buffer{}
	response.Render(chart.PNG, buffer)
	s.chartAll = tgbotapi.FileBytes{Name: "common.png", Bytes: buffer.Bytes()}
}
