package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zloyboy/gobot/internal/user"

	_ "github.com/mattn/go-sqlite3"
)

const (
	idx_country = user.Idx_country
	idx_birth   = iota
	idx_gender
	idx_education
	idx_vacc_opin
	idx_orgn_opin
)

const (
	idx_year   = user.Idx_year
	idx_month  = user.Idx_month
	idx_sign   = user.Idx_sign
	idx_degree = user.Idx_degree
	idx_kind   = user.Idx_kind
	idx_effect = user.Idx_effect
)

type Dbase struct {
	mx sync.Mutex
	db *sql.DB
}

func createDb(db *sql.DB, path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Couldn't read %s", path)
		return err
	}
	scripts := strings.Split(string(content), ";\n")
	for _, script := range scripts {
		if len(script) > 0 {
			stmt, err := db.Prepare(script)
			if err == nil {
				stmt.Exec()
			} else {
				return err
			}
		}
	}
	return nil
}

func InitDb() *Dbase {
	db, err := sql.Open("sqlite3", "data/stat.db")
	if err != nil || db == nil {
		log.Print("Can't init sqlite3")
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Print("Can't ping sqlite3")
		return nil
	}
	var tableNname string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user'").Scan(&tableNname)
	if err != nil {
		log.Print("Table is not found, trying to create DB")
		err = createDb(db, "data/createDb.sql")
	}
	if err != nil {
		return nil
	}

	return &Dbase{db: db}
}

func (dbase *Dbase) ExistId(id int64) bool {
	res := 0
	err := dbase.db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE teleId=?)", id).Scan(&res)
	if err != nil || res != 1 {
		return false
	} else {
		return true
	}
}

func (dbase *Dbase) Insert(teleId int64, date string, usr user.UserData) error {
	dbase.mx.Lock()
	defer dbase.mx.Unlock()

	tx, _ := dbase.db.Begin()
	defer tx.Rollback()

	stmt, _ := tx.Prepare("INSERT INTO user" +
		"(teleId, created, modified, country, birth, gender, education, vaccineOpinion, originOpinion, countIll, countVac)" +
		"values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	stmt.Exec(teleId, date, date, usr.Base[idx_country], usr.Base[idx_birth], usr.Base[idx_gender],
		usr.Base[idx_education], usr.Base[idx_vacc_opin], usr.Base[idx_orgn_opin], usr.CountIll, usr.CountVac)

	age_group := user.GetAgeGroup(time.Now().Year() - usr.Base[idx_birth])
	haveIll := 0
	if 0 < usr.CountIll {
		haveIll = 1
	}
	haveVac := 0
	if 0 < usr.CountVac {
		haveVac = 1
	}

	stmt, _ = tx.Prepare("INSERT INTO userAgeGroup (id, created, teleId, have_ill, have_vac, age_group) values(?, ?, ?, ?, ?, ?)")
	stmt.Exec(nil, date, teleId, haveIll, haveVac, age_group)

	for i := 0; i < usr.CountIll; i++ {
		age := usr.Ill[i][idx_year] - usr.Base[idx_birth]
		stmt, _ = tx.Prepare("INSERT INTO userIllness (id, created, teleId, year, month, sign, degree, age) values(?, ?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(nil, date, teleId, usr.Ill[i][idx_year], usr.Ill[i][idx_month], usr.Ill[i][idx_sign], usr.Ill[i][idx_degree], age)
	}

	for i := 0; i < usr.CountVac; i++ {
		age := usr.Vac[i][idx_year] - usr.Base[idx_birth]
		stmt, _ = tx.Prepare("INSERT INTO userVaccine (id, created, teleId, year, month, kind, effect, age) values(?, ?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(nil, date, teleId, usr.Vac[i][idx_year], usr.Vac[i][idx_month], usr.Vac[i][idx_kind], usr.Vac[i][idx_effect], age)
	}

	tx.Commit()

	return nil
}

func (dbase *Dbase) ReadCaption(table string, arg ...int) [][2]string {
	var caps [][2]string
	query := fmt.Sprintf("SELECT rus from %s", table)
	rows, err := dbase.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var cap [2]string
	idx := 0
	if 0 < len(arg) {
		idx = arg[0]
	}
	for rows.Next() {
		err := rows.Scan(&cap[0])
		if err != nil {
			return nil
		}
		cap[1] = strconv.Itoa(idx)
		//log.Println(cap[0], cap[1])
		caps = append(caps, cap)
		idx++
	}
	err = rows.Err()
	if err != nil {
		return nil
	}

	return caps
}

func (dbase *Dbase) ReadCountAge() (int, int, int, [6][3]int) {
	var cntAll, cntIll, cntVac int
	var stat [6][3]int

	rows, _ := dbase.db.Query("SELECT age_group, have_ill, have_vac, COUNT(*) FROM userAgeGroup GROUP BY age_group, have_ill, have_vac")
	defer rows.Close()

	ageGrp, ill, vac, count := 0, 0, 0, 0
	for rows.Next() {
		err := rows.Scan(&ageGrp, &ill, &vac, &count)
		if err != nil {
			break
		}
		stat[ageGrp][0] += count
		cntAll += count
		if ill == 1 {
			stat[ageGrp][1] += count
			cntIll += count
		}
		if vac == 1 {
			stat[ageGrp][2] += count
			cntVac += count
		}
	}

	//log.Println(cntAll, cntIll, cntVac, stat)
	return cntAll, cntIll, cntVac, stat
}

func (dbase *Dbase) ReadOpinion() ([3]int, [3]int) {
	var vacOpn [3]int
	var orgOpn [3]int

	rows, _ := dbase.db.Query("SELECT vaccineOpinion, originOpinion, COUNT(*) FROM user GROUP BY vaccineOpinion, originOpinion")
	defer rows.Close()

	vac, org, count := 0, 0, 0
	for rows.Next() {
		err := rows.Scan(&vac, &org, &count)
		if err != nil {
			break
		}
		if 0 <= vac && vac <= 2 {
			vacOpn[vac] += count
		}
		if 0 <= org && org <= 2 {
			orgOpn[org] += count
		}
	}

	//log.Println(vacOpn, orgOpn)
	return vacOpn, orgOpn
}

func (dbase *Dbase) ReadChat() []int64 {
	rows, err := dbase.db.Query("SELECT id from chat")
	if err != nil {
		return nil
	}
	defer rows.Close()

	chat := make([]int64, 0)
	var id int64
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil
		}
		chat = append(chat, id)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}

	return chat
}

func (dbase *Dbase) AddChat(id int64) {
	stmt, _ := dbase.db.Prepare("INSERT INTO chat(id) VALUES(?)")
	defer stmt.Close()
	stmt.Exec(id)
}

func (dbase *Dbase) ExistChat(id int64) bool {
	res := 0
	err := dbase.db.QueryRow("SELECT EXISTS(SELECT 1 FROM chat WHERE id=?)", id).Scan(&res)
	if err != nil || res != 1 {
		return false
	} else {
		return true
	}
}
