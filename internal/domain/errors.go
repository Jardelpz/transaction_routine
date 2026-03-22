package domain

import "errors"

// account
var ErrAccountAlreadyExists = errors.New("account with document already exists")
var ErrAccountNotFound = errors.New("account not found")
var ErrDocumentInvalid = errors.New("document should contain 11 digits")
var ErrDocumentNonNumber = errors.New("document should contain only numbers")
var ErrInvalidAccountType = errors.New("account should be only numbers")

// transaction
var ErrInvalidAmount = errors.New("invalid amount")
var ErrOperationTypeNotFound = errors.New("operationType not found")
