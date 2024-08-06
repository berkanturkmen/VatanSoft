package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type City struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique"`
	Regions[] Region `gorm:"foreignKey:CityID"`
}

type Region struct {
	gorm.Model
	Name string `gorm:"type:varchar(255)"`
	City City `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CityID uint `gorm:"index"`;
}

type Hospital struct {
	gorm.Model
	EMailAddress string `gorm:"type:varchar(255);unique"`
	PublicAddress string `gorm:"type:varchar(255)"`
	TaxIdentificationNumber string `gorm:"type:varchar(255);unique"`
	Telephone string `gorm:"type:varchar(255);unique"`

	City City `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CityID uint `gorm:"index"`;
	Region Region `gorm:"foreignKey:RegionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RegionID uint `gorm:"index"`;

	Employees[] Employee `gorm:"foreignKey:HospitalID"`
	Polyclinics[] HospitalPolyclinic `gorm:"foreignKey:HospitalID"`
}

type ResetInformation struct {
	Code string `json:"Code"`
	ExpirationTime time.Time `json:"ExpirationTime"`
	NextCodeTime time.Time `json:"NextCodeTime"`
}

type Permission string

const (
	PermissionOwner Permission = "Owner"
	PermissionEmployee Permission = "Employee"
)

type Employee struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(255)"`
	LastName string `gorm:"type:varchar(255)"`

	EMailAddress string `gorm:"type:varchar(255);unique"`
	IdentityNumber string `gorm:"type:varchar(255);unique"`
	Password string `gorm:"type:varchar(255)"`
	Permission Permission `gorm:"type:varchar(255);check:permission IN ('Owner', 'Employee')"`
	Telephone string `gorm:"type:varchar(255);unique"`

	ResetInformation ResetInformation `gorm:"type:jsonb"`

	Hospital Hospital `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	HospitalID uint `gorm:"index"`;
}

type Session struct {
	gorm.Model
	EmployeeID uint `gorm:"not null"`
	ExpirationTime time.Time `json:"ExpirationTime"`
	Token string `gorm:"not null;unique"`
}

type Polyclinic struct {
	gorm.Model
	Name string `gorm:"type:varchar(255)"`
	HospitalPolyclinics[] HospitalPolyclinic `gorm:"foreignKey:PolyclinicID"`
}

type HospitalPolyclinic struct {
	gorm.Model
	Polyclinic Polyclinic `gorm:"foreignKey:PolyclinicID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PolyclinicID uint `gorm:"index"`;
	Hospital Hospital `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	HospitalID uint `gorm:"index"`;

	Personnels[] Personnel `gorm:"foreignKey:HospitalPolyclinicID"`
}

type Job struct {
	gorm.Model
	Name string `gorm:"type:varchar(255)"`

	Titles[] Title `gorm:"foreignKey:JobID"`
}

type Title struct {
	gorm.Model
	Name string `gorm:"type:varchar(255)"`
	IsUnique bool `gorm:"default:false;not null"`

	Job Job `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	JobID uint `gorm:"index"`;
}

type Personnel struct {
	gorm.Model
	IdentityNumber string `gorm:"type:varchar(255);unique"`
	Name string `gorm:"type:varchar(255)"`
	Surname string `gorm:"type:varchar(255)"`
	Telephone string `gorm:"type:varchar(255);unique"`
	WorkingDays string `gorm:"type:varchar(255)"`

	Hospital Hospital `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	HospitalID uint `gorm:"index"`;

	Job Job `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	JobID uint `gorm:"index"`;

	HospitalPolyclinic HospitalPolyclinic `gorm:"foreignKey:HospitalPolyclinicID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	HospitalPolyclinicID * uint `gorm:"index"`;

	Title Title `gorm:"foreignKey:TitleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	TitleID uint `gorm:"index"`;
}

func(ri * ResetInformation) Scan(value interface {}) error {
	bytes, ok := value.([] byte)

	if !ok {
		return errors.New("Error!")
	}

	return json.Unmarshal(bytes, & ri)
}

func(ri ResetInformation) Value()(driver.Value, error) {
	return json.Marshal(ri)
}

func(e * Employee) BeforeSave(tx * gorm.DB) error {
	if e.Permission != PermissionOwner && e.Permission != PermissionEmployee {
		return errors.New("Error!")
	}

	if len(e.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([] byte(e.Password), bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		e.Password = string(hashedPassword)
	}

	return nil
}
