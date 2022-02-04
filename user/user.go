package user

type UserIll struct {
	Year   int
	Month  int
	Sign   string
	Degree string
}

func MakeIll() UserIll {
	return UserIll{
		Year:   0,
		Month:  0,
		Sign:   "",
		Degree: ""}
}

type UserVac struct {
	Year   int
	Month  int
	Kind   string
	Effect string
}

func MakeVac() UserVac {
	return UserVac{
		Year:   0,
		Month:  0,
		Kind:   "",
		Effect: ""}
}

type UserData struct {
	Country   int
	Birth     int
	Gender    int
	Education string
	Origin    string
	Vaccine   string
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
		Education: "",
		Origin:    "",
		Vaccine:   "",
		CountIll:  0,
		Ill:       nil,
		CountVac:  0,
		Vac:       nil}
}
