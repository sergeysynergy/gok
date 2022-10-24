package cli

import (
	"fmt"
	"github.com/google/uuid"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

// Method pass provide local work with PASS record type.
func (c *CLI) pass() {
	if len(c.args) == 1 {
		err := c.passLs()
		if err != nil {
			fmt.Println("\nText ls failed:", err)
		}
		return
	}

	switch c.args[1] {
	case "add":
		err := c.passAdd()
		if err != nil {
			fmt.Println("\nPass add failed:", err)
		} else {
			fmt.Println("\nSuccessfully added new pass record.")
		}
	case "set":
		err := c.passSet()
		if err != nil {
			fmt.Println("\nPass set failed:", err)
		} else {
			fmt.Println("\nSuccessfully updated pass record.")
		}
	case "ls":
		err := c.passLs()
		if err != nil {
			fmt.Println("\nText ls failed:", err)
		}
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) passAdd() (err error) {
	if len(c.args) < 5 {
		return fmt.Errorf("insufficient arguments for `pass add [description] [login] [password]`")
	}

	desc := c.args[2]
	login := c.args[3]
	password := c.args[4]
	rec := &entity.Record{
		Description: entity.StringField(desc),
		Extension: &entity.Pass{
			Login:    entity.StringField(login),
			Password: entity.StringField(password),
		},
	}
	err = c.uc.RecordAdd(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to add new pass record - %s", err.Error())
	}

	return nil
}

func (c *CLI) passSet() (err error) {
	if len(c.args) < 6 {
		return fmt.Errorf("insufficient arguments for `pass set [ID] [description] [login] [password]`")
	}

	id := c.args[2]
	_, err = uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid record ID - %s", err.Error())
	}

	desc := c.args[3]
	login := c.args[4]
	password := c.args[5]
	rec := &entity.Record{
		ID:          entity.RecordID(id),
		Description: entity.StringField(desc),
		Extension: &entity.Pass{
			Login:    entity.StringField(login),
			Password: entity.StringField(password),
		},
	}

	err = c.uc.RecordSet(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to update pass record - %s", err.Error())
	}

	return nil
}

func (c *CLI) passLs() (err error) {
	if len(c.args) > 2 {
		return fmt.Errorf("invalid arguments for `pass ls`")
	}

	list, err := c.uc.RecordList(c.cfg, gokConsts.PASS)
	if err != nil {
		return fmt.Errorf("failed to get records list - %s", err.Error())
	}

	if len(list) == 0 {
		fmt.Println("No record found")
		return
	}

	fmt.Printf("\n")
	for _, r := range list {
		// Decrypt fields for output.

		desc, errDec := r.Description.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt description - %s", errDec.Error())
		}

		ex := r.Extension.(*entity.Pass)

		login, errDec := ex.Login.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt login - %s", errDec.Error())
		}

		password, errDec := ex.Password.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt password - %s", errDec.Error())
		}

		fmt.Printf("%s\t %s\t %d\t %s\t %s\t %s\n", r.ID, r.Type, r.Head, *desc, *login, *password)
	}

	return nil
}
