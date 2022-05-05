package validation

// Todo: пока не используется, сделать реализацию
type ValidationI interface {
	IsEmailValid(email string) bool
	IsPasswordValid(email string) bool
}

type ValidationData struct {
}

func NewValidationData() *ValidationData {
	return &ValidationData{}
}

func (o ValidationData) IsEmailValid(email string) bool {
	return false
}

func (o ValidationData) IsPasswordValid(email string) bool {
	return false
}
