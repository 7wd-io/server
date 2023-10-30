package password

import "golang.org/x/crypto/bcrypt"

func New() *Manager {
	return new(Manager)
}

type Manager struct{}

func (dst *Manager) Hash(password string, cost int) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(h), err
}

func (dst *Manager) Check(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
