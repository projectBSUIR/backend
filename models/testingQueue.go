package models

import "fiber-apis/databases"

type TestingQueue struct {
	Id           int64 `json:"id"`
	SubmissionId int64 `json:"submission_id"`
}

func (testingQueue *TestingQueue) AddSubmissionToQueue() error {
	row, err := databases.DataBase.Exec("INSERT INTO `testingQueue` (`submission_id`) VALUES (?)", testingQueue.SubmissionId)
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return err
	}
	id, err := row.LastInsertId()
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return err
	}
	testingQueue.Id = id
	return nil
}

func DeleteSubmissionFromTestingQueue(submissionId int64) error {
	_, err := databases.DataBase.Exec("DELETE FROM `testingQueue` WHERE `submission_id`=?", submissionId)
	if err != nil {
		_, err := databases.DataBase.Query("ROLLBACK")
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func GetFirstSubmissionFromTestingQueue() (TestingQueue, error) {
	row, err := databases.DataBase.Query("SELECT `id`, `submission_id` FROM `testingQueue` WHERE `id`=(SELECT MIN(`id`) FROM `testingQueue`)")
	if err != nil {
		return TestingQueue{}, err
	}
	var testingQueue TestingQueue
	row.Next()
	err = row.Scan(&testingQueue.Id, &testingQueue.SubmissionId)
	if err != nil {
		return TestingQueue{}, err
	}
	return testingQueue, nil
}
