package indexbuilding

// city转大写
type User struct {
	Type      string
	City      string
	Lat       float64
	Lng       float64
	HourStart int
	MinStart  int
}

func (u *User) TypeEncode() ([]string, error) {
	return TypeEncoding(u.Type)
}

func (u *User) LocationEncode(x int) ([]string, error) {
	locationCode, err := LocationEncodingUser(u.City, x, u.Lat, u.Lng)
	if err != nil {
		return nil, err
	}
	AddCityIndex(u.City, locationCode)
	return locationCode, nil
}

func (u *User) TimeEncode() ([]string, error) {
	return TimePointEncoding(u.HourStart, u.MinStart)
}
