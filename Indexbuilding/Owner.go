package indexbuilding

// city转大写
type Owner struct {
	Type      string  `json:"type"`
	City      string  `json:"city"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	HourStart int     `json:"hourStart"`
	MinStart  int     `json:"minStart"`
	HourClose int     `json:"hourClose"`
	MinClose  int     `json:"minClose"`
}

func (o *Owner) TypeEncode() ([]string, error) {
	return TypeEncoding(o.Type)
}
func (o *Owner) TypeComplementEncode() ([]string, error) {
	return TypeEncodingComplement(o.Type)
}

func (o *Owner) LocationEncode() ([]string, error) {
	locationCode, err := LocationEncoding(o.City, o.Lat, o.Lng)
	if err != nil {
		return nil, err
	}
	AddCityIndex(o.City, locationCode)
	return locationCode, nil
}
func (o *Owner) LocationComplementEncode() ([]string, error) {
	locationCode, err := LocationEncodingComplement(o.City, o.Lat, o.Lng)
	if err != nil {
		return nil, err
	}
	AddCityIndex(o.City, locationCode)
	return locationCode, nil
}

func (o *Owner) TimeEncode() ([]string, error) {
	return TimeRangeEncoding(o.HourStart, o.MinStart, o.HourClose, o.HourClose)
}
func (o *Owner) TimeComplementEncode() ([]string, error) {
	return TimeRangeEncodingComplement(o.HourStart, o.MinStart, o.HourClose, o.HourClose)
}
