package calculator

import (
    "log"
    "sort"
    "meli-coupon/internal/domain"
)

// FindBestCombination encuentra la mejor combinación de items que no exceda el monto
func FindBestCombination(items []domain.ItemPrice, amount float64) ([]string, float64) {
    n := len(items)
    bestTotal := 0.0
    var bestComb []string
    
    log.Printf("Buscando mejor combinación para %d items con monto %f", n, amount)

    // Ordenamos los items por precio para optimizar la búsqueda
    sort.Slice(items, func(i, j int) bool {
        return items[i].Price < items[j].Price
    })

    // Probamos todas las combinaciones posibles
    for mask := 1; mask < (1 << n); mask++ {
        currentTotal := 0.0
        var currentComb []string

        for i := 0; i < n; i++ {
            if (mask & (1 << i)) != 0 {
                currentTotal += items[i].Price
                currentComb = append(currentComb, items[i].ID)
            }
        }

        if currentTotal <= amount && currentTotal > bestTotal {
            bestTotal = currentTotal
            bestComb = currentComb
            log.Printf("Nueva mejor combinación encontrada: %v con total %f", bestComb, bestTotal)
        }
    }

    return bestComb, bestTotal
}