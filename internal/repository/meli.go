package repository

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
    "meli-coupon/internal/domain"
    "meli-coupon/internal/cache"
    "meli-coupon/internal/monitoring"
)

type MeliRepository interface {
    GetItemPrice(itemID string) (float64, error)
}

type meliRepository struct {
    baseURL string
    cache   cache.Cache
}

func NewMeliRepository(cache cache.Cache) MeliRepository {
    return &meliRepository{
        baseURL: "https://api.mercadolibre.com",
        cache:   cache,
    }
}

func (r *meliRepository) GetItemPrice(itemID string) (float64, error) {
    // Intentar obtener del caché
    if price, found := r.cache.Get(itemID); found {
        monitoring.RecordCacheHit()
        log.Printf("Cache hit for item %s: %f", itemID, price)
        return price, nil
    }

    monitoring.RecordCacheMiss()

    // Si no está en caché, consultar la API
    url := fmt.Sprintf("%s/items/%s", r.baseURL, itemID)
    log.Printf("Cache miss, consulting API: %s", url)
    
    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("error haciendo GET a %s: %v", url, err)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("error leyendo respuesta: %v", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("API respondió con status %d: %s", resp.StatusCode, string(body))
    }

    var item domain.MeliItem
    if err := json.Unmarshal(body, &item); err != nil {
        return 0, fmt.Errorf("error decodificando JSON: %v", err)
    }

    if item.Price == 0 {
        return 0, fmt.Errorf("precio no encontrado para item %s", itemID)
    }

    // Guardar en caché
    r.cache.Set(itemID, item.Price)
    log.Printf("Cached price for item %s: %f", itemID, item.Price)

    return item.Price, nil
}