package model

type File struct {
	BaseModel `xorm:"extends"`
	Key       string `json:"key" xorm:"UNIQUE NOT NULL"`
	UserId    string `json:"user_id" xorm:"varchar(20)"`
	Name      string `json:"name"`
	Hash      string `json:"hash" xorm:"UNIQUE NOT NULL"`
	Size      int64  `json:"size"`
}

func GetFileById(id string) (*File, error) {
	var file *File
	err := GetById(id, file)
	return file, err
}

func DeleteFile(id string) error {
	err := DeleteById(id, &File{})
	return err
}
