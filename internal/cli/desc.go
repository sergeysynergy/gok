package cli

import (
	"fmt"
	"github.com/google/uuid"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

// Method desc provide local work with DESC record type.
func (c *CLI) desc() {
	if len(c.args) == 1 {
		c.descLs()
		return
	}

	switch c.args[1] {
	case "add":
		err := c.descAdd()
		if err != nil {
			fmt.Println("\nDesc add failed:", err)
		} else {
			fmt.Println("\nSuccessfully added new description record")
		}
	case "set":
		err := c.descSet()
		if err != nil {
			fmt.Println("\nDesc set failed:", err)
		} else {
			fmt.Println("\nSuccessfully updated description record.")
		}
	case "ls":
		err := c.descLs()
		if err != nil {
			fmt.Println("\nDesc ls failed:", err)
		}
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) descAdd() (err error) {
	if len(c.args) < 3 {
		return fmt.Errorf("insufficient arguments for `desc add [description]`")
	}

	// concatenate all further arguments in one description
	var description string
	for k, v := range c.args[2:] {
		description += v
		if len(c.args[2:]) != k+1 {
			description += " "
		}
	}

	rec := &entity.Record{
		Description: entity.StringField(description),
		Type:        gokConsts.DESC,
	}
	err = c.uc.RecordAdd(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to add new description record - %s", err.Error())
	}

	return nil
}

func (c *CLI) descSet() (err error) {
	if len(c.args) < 4 {
		return fmt.Errorf("insufficient arguments for `desc set [ID] [description]`")
	}

	id := c.args[2]
	_, err = uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid record ID - %s", err.Error())
	}

	// concatenate all further arguments in one description
	var description string
	for k, v := range c.args[3:] {
		description += v
		if len(c.args[2:]) != k+1 {
			description += " "
		}
	}

	rec := &entity.Record{
		ID:          entity.RecordID(id),
		Description: entity.StringField(description),
	}

	err = c.uc.RecordSet(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to update description record - %s", err.Error())
	}

	return nil
}

func (c *CLI) descLs() (err error) {
	if len(c.args) > 2 {
		return fmt.Errorf("invalid arguments for `desc ls`")
	}

	list, err := c.uc.RecordList(c.cfg, gokConsts.DESC)
	if err != nil {
		return fmt.Errorf("failed to get records list - %s", err.Error())
	}

	if len(list) == 0 {
		fmt.Println("No record found")
		return
	}

	fmt.Printf("\n")
	for _, r := range list {
		// Decrypt description for output.
		desc, errDec := r.Description.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt description - %s", errDec.Error())
		}

		fmt.Printf("%s\t %s\t %d\t %s\n", r.ID, r.Type, r.Head, *desc)
	}

	return nil
}
