package service

import (
	"cat-agency/internal/cats/models"
	"cat-agency/internal/cats/repository"
	"cat-agency/internal/common"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
)

type CatsService interface {
	CreateCat(int, string, float64) (models.Cat, error)
	UpdateCat(int, float64) error
	DeleteCat(int) error
	GetCats() ([]models.Cat, error)
	GetCat(int) (models.Cat, error)
}

type service struct {
	repo repository.CatsRepository
}

func New(repo repository.CatsRepository) CatsService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateCat(experience int, breed string, salary float64) (models.Cat, error) {
	if err := validateCatBreed(breed); err != nil {
		return models.Cat{}, err
	}
	cat, err := s.repo.CreateCat(experience, breed, salary)
	if err != nil {
		return models.Cat{}, err
	}

	return cat, nil
}

var validBreeds map[string]bool

func validateCatBreed(breed string) error {
	if validBreeds == nil {
		getBreeds()
		fmt.Println(validBreeds)
	}

	if exist := validBreeds[breed]; !exist {
		return common.NewInvalidBreedError(breed)
	}
	return nil
}

func getBreeds() error {
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return common.NewHTTPRequestError(err)
	}

	type Breed struct {
		Name string `json:"name"`
	}
	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return common.NewJSONError(err)
	}

	validBreeds = make(map[string]bool, len(breeds))
	for _, breed := range breeds {
		validBreeds[breed.Name] = true
	}
	return nil
}

func (s *service) UpdateCat(catID int, salary float64) error {
	return s.repo.UpdateCat(catID, salary)
}

func (s *service) DeleteCat(id int) error {
	return s.repo.DeleteCat(id)
}

func (s *service) GetCat(id int) (models.Cat, error) {
	return s.repo.GetCat(id)
}

func (s *service) GetCats() ([]models.Cat, error) {
	cats, err := s.repo.GetCats()
	if errors.As(err, &common.NoRowsError{}){
		return make([]models.Cat, 0), nil
	}
	return cats, nil
}
