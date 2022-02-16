package user

const (
	Idx_country = 0
	Idx_birth   = iota
	Idx_gender
	Idx_education
	Idx_vacc_opin
	Idx_orgn_opin
)

const (
	Idx_year   = 0
	Idx_month  = 1
	Idx_sign   = 2
	Idx_degree = 3
	Idx_kind   = 2
	Idx_effect = 3
)

type UserData struct {
	Base     [6]int
	CountIll int
	Ill      [][4]int
	CountVac int
	Vac      [][4]int
}

func MakeUser() UserData {
	return UserData{
		Base:     [6]int{-1, -1, -1, -1, -1, -1},
		CountIll: 0,
		Ill:      nil,
		CountVac: 0,
		Vac:      nil}
}

func MakeSubUser() [4]int {
	return [4]int{-1, -1, -1, -1}
}

func GetAgeGroup(age int) int {
	switch {
	case age < 20:
		return 0
	case 20 <= age && age < 30:
		return 1
	case 30 <= age && age < 40:
		return 2
	case 40 <= age && age < 50:
		return 3
	case 50 <= age && age < 60:
		return 4
	case 60 <= age:
		return 5
	default:
		return -1
	}
}
