package mocks

type DocumentProtectorMock struct {
	HashFunc     func(document string) string
	EncryptFunc  func(document string) (string, error)
	DecryptFunc  func(cipherText string) (string, error)
}

func (m *DocumentProtectorMock) Hash(document string) string {
	if m.HashFunc != nil {
		return m.HashFunc(document)
	}
	return "hash-" + document
}

func (m *DocumentProtectorMock) Encrypt(document string) (string, error) {
	if m.EncryptFunc != nil {
		return m.EncryptFunc(document)
	}
	return "encrypted-" + document, nil
}

func (m *DocumentProtectorMock) Decrypt(cipherText string) (string, error) {
	if m.DecryptFunc != nil {
		return m.DecryptFunc(cipherText)
	}
	return cipherText, nil
}
