package model

import "time"

type File struct {
	ID               int64
	FileName         string
	FileExt          string
	UserID           int64
	ShortUrlOriginal *string
	ShortUrlUpgrade  *string
	CreatedAt        int64
	StateFile        string
}

func NewFile(fileName string, fileExt string, userID int64) *File {
	return &File{
		FileName:  fileName,
		FileExt:   fileExt,
		UserID:    userID,
		CreatedAt: time.Now().UnixNano(),
		//StateFile: StateFileNotSend,
	}
}

func (f *File) SetShortUrlOriginal(url *string) {
	f.ShortUrlOriginal = url
}
func (f *File) SetShortUrlUpgrade(url *string) {
	f.ShortUrlUpgrade = url
}

const (
	StateFileSend    = "state_file_send"
	StateFileNotSend = "state_file_not_send"
)
