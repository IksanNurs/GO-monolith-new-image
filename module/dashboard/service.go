package dashboard

type DashboardService interface {
	CountUser() (int, error)
	CountInstitution() (int, error)
	CountRequest() (int, error)
}

type service struct {
	repository DashboardRepository
}

func NewService(repository DashboardRepository) *service {
	return &service{repository}
}

func (s *service) CountUser() (int, error) {
	return s.repository.CountUser()
}

func (s *service) CountInstitution() (int, error) {
	return s.repository.CountInstitution()
}

func (s *service) CountRequest() (int, error) {
	return s.repository.CountRequest()
}
