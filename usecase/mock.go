package usecase

import (
	"consumer/models"

	"github.com/stretchr/testify/mock"
)

type UsecaseMock struct {
	mock.Mock
}

func (m *UsecaseMock) Upsert(product models.Product) (models.Product, error) {
	args := m.Called(product)

	return args.Get(0).(models.Product), args.Error(1)
}

func (f *UsecaseMock) Delete(product models.Product) (models.Product, error) {
	args := f.Called(product)

	return args.Get(0).(models.Product), args.Error(1)
}
