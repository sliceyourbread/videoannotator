package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string       `json:"username"`
	Annotation []Annotation `json:"annotation"`
}

type Video struct {
	gorm.Model
	Link       string `json:"link"`
	Provider   string `json:"provider"`
	Length     int64  `json:"length"`
	Quality    string `json:"quality"`
	Annotation []Annotation
}

type Annotation struct {
	gorm.Model
	UserID   uint
	User     User `gorm:"foreignKey:UserID;references:ID"`
	VideoID  uint
	Video    Video  `gorm:"foreignKey:VideoID;references:ID"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Comment  string `json:"comment"`
	Type     string `json:"type"`
	Language string `json:"language"`
}

func (u *User) ListAnnotations(db *gorm.DB, link string) ([]Annotation, error) {

	var err error

	var v Video
	result := db.Where("link=?", link).First(&v)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no video found aborting delete")
	}

	result = db.Debug().Where("username=?", u.Username).First(&u)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no user found")
	}

	a := []Annotation{}
	err = db.Model(&u).Association("Annotation").Find(&a)
	if err != nil {
		return nil, err
	}

	return a, err
}

func (u *User) DeleteAnnotation(db *gorm.DB, link string) (int64, error) {

	var err error
	var v Video
	result := db.Where("link=?", link).First(&v)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no video found aborting delete")
	}

	result = db.Debug().Where("username=?", u.Username).First(&u)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no user found")
	}

	var rows int64
	for i := range u.Annotation {
		u.Annotation[i].VideoID = v.ID
		u.Annotation[i].UserID = u.ID
		result = db.Delete(&u.Annotation[i])
		if result.Error != nil {
			err = result.Error
		}

		rows += result.RowsAffected
	}

	return rows, err

}

func (u *User) AddAnnotation(db *gorm.DB, link string) (int64, error) {

	var v Video
	result := db.Where("link=?", link).First(&v)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no video found aborting delete")
	}

	for i := range u.Annotation {
		u.Annotation[i].VideoID = v.ID
	}
	// result := db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "link"}},
	// 	UpdateAll: true,
	// }).Create(&u)

	result = db.Debug().Create(&u)
	return result.RowsAffected, result.Error

}

func (u *User) UpdateAnnotation(db *gorm.DB, link string) (int64, error) {
	var v Video
	result := db.Debug().Where("link=?", link).First(&v)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no video found aborting delete")
	}

	result = db.Debug().Where("username=?", u.Username).First(&u)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no user found")
	}

	var err error
	var rows int64
	for i, a := range u.Annotation {
		a.VideoID = v.ID
		// result = db.Debug().Where("video_id", v.ID).Where("user_id", u.ID).Where("type", a.Type).First(&tmp)
		// if result.Error != nil {
		// 	err = result.Error
		// }

		result = db.Debug().Model(&u.Annotation[i]).Where("video_id", v.ID).Where("user_id", u.ID).Where("type", a.Type).Updates(Annotation{Start: u.Annotation[i].Start, End: u.Annotation[i].End, Comment: u.Annotation[i].Comment, Language: u.Annotation[i].Language})
		rows += result.RowsAffected
	}

	return rows, err
}

func (v *Video) AddVideo(db *gorm.DB) (int64, error) {

	result := db.Where("link=?", v.Link).First(&v)
	if result.RowsAffected != 0 {
		return 0, fmt.Errorf("video already present")
	}

	result = db.Create(&v)
	return result.RowsAffected, result.Error
}

func (v *Video) DeleteVideo(db *gorm.DB) (int64, error) {

	result := db.Where("link=?", v.Link).First(&v)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("no video found aborting delete")
	}

	result = db.Select("Annotation").Delete(&v)
	return result.RowsAffected, result.Error
}
