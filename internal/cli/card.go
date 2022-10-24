package cli

import (
	"fmt"
	"github.com/google/uuid"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
	"strconv"
)

// Method card provide local work with CARD record type.
func (c *CLI) card() {
	if len(c.args) == 1 {
		err := c.cardLs()
		if err != nil {
			fmt.Println("\nText ls failed:", err)
		}
		return
	}

	switch c.args[1] {
	case "add":
		err := c.cardAdd()
		if err != nil {
			fmt.Println("\nCard add failed:", err)
		} else {
			fmt.Println("\nSuccessfully added new card record.")
		}
	case "set":
		err := c.cardSet()
		if err != nil {
			fmt.Println("\nCard set failed:", err)
		} else {
			fmt.Println("\nSuccessfully updated card record.")
		}
	case "ls":
		err := c.cardLs()
		if err != nil {
			fmt.Println("\nText ls failed:", err)
		}
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) cardAdd() (err error) {
	if len(c.args) < 7 {
		return fmt.Errorf("insufficient arguments for `card add [description] [number] [expired] [code] [owner]`")
	}

	desc := c.args[2]
	number, err := strconv.Atoi(c.args[3])
	if err != nil {
		return err
	}
	expired := c.args[4]
	code, err := strconv.Atoi(c.args[5])
	if err != nil {
		return err
	}
	owner := c.args[6]

	rec := &entity.Record{
		Description: entity.StringField(desc),
		Extension: &entity.Card{
			Number:  entity.NumberField(number),
			Code:    entity.NumberField(code),
			Expired: entity.StringField(expired),
			Owner:   entity.StringField(owner),
		},
	}
	err = c.uc.RecordAdd(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to add new card record - %s", err.Error())
	}

	return nil
}

func (c *CLI) cardSet() (err error) {
	if len(c.args) < 7 {
		return fmt.Errorf("insufficient arguments for `card set [ID] [description] [number] [expired] [code] [owner]`")
	}

	id := c.args[2]
	_, err = uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid record ID - %s", err.Error())
	}

	desc := c.args[3]
	number, err := strconv.Atoi(c.args[4])
	if err != nil {
		return err
	}
	expired := c.args[5]
	code, err := strconv.Atoi(c.args[6])
	if err != nil {
		return err
	}
	owner := c.args[7]

	rec := &entity.Record{
		ID:          entity.RecordID(id),
		Description: entity.StringField(desc),
		Extension: &entity.Card{
			Number:  entity.NumberField(number),
			Code:    entity.NumberField(code),
			Expired: entity.StringField(expired),
			Owner:   entity.StringField(owner),
		},
	}

	err = c.uc.RecordSet(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to update card record - %s", err.Error())
	}

	return nil
}

func (c *CLI) cardLs() (err error) {
	if len(c.args) > 2 {
		return fmt.Errorf("invalid arguments for `card ls`")
	}

	list, err := c.uc.RecordList(c.cfg, gokConsts.CARD)
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

		ex := r.Extension.(*entity.Card)

		number, errDec := ex.Number.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt number - %s", errDec.Error())
		}
		expired, errDec := ex.Expired.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt expired - %s", errDec.Error())
		}
		code, errDec := ex.Code.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt code - %s", errDec.Error())
		}
		owner, errDec := ex.Owner.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt owner - %s", errDec.Error())
		}

		fmt.Printf("%s\t %s\t %d\t %s\t %d\t %s\t %d\t %s\n", r.ID, r.Type, r.Head, *desc, *number, *expired, *code, *owner)
	}

	return nil
}
