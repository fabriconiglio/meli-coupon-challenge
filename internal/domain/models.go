package domain

// MeliItem representa un item de MercadoLibre
type MeliItem struct {
    ID    string  `json:"id"`
    Price float64 `json:"price"`
    Title string  `json:"title"`
}

// CouponRequest representa la petición del cupón
type CouponRequest struct {
    ItemIDs []string `json:"item_ids"`
    Amount  float64  `json:"amount"`
}

// CouponResponse representa la respuesta del cupón
type CouponResponse struct {
    ItemIDs []string `json:"item_ids"`
    Total   float64  `json:"total"`
}

// ItemStats representa las estadísticas de un item
type ItemStats struct {
    ID       string `json:"id"`
    Quantity int    `json:"quantity"`
}

// ItemPrice representa un item con su precio para cálculos internos
type ItemPrice struct {
    ID    string
    Price float64
}