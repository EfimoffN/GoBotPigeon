package sqlapi

import (
	"context"
	"strings"
	"time"

	apitypes "GoBotPigeon/types/apitypes"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CheckingPigeonWork ...
func (api *API) CheckingPigeonWork(userN string) (bool, error) {

	user, err := api.GetUserByName(userN)
	if err != nil {
		return false, errors.Wrap(err, "Get user by name failed")
	}

	var botWork = []apitypes.BotWork{}
	err = api.db.Select(&botWork, "SELECT * FROM prj_botwork WHERE userid = $1 LIMIT 1;", user.UserID)
	if err != nil {
		return false, errors.Wrap(err, "SELECT * FROM prj_botwork failed")
	}

	return botWork[0].BotWorkFlag, err
}

// StartPigeonWork ...
func (api *API) StartPigeonWork(userN string) error {

	user, err := getUserByName(context.Background(), api.db, userN)
	if err != nil {
		return errors.Wrap(err, "Get user by name failed")
	}

	work := func(ctx context.Context, db TxContext) error {
		var botWork = []apitypes.BotWork{}
		err = db.SelectContext(ctx, &botWork, "SELECT * FROM prj_botwork WHERE userid = $1 LIMIT 1;", user.UserID)
		if err != nil {
			return errors.Wrap(err, "SELECT * FROM prj_botwork failed")
		}

		if len(botWork) == 0 {
			err = api.CreatePigeonWorkFlag(userN)
			if err != nil {
				return errors.Wrap(err, "Create pigeon work flag failed")
			}
		} else {
			if _, err := db.ExecContext(ctx, `UPDATE prj_botwork SET botworkflag = $1  WHERE botworkid = $2`, true, botWork[0].BotWorkID); err != nil {
				return errors.Wrap(err, "UPDATE prj_botwork SET failed")
			}
		}

		return nil
	}

	if err := RunInTransaction(context.Background(), api.db, work); err != nil {
		return errors.Wrap(err, "RunInTransaction failed")
	}

	return nil
}

// StopPigeonWork ...
func (api *API) StopPigeonWork(userN string) error {
	user, err := api.GetUserByName(userN)
	if err != nil {
		return errors.Wrap(err, "Get user by name failed")
	}

	work := func(ctx context.Context, db TxContext) error {

		var botWork = []apitypes.BotWork{}
		err = db.SelectContext(ctx, &botWork, "SELECT * FROM prj_botwork WHERE userid = $1 LIMIT 1", user.UserID)
		if err != nil {
			return errors.Wrap(err, "SELECT * FROM prj_botwork failed")
		}

		if _, err := db.ExecContext(ctx, `UPDATE prj_botwork SET botworkflag = $1 WHERE botworkid = $2`, false, botWork[0].BotWorkID); err != nil {
			return errors.Wrap(err, "UPDATE prj_botwork failed")
		}
		return nil
	}

	if err := RunInTransaction(context.Background(), api.db, work); err != nil {
		return err
	}

	return nil
}

// CreatePigeonWorkFlag ...
func (api *API) CreatePigeonWorkFlag(userN string) error {
	user, err := api.GetUserByName(userN)
	if err != nil {
		return errors.Wrap(err, "Get user by name failed")
	}

	work := func(ctx context.Context, db TxContext) error {
		var botWork = []apitypes.BotWork{}
		err = db.SelectContext(ctx, &botWork, "SELECT * FROM prj_botwork WHERE userid = $1", user.UserID)

		if err != nil {
			return errors.Wrap(err, "SELECT * FROM prj_botwork failed")
		}

		if len(botWork) == 0 {
			uuidWithHyphen := uuid.New()
			uid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

			if _, err := db.ExecContext(ctx, `INSERT INTO prj_botwork ("botworkid", "userid", "botworkflag") VALUES ($1, $2, $3)`, uid, user.UserID, true); err != nil {
				return err
			}
		}

		return nil
	}

	if err := RunInTransaction(context.Background(), api.db, work); err != nil {
		return err
	}

	return nil
}

// SetLastComandUser ...
func (api *API) SetLastComandUser(userN string, command string) error {
	user, err := api.GetUserByName(userN)
	if err != nil {
		return errors.Wrap(err, "Get user by name failed")
	}

	work := func(ctx context.Context, db TxContext) error {
		uuidWithHyphen := uuid.New()
		uid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
		today := time.Now()
		tTime := today.Add(10 * time.Minute).Format("2006/1/2 15:04")
		if _, err := db.ExecContext(ctx, `INSERT INTO prj_lastusercommand ("commandid", "userid", "command", "datacommand") VALUES ($1, $2, $3, $4)`, uid, user.UserID, command, tTime); err != nil {
			return err
		}
		return nil
	}
	if err := RunInTransaction(context.Background(), api.db, work); err != nil {
		return err
	}
	return nil
}

// GetLastCommandByUserName ...
func (api *API) GetLastCommandByUserName(userN string) (*apitypes.LastUserCommand, error) {
	return getLastCommandByUserName(context.Background(), api.db, userN)
}

func getLastCommandByUserName(ctx context.Context, db TxContext, userN string) (*apitypes.LastUserCommand, error) {
	user, err := getUserByName(context.Background(), db, userN) // использовать context.Background() или ctx ?

	if err != nil {
		return nil, errors.Wrap(err, "Get user by name failed")
	}

	if user == nil {
		return nil, nil // создать ошибку
	}

	var arrCommand = []apitypes.LastUserCommand{}

	err = db.SelectContext(ctx, &arrCommand, "SELECT * FROM prj_lastusercommand WHERE (userid = $1) ORDER BY datacommand DESC", user.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "SELECT * FROM prj_lastusercommand failed")
	}

	if len(arrCommand) == 0 {
		return nil, nil // создать ошибку
	}

	return &arrCommand[0], nil
}

// DeleteLastCommand ...
func (api *API) DeleteLastCommand(userN string, command string) error {
	user, err := api.GetUserByName(userN)
	if err != nil {
		return errors.Wrap(err, "GetUserByName failed")
	}

	work := func(ctx context.Context, db TxContext) error {
		if _, err := db.ExecContext(ctx, `DELETE FROM prj_lastusercommand WHERE userid = $1`, user.UserID); err != nil {
			return errors.Wrap(err, "DELETE FROM prj_lastusercommand")
		}

		return nil
	}

	if err := RunInTransaction(context.Background(), api.db, work); err != nil {
		return err
	}
	return nil
}
