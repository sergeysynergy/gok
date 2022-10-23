package cli

import (
	"fmt"
	"github.com/google/uuid"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

// Method desc provide local work with DESC record type.
func (c *CLI) text() {
	if len(c.args) == 1 {
		c.descLs()
		return
	}

	switch c.args[1] {
	case "add":
		err := c.textAdd()
		if err != nil {
			fmt.Println("\nText add failed:", err)
		} else {
			fmt.Println("\nSuccessfully added new text record.")
		}
	case "set":
		err := c.textSet()
		if err != nil {
			fmt.Println("\nText set failed:", err)
		} else {
			fmt.Println("\nSuccessfully updated text record.")
		}
	case "ls":
		c.textLs()
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) textAdd() (err error) {
	if len(c.args) < 4 {
		return fmt.Errorf("insufficient arguments for `text add [description] [text]`")
	}

	desc := c.args[2]
	text := c.args[3]
	rec := &entity.Record{
		Description: entity.StringField(desc),
		//Type:        gokConsts.TEXT,
		Extension: &entity.Text{
			Text: entity.StringField(text),
		},
	}
	err = c.uc.RecordAdd(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to add new text record - %s", err.Error())
	}

	return nil
}

func (c *CLI) textSet() (err error) {
	if len(c.args) < 5 {
		return fmt.Errorf("insufficient arguments for `text set [ID] [description] [text]`")
	}

	id := c.args[2]
	_, err = uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid record ID - %s", err.Error())
	}

	rec := &entity.Record{
		ID:          entity.RecordID(id),
		Description: entity.StringField(c.args[3]),
		Extension: &entity.Text{
			Text: entity.StringField(c.args[4]),
		},
	}

	err = c.uc.RecordSet(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to update description record - %s", err.Error())
	}

	return nil
}

func (c *CLI) textLs() (err error) {
	if len(c.args) > 2 {
		return fmt.Errorf("invalid arguments for `desc ls`")
	}

	list, err := c.uc.RecordList(c.cfg, gokConsts.TEXT)
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

		//Decrypt text for output.
		ex := r.Extension.(*entity.Text)
		text, errDec := ex.Text.Decrypt(c.cfg.Key)
		if errDec != nil {
			return fmt.Errorf("failed to decrypt description - %s", errDec.Error())
		}

		fmt.Printf("%s\t %s\t %d\t %s\t %s\n", r.ID, r.Type, r.Head, *desc, *text)
	}

	return nil
}
