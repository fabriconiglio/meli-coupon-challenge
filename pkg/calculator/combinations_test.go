package calculator

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "meli-coupon/internal/domain"
)

func TestFindBestCombination(t *testing.T) {
    tests := []struct {
        name           string
        items          []domain.ItemPrice
        amount         float64
        expectedIDs    []string
        expectedTotal  float64
    }{
        {
            name: "ejemplo del challenge",
            items: []domain.ItemPrice{
                {ID: "MLA1", Price: 100},
                {ID: "MLA2", Price: 210},
                {ID: "MLA3", Price: 260},
                {ID: "MLA4", Price: 80},
                {ID: "MLA5", Price: 90},
            },
            amount:        500,
            expectedIDs:   []string{"MLA1", "MLA2", "MLA4", "MLA5"},
            expectedTotal: 480,
        },
        {
            name: "monto exacto",
            items: []domain.ItemPrice{
                {ID: "MLA1", Price: 100},
                {ID: "MLA2", Price: 200},
            },
            amount:        300,
            expectedIDs:   []string{"MLA1", "MLA2"},
            expectedTotal: 300,
        },
        {
            name: "sin combinación válida",
            items: []domain.ItemPrice{
                {ID: "MLA1", Price: 200},
                {ID: "MLA2", Price: 200},
            },
            amount:        100,
            expectedIDs:   []string{},
            expectedTotal: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ids, total := FindBestCombination(tt.items, tt.amount)
            assert.Equal(t, tt.expectedIDs, ids)
            assert.Equal(t, tt.expectedTotal, total)
        })
    }
}