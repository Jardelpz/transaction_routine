package domain

import "unicode"

type Account struct {
	AccountId         int64
	DocumentHash      string
	DocumentEncrypted string
}

func ValidateDocument(document string) error {
	if len(document) != 11 { //cnpj too?
		return ErrDocumentInvalid
	}

	for _, d := range document {
		if !unicode.IsDigit(d) {
			return ErrDocumentNonNumber
		}
	}

	return nil
}
