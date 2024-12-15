package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "meli-coupon/internal/domain"
)

type mockMeliRepository struct {
    mock.Mock
}

func (m *mockMeliRepository) GetItemPrice(itemID string) (float64, error) {
    args := m.Called(itemID)
    return args.Get(0).(float64), args.Error(1)
}

func TestProcessCoupon(t *testing.T) {
    // Crear el mock y el servicio
    mockRepo := new(mockMeliRepository)
    service := NewCouponService(mockRepo)

    // Caso 1: El ejemplo del challenge
    t.Run("ejemplo del challenge", func(t *testing.T) {
        // Configurar los precios del ejemplo
        mockRepo.On("GetItemPrice", "MLA1").Return(100.0, nil).Once()
        mockRepo.On("GetItemPrice", "MLA2").Return(210.0, nil).Once()
        mockRepo.On("GetItemPrice", "MLA3").Return(260.0, nil).Once()

        req := domain.CouponRequest{
            ItemIDs: []string{"MLA1", "MLA2", "MLA3"},
            Amount:  350.0, // Ajustado para que tome MLA1 + MLA2
        }

        response, err := service.ProcessCoupon(req)
        assert.NoError(t, err)
        assert.Equal(t, []string{"MLA1", "MLA2"}, response.ItemIDs)
        assert.Equal(t, 310.0, response.Total)
    })

    // Caso 2: Probar con límite mayor
    t.Run("monto mayor", func(t *testing.T) {
        mockRepo.On("GetItemPrice", "MLA1").Return(100.0, nil).Once()
        mockRepo.On("GetItemPrice", "MLA2").Return(210.0, nil).Once()
        mockRepo.On("GetItemPrice", "MLA3").Return(260.0, nil).Once()

        req := domain.CouponRequest{
            ItemIDs: []string{"MLA1", "MLA2", "MLA3"},
            Amount:  500.0,
        }

        response, err := service.ProcessCoupon(req)
        assert.NoError(t, err)
        assert.Equal(t, []string{"MLA2", "MLA3"}, response.ItemIDs)
        assert.Equal(t, 470.0, response.Total)
    })
}

func TestGetTopFavorites(t *testing.T) {
    mockRepo := new(mockMeliRepository)
    service := NewCouponService(mockRepo)

    mockRepo.On("GetItemPrice", mock.Anything).Return(100.0, nil)

    // Primera petición
    req1 := domain.CouponRequest{
        ItemIDs: []string{"MLA1", "MLA2"},
        Amount:  500,
    }
    service.ProcessCoupon(req1)

    // Segunda petición
    req2 := domain.CouponRequest{
        ItemIDs: []string{"MLA1", "MLA3"},
        Amount:  500,
    }
    service.ProcessCoupon(req2)

    stats := service.GetTopFavorites()
    
    // Verificar que MLA1 es el más favoriteado
    assert.GreaterOrEqual(t, len(stats), 1, "Debería haber al menos un item en las estadísticas")
    if len(stats) > 0 {
        assert.Equal(t, "MLA1", stats[0].ID)
        assert.Equal(t, 2, stats[0].Quantity)
    }
}