package types

type UserStatus int

const (
	Participant  UserStatus = 1
	Coach        UserStatus = 2
	Admin        UserStatus = 3
	UnAuthorized UserStatus = 0
)

func GetStatusById(statusId int) UserStatus {
	switch statusId {
	case 1:
		return Participant
	case 2:
		return Coach
	case 3:
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
