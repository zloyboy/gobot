package user

const Idx_country = 0
const Idx_birth = 1
const Idx_gender = 2
const Idx_education = 3
const Idx_vacc_opin = 4
const Idx_orgn_opin = 5

const Idx_year = 0
const Idx_month = 1
const Idx_sign = 2
const Idx_degree = 3
const Idx_kind = 2
const Idx_effect = 3

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
