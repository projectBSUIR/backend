package types

type UserStatus int

const (
	Participant  UserStatus = 0
	Coach        UserStatus = 1
	Admin        UserStatus = 2
	UnAuthorized UserStatus = 3
)

func GetStatusById(statusId int) UserStatus {
	switch statusId {
	case 0:
		return Participant
	case 1:
		return Coach
	case 2:
		return Admin
	}
	return UnAuthorized
}

func GetStatusByString(statusString string) UserStatus {
	switch statusString {
	case "Participant":
		return Participant
	case "Coach":
		return Coach
	case "Admin":
		return Admin
	default:
		return UnAuthorized
	}
}
