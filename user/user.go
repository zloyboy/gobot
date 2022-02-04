package user

type UserIll struct {
	Year   int
	Month  int
	Sign   int
	Degree int
}

func MakeIll() UserIll {
	return UserIll{
		Year:   0,
		Month:  0,
		Sign:   -1,
		Degree: -1}
}

type UserVac struct {
	Year   int
	Month  int
	Kind   int
	Effect int
}

func MakeVac() UserVac {
	return UserVac{
		Year:   0,
		Month:  0,
		Kind:   -1,
		Effect: -1}
}

type UserData struct {
	Country   int
	Birth     int
	Gender    int
	Education int
	Vaccine   int
	Origin    int
	CountIll  int
	Ill       []UserIll
	CountVac  int
	Vac       []UserVac
}

func MakeUser() UserData {
	return UserData{
		Country:   -1,
		Birth:     -1,
		Gender:    -1,
		Education: -1,
		Vaccine:   -1,
		Origin:    -1,
		CountIll:  0,
		Ill:       nil,
		CountVac:  0,
		Vac:       nil}
}
