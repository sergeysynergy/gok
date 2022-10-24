package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
	"os"
	"path"
)

// Method file provide local work with FILE record type.
func (c *CLI) file() {
	if len(c.args) == 1 {
		err := c.fileLs()
		if err != nil {
			fmt.Println("\nText ls failed:", err)
		}
		return
	}

	switch c.args[1] {
	case "add":
		err := c.fileAdd()
		if err != nil {
			fmt.Println("\nFile add failed:", err)
		} else {
			fmt.Println("\nSuccessfully added new file record.")
		}
	case "set":
		err := c.fileSet()
		if err != nil {
			fmt.Println("\nFile set failed:", err)
		} else {
			fmt.Println("\nSuccessfully updated file record.")
		}
	case "ls":
		err := c.fileLs()
		if err != nil {
			fmt.Println("\nFile ls failed:", err)
		}
	case "get":
		err := c.fileGet()
		if err != nil {
			fmt.Println("\nFile get failed:", err)
		}
	default:
		fmt.Println(c.helpMsg)
	}
}

func (c *CLI) fileAdd() (err error) {
	if len(c.args) < 3 {
		return fmt.Errorf("insufficient arguments for `file add [description] [file_path]`")
	}

	filePath := c.args[2]
	fileName := path.Base(filePath)

	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	rec := &entity.Record{
		Description: entity.StringField(fileName),
		Extension: &entity.File{
			File: dst,
		},
	}
	err = c.uc.RecordAdd(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to add new file record - %s", err.Error())
	}

	return nil
}

func (c *CLI) fileSet() (err error) {
	if len(c.args) < 4 {
		return fmt.Errorf("insufficient arguments for `file set [ID] [file_path]`")
	}

	id := c.args[2]
	_, err = uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid record ID - %s", err.Error())
	}

	filePath := c.args[3]
	fileName := path.Base(filePath)

	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	rec := &entity.Record{
		ID:          entity.RecordID(id),
		Description: entity.StringField(fileName),
		Extension: &entity.File{
			File: dst,
		},
	}

	err = c.uc.RecordSet(c.cfg, rec)
	if err != nil {
		return fmt.Errorf("failed to update file record - %s", err.Error())
	}

	return nil
}

func (c *CLI) fileLs() (err error) {
	if len(c.args) > 2 {
		return fmt.Errorf("invalid arguments for `file ls`")
	}

	list, err := c.uc.RecordList(c.cfg, gokConsts.FILE)
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

func (c *CLI) fileGet() (err error) {
	if len(c.args) < 3 {
		return fmt.Errorf("insufficient arguments for `file set [ID]`")
	}

	id := c.args[2]
	_, err = uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid record ID - %s", err.Error())
	}

	rec, err := c.uc.RecordGet(entity.RecordID(id))
	if err != nil {
		return fmt.Errorf("failed to update file record - %s", err.Error())
	}

	src := rec.Extension.(*entity.File).File
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err = hex.Decode(dst, src)
	if err != nil {
		return err
	}

	err = os.WriteFile(string(rec.Description), dst, 0777)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully safe file " + rec.Description)
	return nil
}
