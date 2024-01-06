package province

type Province struct {
	ID             int
	Name           string
	SelectDistrict func(province_id int) string
}
