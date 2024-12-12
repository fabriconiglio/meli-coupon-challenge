package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "log"
    "sort"
    "github.com/gorilla/mux"
)

// MeliItemResponse estructura para la respuesta de la API de ML
type MeliItemResponse struct {
    ID    string  `json:"id"`
    Price float64 `json:"price"`
    Title string  `json:"title"`
}

// CouponRequest estructura para la petición
type CouponRequest struct {
    ItemIDs []string  `json:"item_ids"`
    Amount  float64   `json:"amount"`
}

// CouponResponse estructura para la respuesta
type CouponResponse struct {
    ItemIDs []string  `json:"item_ids"`
    Total   float64   `json:"total"`
}

// ItemPrice estructura para almacenar el precio de un item
type ItemPrice struct {
    ID    string
    Price float64
}

// ItemStats estructura para las estadísticas de un item
type ItemStats struct {
    ID       string `json:"id"`
    Quantity int    `json:"quantity"`
}

// FavoriteStats maneja el conteo de items favoritos
type FavoriteStats struct {
    stats map[string]int
}

// Variable global para mantener las estadísticas
var favoriteStats = &FavoriteStats{
    stats: make(map[string]int),
}

// getItemPrice obtiene el precio de un item desde la API de ML
func getItemPrice(itemID string) (float64, error) {
    url := fmt.Sprintf("https://api.mercadolibre.com/items/%s", itemID)
    log.Printf("Consultando API de ML: %s", url)
    
    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("error haciendo GET a %s: %v", url, err)
    }
    defer resp.Body.Close()
    
    // Leemos el body completo para logging
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("error leyendo respuesta: %v", err)
    }
    
    log.Printf("Respuesta de ML para %s: %s", itemID, string(body))
    
    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("API respondió con status %d: %s", resp.StatusCode, string(body))
    }

    var item MeliItemResponse
    if err := json.Unmarshal(body, &item); err != nil {
        return 0, fmt.Errorf("error decodificando JSON: %v", err)
    }

    if item.Price == 0 {
        return 0, fmt.Errorf("precio no encontrado para item %s", itemID)
    }

    log.Printf("Precio obtenido para %s: %f", itemID, item.Price)
    return item.Price, nil
}

// findBestCombination encuentra la mejor combinación de items
func findBestCombination(items []ItemPrice, amount float64) ([]string, float64) {
    n := len(items)
    bestTotal := 0.0
    var bestComb []string
    
    log.Printf("Buscando mejor combinación para %d items con monto %f", n, amount)
    log.Printf("Items disponibles: %+v", items)

    // Ordenamos los items por precio para optimizar la búsqueda
    sort.Slice(items, func(i, j int) bool {
        return items[i].Price < items[j].Price
    })

    // Probamos todas las combinaciones posibles
    for mask := 1; mask < (1 << n); mask++ {
        currentTotal := 0.0
        var currentComb []string

        // Verificamos cada bit de la máscara
        for i := 0; i < n; i++ {
            if (mask & (1 << i)) != 0 {
                currentTotal += items[i].Price
                currentComb = append(currentComb, items[i].ID)
            }
        }

        // Si esta combinación es mejor que la anterior y no excede el monto
        if currentTotal <= amount && currentTotal > bestTotal {
            bestTotal = currentTotal
            bestComb = currentComb
            log.Printf("Nueva mejor combinación encontrada: %v con total %f", bestComb, bestTotal)
        }
    }

    return bestComb, bestTotal
}


// Método para obtener el top 5 de items
func (fs *FavoriteStats) getTopItems() []ItemStats {
    // Convertir el map a slice para poder ordenarlo
    var items []ItemStats
    for id, quantity := range fs.stats {
        items = append(items, ItemStats{
            ID:       id,
            Quantity: quantity,
        })
    }

    // Ordenar por cantidad (descendente)
    sort.Slice(items, func(i, j int) bool {
        return items[i].Quantity > items[j].Quantity
    })

    // Devolver solo los primeros 5 items
    if len(items) > 5 {
        items = items[:5]
    }
    return items
}

// handleStats maneja el endpoint de estadísticas
func handleStats(w http.ResponseWriter, r *http.Request) {
    topItems := favoriteStats.getTopItems()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(topItems)
}

// handleCoupon procesa la petición del cupón
func handleCoupon(w http.ResponseWriter, r *http.Request) {
    var req CouponRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Petición recibida: %+v", req)

	 // Registrar los items como favoritos
	 for _, id := range req.ItemIDs {
        favoriteStats.stats[id]++
    }

    // Obtenemos los precios de todos los items
    var items []ItemPrice
    validItems := false // Flag para verificar si al menos un item es válido

    for _, id := range req.ItemIDs {
        price, err := getItemPrice(id)
        if err != nil {
            log.Printf("Error obteniendo precio para item %s: %v", id, err)
            continue // Seguimos con el siguiente item en lugar de fallar
        }
        items = append(items, ItemPrice{ID: id, Price: price})
        validItems = true
    }

    if !validItems {
        http.Error(w, "No se encontraron items válidos", http.StatusBadRequest)
        return
    }

    // Encontramos la mejor combinación
    selectedItems, total := findBestCombination(items, req.Amount)

    response := CouponResponse{
        ItemIDs: selectedItems,
        Total:   total,
    }

    log.Printf("Enviando respuesta: %+v", response)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Configurar logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    router := mux.NewRouter()
    
    router.HandleFunc("/coupon/", handleCoupon).Methods("POST")
	router.HandleFunc("/coupon/stats", handleStats).Methods("GET")
    
    srv := &http.Server{
        Handler:      router,
        Addr:         ":8080",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    
    log.Printf("Server starting on port 8080...")
    log.Fatal(srv.ListenAndServe())
}