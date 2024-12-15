package repository

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
    "meli-coupon/internal/domain"
)

// MeliRepository interface para interactuar con la API de ML
type MeliRepository interface {
    GetItemPrice(itemID string) (float64, error)
}

type meliRepository struct {
    baseURL string
}

// NewMeliRepository crea una nueva instancia del repositorio
func NewMeliRepository() MeliRepository {
    return &meliRepository{
        baseURL: "https://api.mercadolibre.com",
    }
}

// GetItemPrice obtiene el precio de un item desde la API de ML
func (r *meliRepository) GetItemPrice(itemID string) (float64, error) {
    url := fmt.Sprintf("%s/items/%s", r.baseURL, itemID)
    log.Printf("Consultando API de ML: %s", url)
    
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

    log.Printf("Precio obtenido para %s: %f", itemID, item.Price)
    return item.Price, nil
}