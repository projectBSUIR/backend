package types

type UserStatus int

const (
	TestMachine  UserStatus = 1
	Participant  UserStatus = 2
	Coach        UserStatus = 3
	Admin        UserStatus = 4
	UnAuthorized UserStatus = 0
)

func GetStatusById(statusId int) UserStatus {
	switch statusId {
	case 1:
		return TestMachine
	case 2:
		return Participant
	case 3:
		return Coach
	case 4:
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
	case "TestMachine":
		return TestMachine
	default:
		return UnAuthorized
	}
}
