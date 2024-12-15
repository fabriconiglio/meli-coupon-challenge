package services

import (
    "sort"
    "sync"
    "meli-coupon/internal/domain"
    "meli-coupon/internal/repository"
    "meli-coupon/pkg/calculator"
)

// CouponService interfaz para el servicio de cupones
type CouponService interface {
    ProcessCoupon(req domain.CouponRequest) (domain.CouponResponse, error)
    GetTopFavorites() []domain.ItemStats
}

type couponService struct {
    meliRepo repository.MeliRepository
    stats    *statsManager
}

// statsManager maneja las estadísticas de items favoritos
type statsManager struct {
    mu    sync.RWMutex
    stats map[string]int
}

// NewCouponService crea una nueva instancia del servicio
func NewCouponService(meliRepo repository.MeliRepository) CouponService {
    return &couponService{
        meliRepo: meliRepo,
        stats:    &statsManager{
            stats: make(map[string]int),
        },
    }
}

// ProcessCoupon procesa la petición del cupón
func (s *couponService) ProcessCoupon(req domain.CouponRequest) (domain.CouponResponse, error) {
    // Registrar favoritos
    s.stats.registerFavorites(req.ItemIDs)

    // Obtener precios
    var items []domain.ItemPrice
    for _, id := range req.ItemIDs {
        price, err := s.meliRepo.GetItemPrice(id)
        if err != nil {
            continue
        }
        items = append(items, domain.ItemPrice{
            ID:    id,
            Price: price,
        })
    }

    if len(items) == 0 {
        return domain.CouponResponse{}, domain.ErrNoValidItems
    }

    // Encontrar mejor combinación y ordenar el resultado
    selectedItems, total := calculator.FindBestCombination(items, req.Amount)
    sort.Strings(selectedItems)

    return domain.CouponResponse{
        ItemIDs: selectedItems,
        Total:   total,
    }, nil
}

// GetTopFavorites obtiene los 5 items más favoritados
func (s *couponService) GetTopFavorites() []domain.ItemStats {
    return s.stats.getTopItems(5)
}

// registerFavorites registra items como favoritos
func (sm *statsManager) registerFavorites(itemIDs []string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    for _, id := range itemIDs {
        sm.stats[id]++
    }
}

// getTopItems obtiene los n items más favoritados
func (sm *statsManager) getTopItems(n int) []domain.ItemStats {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    var items []domain.ItemStats
    for id, quantity := range sm.stats {
        items = append(items, domain.ItemStats{
            ID:       id,
            Quantity: quantity,
        })
    }

    // Ordenar por cantidad (descendente)
    sort.Slice(items, func(i, j int) bool {
        return items[i].Quantity > items[j].Quantity
    })

    // Retornar solo los primeros n items
    if len(items) > n {
        items = items[:n]
    }
    return items
}