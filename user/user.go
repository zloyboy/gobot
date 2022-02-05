package user

const (
	Idx_country   = 0
	Idx_birth     = 1
	Idx_gender    = 2
	Idx_education = 3
	Idx_vacc_opin = 4
	Idx_orgn_opin = 5

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
