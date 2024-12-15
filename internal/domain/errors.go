package domain

import "errors"

var (
    // ErrNoValidItems cuando no hay items válidos para procesar
    ErrNoValidItems = errors.New("no se encontraron items válidos")

    // ErrInvalidAmount cuando el monto del cupón es inválido
    ErrInvalidAmount = errors.New("monto del cupón inválido")

    // ErrMeliAPIError cuando hay un error en la API de MercadoLibre
    ErrMeliAPIError = errors.New("error en la API de MercadoLibre")

    // ErrInvalidItemID cuando el ID del item es inválido
    ErrInvalidItemID = errors.New("ID de item inválido")
)