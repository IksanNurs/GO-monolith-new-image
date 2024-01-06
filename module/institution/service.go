package institution

import (
	"akuntansi/module/province"
	"errors"
)

type Service interface {
	FetchAllInstitution() ([]Institution, error)
	FetchAllProvince() ([]province.Province, error)
	Create(input FormCreateInstitution) (Institution, error)
	Update(input FormCreateInstitution, ID int) (Institution, error)
	Delete(ID int) (Institution, error)
	GetInstitutionByID(ID int) (Institution, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Create(input FormCreateInstitution) (Institution, error) {
	institution := Institution{}
	institution.Name.String = input.Name
	institution.Type.Int64 = input.Type
	institution.Code.String = input.Code
	institution.Province_id.String = input.Province_ID
	institution.Province_name.String = input.Province_Name
	NewInstitution, err := s.repository.Save(institution)
	if err != nil {
		return NewInstitution, err
	}

	return NewInstitution, nil
}

func (s *service) Update(input FormCreateInstitution, ID int) (Institution, error) {
	institution := Institution{}
	institution.ID = ID
	institution.Name.String = input.Name
	institution.Type.Int64 = input.Type
	institution.Code.String = input.Code
	institution.Province_id.String = input.Province_ID
	institution.Province_name.String = input.Province_Name
	NewInstitution, err := s.repository.Update(institution)
	if err != nil {
		return NewInstitution, err
	}

	return NewInstitution, nil
}

func (s *service) FetchAllInstitution() ([]Institution, error) {
	institutions, err := s.repository.FetchAllInstitution()
	if err != nil {
		return institutions, err
	}

	return institutions, err
}

func (s *service) FetchAllProvince() ([]province.Province, error) {
	pprovinces, err := s.repository.FetchAllProvince()
	if err != nil {
		return pprovinces, err
	}

	return pprovinces, err
}

func (s *service) Delete(ID int) (Institution, error) {
	institution, err := s.repository.Delete(ID)
	if err != nil {
		return institution, err
	}

	return institution, nil
}

func (s *service) GetInstitutionByID(ID int) (Institution, error) {
	institution, err := s.repository.FindByID(ID)
	if err != nil {
		return institution, err
	}

	if institution.ID == 0 {
		return institution, errors.New("no institution found on with that ID")
	}

	return institution, nil
}
