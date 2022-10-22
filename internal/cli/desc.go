package cli

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"time"

	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
)

// Method desc provide local work with DESC record type.
func (c *CLI) desc() {
	if len(c.args) == 1 {
		c.descLs()
		return
	}

	switch c.args[1] {
	case "add":
		c.descAdd()
	case "set":
		c.descSet()
	case "get":
	case "ls":
		c.descLs()
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) descAdd() {
	if len(c.args) < 3 {
		fmt.Println("Insufficient arguments for `desc add [description]` command: need description field.")
		return
	}

	// concatenate all further arguments in one description
	var description string
	for k, v := range c.args[2:] {
		description += v
		if len(c.args[2:]) != k+1 {
			description += " "
		}
	}

	dsr := entity.NewRecord(
		c.cfg.Key,
		c.cfg.LocalHead+1, // increase head counter for new records
		c.cfg.Branch,
		description,
		time.Now(),
		nil,
	)
	err := c.uc.DescAdd(dsr)
	if err != nil {
		c.lg.Error(err.Error())
		fmt.Println("Failed to add new description record -", err.Error())
		return
	}

	fmt.Println("Successfully added new description.")
}

func (c *CLI) descSet() {
	if len(c.args) < 4 {
		fmt.Println("Insufficient arguments for `desc set [record_id] [description]` command.")
		return
	}

	id := c.args[2]
	_, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("Invalid record ID - ", err.Error())
		return
	}

	// concatenate all further arguments in one description
	var description string
	for k, v := range c.args[3:] {
		description += v
		if len(c.args[2:]) != k+1 {
			description += " "
		}
	}

	dsr := &entity.Record{
		ID:          entity.RecordID(id),
		Head:        c.cfg.LocalHead + 1, // increase head counter
		Branch:      c.cfg.Branch,
		Description: entity.Description(description),
		Type:        gokConsts.DESC,
		UpdatedAt:   time.Now(),
	}

	err = c.uc.DescSet(dsr)
	if err != nil {
		c.lg.Error(err.Error())
		if errors.Is(err, gokErrors.ErrRecordNotFound) {
			fmt.Println("Failed to update description record: record not found -", err.Error())
		} else {
			fmt.Println("Failed to update description record -", err.Error())
		}
		return
	}

	fmt.Println("Successfully updated description.")
}

func (c *CLI) descLs() {
	if len(c.args) > 2 {
		fmt.Println("Invalid arguments for `desc ls` command.")
		return
	}

	list, err := c.uc.DescList()
	if err != nil {
		fmt.Println("Failed to get records list -", err.Error(), list)
		return
	}

	if len(list) == 0 {
		fmt.Println("No record found")
		return
	}

	for _, r := range list {
		desc, errDec := r.Description.Decrypt(c.cfg.Key)
		if errDec != nil {
			fmt.Println("Failed to decrypt description -", errDec.Error())
			return
		}
		fmt.Printf("%s\t %s\t %d\t %s\t %s\n", r.ID, r.Type, r.Head, r.UpdatedAt, *desc)
	}
}
