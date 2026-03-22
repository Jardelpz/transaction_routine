package ports

type DocumentProtector interface {
	Hash(document string) string
	Encrypt(document string) (string, error)
	Decrypt(cipherText string) (string, error)
}
