package models

import (
	"fiber-apis/databases"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Contest struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	Duration  int64  `json:"duration"`
}

type ContestInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (contest *Contest) NotStarted() bool {
	timeNow := time.Now()
	startTimeContest, _ := time.Parse(time.RFC3339, contest.StartTime)
	return timeNow.Before(startTimeContest)
}

func (contest *Contest) Create(c *fiber.Ctx) (ContestInfo, error) {
	row, err := databases.DataBase.Exec("INSERT INTO `contest` (`contest_name`, `start_time`, `duration`) VALUES (?, ?, ?)",
		contest.Name, contest.StartTime, contest.Duration)
	if err != nil {
		prevErr := err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return ContestInfo{}, err
		}
		return ContestInfo{}, prevErr
	}
	id, err := row.LastInsertId()
	if err != nil {
		prevErr := err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return ContestInfo{}, err
		}
		return ContestInfo{}, prevErr
	}
	contest.Id = id
	err = SetAuthorOfContest(contest.Id, c)
	if err != nil {
		prevErr := err
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return ContestInfo{}, err
		}
		return ContestInfo{}, prevErr
	}
	return ContestInfo{Id: contest.Id, Name: contest.Name}, nil
}

func FetchAllContests() ([]Contest, error) {
	rows, err := databases.DataBase.Query("SELECT * FROM `contest`")
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	contests := make([]Contest, 0)
	for rows.Next() {
		var contest Contest
		rows.Scan(&contest.Id, &contest.Name, &contest.StartTime, &contest.Duration)
		contests = append(contests, contest)
	}
	return contests, nil
}

func (contest *Contest) FetchContest(contestId int) error {
	row, err := databases.DataBase.Query("SELECT `contest_name`, `start_time`, `duration` FROM `contest` WHERE id = ?", contestId)
	if err != nil {
		return err
	}

	row.Next()
	err = row.Scan(&contest.Name, &contest.StartTime, &contest.Duration)
	if err != nil {
		return err
	}

	return nil
}

func FetchAllContestsForAuthor(userId int64) ([]ContestInfo, error) {
	rows, err := databases.DataBase.Query("SELECT `id`, `name` FROM `contest` WHERE `id` IN (SELECT `contest_id` FROM `contestAuthor` WHERE `user_id`=?)", userId)
	if err != nil {
		return nil, err
	}
	var contests []ContestInfo
	for rows.Next() {
		var contestInfo ContestInfo
		err = rows.Scan(&contestInfo.Id, &contestInfo.Name)
		if err != nil {
			return nil, err
		}
		contests = append(contests, contestInfo)
	}
	return contests, nil
}

func GetParticipantsIds(contestId int64) ([]int64, error) {
	var participantsIds []int64
	res, err := databases.DataBase.Query("SELECT `user_id` FROM `contestResult` WHERE `contest_id`= ?", contestId)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var id int64
		err := res.Scan(&id)
		if err != nil {
			return nil, err
		}
		participantsIds = append(participantsIds, id)
	}
	return participantsIds, nil
}
