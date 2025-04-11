package enums

type ErrorType string

const (
	ErrInvalidToken            ErrorType = "invalid or expired token"
	ErrWrongCredentials        ErrorType = "wrong username or password"
	ErrNoAuthToken             ErrorType = "no auth token provided"
	ErrUserNotFound            ErrorType = "user not found"
	ErrPVZNotFound             ErrorType = "pvz not found"
	ErrWrongTokenFormat        ErrorType = "wrong token format"
	ErrInvalidCity             ErrorType = "invalid city"
	ErrNoProductsToDelete      ErrorType = "no products to delete"
	ErrNoOpenReceptionToDelete ErrorType = "no open reception to delete product"
	ErrNoEmployeeRights        ErrorType = "only employees can do that"
	ErrNoModeratorRights       ErrorType = "only moderators can do that"
	ErrNoOpenReceptionsToAdd   ErrorType = "no open reception to add product"
	ErrInvalidRole             ErrorType = "invalid role"
	ErrOpenReception           ErrorType = "there is already an open reception"
	ErrNoOpenReceptionToClose  ErrorType = "no open reception to close"
	ErrInvalidStartDate        ErrorType = "invalid startDate"
	ErrInvalidEndDate          ErrorType = "invalid endDate"
)

func (et ErrorType) Error() string {
	return string(et)
}
