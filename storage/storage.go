package storage

type TokenStatus int

const (
	DoesntExist TokenStatus = iota
	Expired
	Valid
)

//expired, doesnt exist
func CheckToken() TokenStatus {
	return DoesntExist
}
