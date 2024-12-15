package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "meli-coupon/internal/domain"
    "meli-coupon/internal/repository"
)

// Mock del repositorio
type mockMeliRepository struct {
    mock.Mock
}

func (m *mockMeliRepository) GetItemPrice(itemID string) (float64, error) {
    args := m.Called(itemID)
    return args.Get(0).(float64), args.Error(1)
}

func TestProcessCoupon(t *testing.T) {
    mockRepo := new(mockMeliRepository)
    service := NewCouponService(mockRepo)

    // Configurar el comportamiento del mock
    mockRepo.On("GetItemPrice", "MLA1").Return(100.0, nil)
    mockRepo.On("GetItemPrice", "MLA2").Return(210.0, nil)
    mockRepo.On("GetItemPrice", "MLA3").Return(260.0, nil)

    req := domain.CouponRequest{
        ItemIDs: []string{"MLA1", "MLA2", "MLA3"},
        Amount:  500,
    }

    response, err := service.ProcessCoupon(req)

    assert.NoError(t, err)
    assert.Equal(t, []string{"MLA1", "MLA2"}, response.ItemIDs)
    assert.Equal(t, 310.0, response.Total)
}

func TestGetTopFavorites(t *testing.T) {
    mockRepo := new(mockMeliRepository)
    service := NewCouponService(mockRepo)

    // Simular algunas peticiones
    req1 := domain.CouponRequest{
        ItemIDs: []string{"MLA1", "MLA2"},
        Amount:  500,
    }
    req2 := domain.CouponRequest{
        ItemIDs: []string{"MLA1", "MLA3"},
        Amount:  500,
    }

    mockRepo.On("GetItemPrice", mock.Anything).Return(100.0, nil)

    service.ProcessCoupon(req1)
    service.ProcessCoupon(req2)

    stats := service.GetTopFavorites()

    assert.Equal(t, "MLA1", stats[0].ID)
    assert.Equal(t, 2, stats[0].Quantity)
}