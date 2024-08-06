package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/berkanturkmen/VatanSoft/database"
	"github.com/berkanturkmen/VatanSoft/models"
	"gorm.io/gorm"
	"time"
	"math/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"github.com/berkanturkmen/VatanSoft/cache"
	"encoding/json"
	"errors"
)

func GetCities(c * fiber.Ctx) error {
	val, err := cache.RDB.Get(cache.Ctx, "cities").Result()

	if err != nil {
		var cities[] struct {
			ID uint `json:"ID"`
			Name string `json:"Name"`
		}

		if err := database.DB.Model( & models.City {}).Select("id", "name").Find( & cities).Error;
		err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		data, err := json.Marshal(cities)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		cache.RDB.Set(cache.Ctx, "cities", data, 24 * time.Hour)

		return c.JSON(cities)
	}

	return c.SendString(val)
}

func GetRegionsByCity(c * fiber.Ctx) error {
	cityID, err := strconv.ParseUint(c.Params("cityID"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	key := "regions:" + strconv.FormatUint(cityID, 10)

	val, err := cache.RDB.Get(cache.Ctx, key).Result()

	if err != nil {
		var regions[] struct {
			ID uint `json:"ID"`
			Name string `json:"Name"`
			CityID uint `json:"CityID"`
		}

		if err := database.DB.Model( & models.Region {}).Select("id", "name", "city_id").Where("city_id = ?", cityID).Find( & regions).Error;
		err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		data, err := json.Marshal(regions)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		cache.RDB.Set(cache.Ctx, key, data, 24 * time.Hour)

		return c.JSON(regions)
	}

	return c.SendString(val)
}

func GetJobs(c * fiber.Ctx) error {
	val, err := cache.RDB.Get(cache.Ctx, "jobs").Result()

	if err != nil {
		var jobs[] struct {
			ID uint `json:"ID"`
			Name string `json:"Name"`
		}

		if err := database.DB.Model( & models.Job {}).Select("id", "name").Find( & jobs).Error;
		err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		data, err := json.Marshal(jobs)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		cache.RDB.Set(cache.Ctx, "jobs", data, 24 * time.Hour)

		return c.JSON(jobs)
	}

	return c.SendString(val)
}

func GetTitlesByJob(c * fiber.Ctx) error {
	jobID, err := strconv.ParseUint(c.Params("jobID"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	key := "titles:" + strconv.FormatUint(jobID, 10)

	val, err := cache.RDB.Get(cache.Ctx, key).Result()

	if err != nil {
		var titles[] struct {
			ID uint `json:"ID"`
			Name string `json:"Name"`
			JobID uint `json:"JobID"`
		}

		if err := database.DB.Model( & models.Title {}).Select("id", "name", "job_id").Where("job_id = ?", jobID).Find( & titles).Error;
		err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		data, err := json.Marshal(titles)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		cache.RDB.Set(cache.Ctx, key, data, 24 * time.Hour)

		return c.JSON(titles)
	}

	return c.SendString(val)
}

func Register(c * fiber.Ctx) error {
	type request struct {
		Hospital models.Hospital `json:"Hospital"`
		Employee models.Employee `json:"Employee"`
	}

	var req request

	if err := c.BodyParser( & req);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	err := database.DB.Transaction(func(tx * gorm.DB) error {
		if err := tx.Create( & req.Hospital).Error;
		err != nil {
			return err
		}

		var employeeCount int64

		tx.Model( & models.Employee {}).Where("hospital_id = ?", req.Hospital.ID).Count( & employeeCount)

		if employeeCount == 0 {
			req.Employee.Permission = "Owner"
		}

		req.Employee.HospitalID = req.Hospital.ID

		if err := tx.Create( & req.Employee).Error;
		err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map {
		"Message": "",
		"Success": true,
	})
}

func GetResetCode(c * fiber.Ctx) error {
	type request struct {
		Telephone string `json:"Telephone"`
	}

	var req request

	if err := c.BodyParser( & req);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	result := database.DB.Where("telephone = ?", req.Telephone).First( & employee)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if time.Now().Before(employee.ResetInformation.NextCodeTime) {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	rand.Seed(time.Now().UnixNano())

	resetCode := fmt.Sprintf("%06d", rand.Intn(900000) + 100000)

	employee.ResetInformation = models.ResetInformation {
		Code: resetCode,
		ExpirationTime: time.Now().Add(10 * time.Minute),
		NextCodeTime: time.Now().Add(1 * time.Minute),
	}

	database.DB.Save( & employee)

	return c.JSON(fiber.Map {
		"Message": resetCode,
		"Success": true,
	})
}

func ResetPassword(c * fiber.Ctx) error {
	type request struct {
		Telephone string `json:"Telephone"`
		Code string `json:"Code"`
		NewPassword string `json:"NewPassword"`
		ConfirmPassword string `json:"ConfirmPassword"`
	}

	var req request

	if err := c.BodyParser( & req);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	result := database.DB.Where("telephone = ?", req.Telephone).First( & employee)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if employee.ResetInformation.Code != req.Code || employee.ResetInformation.ExpirationTime.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employee.Password = req.NewPassword

	database.DB.Save( & employee)

	return c.JSON(fiber.Map {
		"Message": "",
		"Success": true,
	})
}

func Login(c * fiber.Ctx) error {
	type request struct {
		EMailAddress string `json:"EMailAddress,emitempty"`
		Telephone string `json:"Telephone,emitempty"`
		Password string `json:"Password,emitempty"`
	}

	var req request

	if err := c.BodyParser( & req);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	result := database.DB.Where("e_mail_address = ? OR telephone = ?", req.EMailAddress, req.Telephone).First( & employee)

	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if err := bcrypt.CompareHashAndPassword([] byte(employee.Password), [] byte(req.Password));
	err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var sessionCount int64

	database.DB.Model( & models.Session {}).Where("employee_id = ?", employee.ID).Count( & sessionCount)

	if sessionCount >= 2 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	sessionToken, err := generateRandomString(45)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	session := models.Session {
		EmployeeID: employee.ID,
		ExpirationTime: time.Now().Add(10 * time.Minute),
		Token: sessionToken,
	}

	if err := database.DB.Create( & session).Error;
	err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.JSON(fiber.Map {
		"Message": sessionToken,
		"Success": true,
	})
}

func GetCurrentEmployee(c * fiber.Ctx) error {
	employeeID := c.Locals("EmployeeID")

	if employeeID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	result := database.DB.Preload("Hospital").First( & employee, employeeID)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employeeDTO := EmployeeDTO {
		FirstName: employee.FirstName,
		LastName: employee.LastName,
		EMailAddress: employee.EMailAddress,
		IdentityNumber: employee.IdentityNumber,
		Permission: string(employee.Permission),
		Telephone: employee.Telephone,
	}

		return c.JSON(fiber.Map {
		"Message": employeeDTO,
		"Success": true,
	})
}

func Logout(c * fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if len(tokenString) > 7 && tokenString[: 7] == "Bearer " {
		tokenString = tokenString[7: ]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var session models.Session

	if err := database.DB.Where("token = ?", tokenString).First( & session).Error;err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if err := database.DB.Unscoped().Delete( & session).Error;err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.JSON(fiber.Map {
		"Message": "",
		"Success": true,
	})
}

func generateRandomString(n int)(string, error) {
	b := make([] byte, n)
	_,
	err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b),
	nil
}

type EmployeeDTO struct {
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	EMailAddress string `json:"EMailAddress"`
	IdentityNumber string `json:"IdentityNumber"`
	Permission string `json:"Permission"`
	Telephone string `json:"Telephone"`
}

func ListEmployees(c * fiber.Ctx) error {
	var employees[] models.Employee

	if result := database.DB.Find( & employees);
	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employeeDTOs := make([] EmployeeDTO, len(employees))
	for i, emp := range employees {
		employeeDTOs[i] = EmployeeDTO {
			FirstName: emp.FirstName,
			LastName: emp.LastName,
			EMailAddress: emp.EMailAddress,
			IdentityNumber: emp.IdentityNumber,
			Permission: string(emp.Permission),
			Telephone: emp.Telephone,
		}
	}

	return c.JSON(fiber.Map {
		"Message": employeeDTOs,
		"Success": true,
	})
}

func CreateEmployee(c * fiber.Ctx) error {
	employee := new(models.Employee)

	if err := c.BodyParser(employee);err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if employee.EMailAddress == "" || employee.IdentityNumber == "" || employee.Password == "" || employee.Telephone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalID,
	ok := c.Locals("HospitalID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employee.HospitalID = hospitalID

	if result := database.DB.Create(employee);result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employeeDTO := EmployeeDTO {
		FirstName: employee.FirstName,
		LastName: employee.LastName,
		EMailAddress: employee.EMailAddress,
		IdentityNumber: employee.IdentityNumber,
		Permission: string(employee.Permission),
		Telephone: employee.Telephone,
	}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map {
		"Message": employeeDTO,
		"Success": true,
	})
}

func GetEmployee(c * fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	if result := database.DB.First( & employee, id);
	result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employeeDTO := EmployeeDTO {
		FirstName: employee.FirstName,
		LastName: employee.LastName,
		EMailAddress: employee.EMailAddress,
		IdentityNumber: employee.IdentityNumber,
		Permission: string(employee.Permission),
		Telephone: employee.Telephone,
	}

	return c.JSON(fiber.Map {
		"Message": employeeDTO,
		"Success": true,
	})
}

func UpdateEmployee(c * fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	if result := database.DB.First( & employee, id);
	result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalID, ok := c.Locals("HospitalID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if employee.HospitalID != hospitalID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var input models.Employee

	if err := c.BodyParser( & input);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if input.EMailAddress == "" || input.IdentityNumber == "" || input.Telephone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if input.Password != "" {
		employee.Password = input.Password
	}

	employee.FirstName = input.FirstName
	employee.LastName = input.LastName
	employee.EMailAddress = input.EMailAddress
	employee.IdentityNumber = input.IdentityNumber
	employee.Permission = input.Permission
	employee.Telephone = input.Telephone

	if result := database.DB.Save( & employee);
	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	employeeDTO := EmployeeDTO {
		FirstName: employee.FirstName,
		LastName: employee.LastName,
		EMailAddress: employee.EMailAddress,
		IdentityNumber: employee.IdentityNumber,
		Permission: string(employee.Permission),
		Telephone: employee.Telephone,
	}

	return c.JSON(fiber.Map {
		"Message": employeeDTO,
		"Success": true,
	})
}

func DeleteEmployee(c * fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var employee models.Employee

	if result := database.DB.Delete( & employee, id);
	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.JSON(fiber.Map {
		"Message": "",
		"Success": true,
	})
}

func AttachPolyclinic(c * fiber.Ctx) error {
	polyclinicID, err := strconv.ParseUint(c.Params("polyclinicID"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalID, ok := c.Locals("HospitalID").(uint)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var existing models.HospitalPolyclinic

	if result := database.DB.Where("polyclinic_id = ? AND hospital_id = ?", polyclinicID, hospitalID).First( & existing);

	result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalPolyclinic := models.HospitalPolyclinic {
		PolyclinicID: uint(polyclinicID),
		HospitalID: hospitalID,
	}

	if result := database.DB.Create( & hospitalPolyclinic);

	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map {
		"message": "",
		"success": true,
	})
}

func DetachPolyclinic(c * fiber.Ctx) error {
	polyclinicID, err := strconv.ParseUint(c.Params("polyclinicID"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalID, ok := c.Locals("HospitalID").(uint)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if result := database.DB.Where("polyclinic_id = ? AND hospital_id = ?", polyclinicID, hospitalID).Delete( & models.HospitalPolyclinic {});

	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.JSON(fiber.Map {
		"Message": "",
		"Success": true,
	})
}

func GetAllPolyclinics(c * fiber.Ctx) error {
	val, err := cache.RDB.Get(cache.Ctx, "polyclinics").Result()

	if err != nil {
		var polyclinics[] struct {
			ID uint `json:"ID"`
			Name string `json:"Name"`
		}

		if err := database.DB.Model( & models.Polyclinic {}).Select("id", "name").Find( & polyclinics).Error;
		err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		data, err := json.Marshal(polyclinics)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}

		cache.RDB.Set(cache.Ctx, "polyclinics", data, 24 * time.Hour)

		return c.JSON(polyclinics)
	}

	return c.SendString(val)
}

type PolyclinicPersonnelDTO struct {
	PolyclinicID uint `json:"PolyclinicID"`
	Name string `json:"Name"`
	Count int `json:"Count"`
}

func GetPolyclinicsByHospital(c * fiber.Ctx) error {
	hospitalID, ok := c.Locals("HospitalID").(uint)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var hospitalPolyclinics[] models.HospitalPolyclinic

	result := database.DB.Preload("Polyclinic").Preload("Personnels").Where("hospital_id = ?", hospitalID).Find( & hospitalPolyclinics)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	polyclinicsCounts := [] PolyclinicPersonnelDTO {}
	polyclinicMap := make(map[uint] PolyclinicPersonnelDTO)
	totalPersonnelCount := 0

	for _, hp := range hospitalPolyclinics {
		polyclinicMap[hp.Polyclinic.ID] = PolyclinicPersonnelDTO {
			PolyclinicID: hp.Polyclinic.ID,
			Name: hp.Polyclinic.Name,
			Count: 0,
		}
	}

	for _, hp := range hospitalPolyclinics {
		entry := polyclinicMap[hp.Polyclinic.ID]
		personnelCount := len(hp.Personnels)
		entry.Count += personnelCount
		polyclinicMap[hp.Polyclinic.ID] = entry
		totalPersonnelCount += personnelCount
	}

	for _, count := range polyclinicMap {
		polyclinicsCounts = append(polyclinicsCounts, count)
	}

	return c.JSON(fiber.Map {
		"Message": polyclinicsCounts,
		"Success": true,
		"TotalPersonnel": totalPersonnelCount,
	})
}

type PersonnelDTO struct {
	IdentityNumber string `json:"IdentityNumber"`
	Name string `json:"Name"`
	Surname string `json:"Surname"`
	Telephone string `json:"Telephone"`
	WorkingDays string `json:"WorkingDays"`
	JobID uint `json:"JobID"`
	JobName string `json:"JobName"`
	HospitalPolyclinicID * uint `json:"HospitalPolyclinicID"`
	TitleID uint `json:"TitleID"`
	TitleName string `json:"TitleName"`
}

func ensureUniqueTitle(titleID uint) error {
	var title models.Title

	if err := database.DB.First( & title, titleID).Error;
	err != nil {
		return err
	}

	if title.IsUnique {
		var existingCount int64

		database.DB.Model( & models.Personnel {}).Where("title_id = ?", titleID).Count( & existingCount)

		if existingCount > 0 {
			return errors.New("Error!")
		}
	}

	return nil
}

func ListPersonnels(c * fiber.Ctx) error {
	var personnels[] models.Personnel

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 10

	identityNumberFilter := c.Query("IdentityNumber")
	nameFilter := c.Query("Name")
	surnameFilter := c.Query("Surname")
	jobNameFilter := c.Query("Job")
	titleNameFilter := c.Query("Title")

	offset := (page - 1) * limit

	var totalRecords int64

	countQuery := database.DB.Table("personnels")

	if identityNumberFilter != "" {
		countQuery = countQuery.Where("identity_number LIKE ?", "%" + identityNumberFilter + "%")
	}

	if nameFilter != "" {
		countQuery = countQuery.Where("name LIKE ?", "%" + nameFilter + "%")
	}

	if surnameFilter != "" {
		countQuery = countQuery.Where("surname LIKE ?", "%" + surnameFilter + "%")
	}

	if jobNameFilter != "" {
		countQuery = countQuery.Joins("JOIN jobs ON jobs.id = personnels.job_id").Where("jobs.name LIKE ?", "%" + jobNameFilter + "%")
	}

	if titleNameFilter != "" {
		countQuery = countQuery.Joins("JOIN titles ON titles.id = personnels.title_id").Where("titles.name LIKE ?", "%" + titleNameFilter + "%")
	}

	countQuery.Count( & totalRecords)

	totalPages := (totalRecords + int64(limit) - 1) / int64(limit)

	if totalPages == 0 {
		totalPages = 1
	}

	mainQuery := database.DB.Preload("Job").Preload("Title").Table("personnels").Offset(offset).Limit(limit)

	if identityNumberFilter != "" {
		mainQuery = mainQuery.Where("identity_number LIKE ?", "%" + identityNumberFilter + "%")
	}

	if nameFilter != "" {
		mainQuery = mainQuery.Where("name LIKE ?", "%" + nameFilter + "%")
	}

	if surnameFilter != "" {
		mainQuery = mainQuery.Where("surname LIKE ?", "%" + surnameFilter + "%")
	}

	if jobNameFilter != "" {
		mainQuery = mainQuery.Joins("JOIN jobs ON jobs.id = personnels.job_id").Where("jobs.name LIKE ?", "%" + jobNameFilter + "%")
	}

	if titleNameFilter != "" {
		mainQuery = mainQuery.Joins("JOIN titles ON titles.id = personnels.title_id").Where("titles.name LIKE ?", "%" + titleNameFilter + "%")
	}

	if result := mainQuery.Find( & personnels);

	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	personnelDTOs := make([] PersonnelDTO, len(personnels))

	for i, per := range personnels {
		personnelDTOs[i] = PersonnelDTO {
			IdentityNumber: per.IdentityNumber,
			Name: per.Name,
			Surname: per.Surname,
			Telephone: per.Telephone,
			WorkingDays: per.WorkingDays,
			JobID: per.JobID,
			JobName: per.Job.Name,
			HospitalPolyclinicID: per.HospitalPolyclinicID,
			TitleID: per.TitleID,
			TitleName: per.Title.Name,
		}
	}

	return c.JSON(fiber.Map {
		"Message": personnelDTOs,
		"Success": true,
		"Pagination": fiber.Map {
			"TotalRecords": totalRecords,
			"TotalPages": totalPages,
			"CurrentPage": page,
			"Limit": limit,
		},
	})
}

func CreatePersonnel(c * fiber.Ctx) error {
	personnel := new(models.Personnel)

	if err := c.BodyParser(personnel);err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if personnel.IdentityNumber == "" || personnel.Telephone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalID,
	ok := c.Locals("HospitalID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	personnel.HospitalID = hospitalID

	if personnel.HospitalPolyclinicID != nil {
		var hp models.HospitalPolyclinic

		if err := database.DB.First( & hp, * personnel.HospitalPolyclinicID).Error;
		err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}
	}

	if err := ensureUniqueTitle(personnel.TitleID);err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if result := database.DB.Create(personnel);result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	personnelDTO := PersonnelDTO {
		IdentityNumber: personnel.IdentityNumber,
		Name: personnel.Name,
		Surname: personnel.Surname,
		Telephone: personnel.Telephone,
		WorkingDays: personnel.WorkingDays,
		JobID: personnel.JobID,
		HospitalPolyclinicID: personnel.HospitalPolyclinicID,
		TitleID: personnel.TitleID,
	}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map {
		"Message": personnelDTO,
		"Success": true,
	})
}

func GetPersonnel(c * fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var personnel models.Personnel

	if result := database.DB.First( & personnel, id);
	result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	personnelDTO := PersonnelDTO {
		IdentityNumber: personnel.IdentityNumber,
		Name: personnel.Name,
		Surname: personnel.Surname,
		Telephone: personnel.Telephone,
		WorkingDays: personnel.WorkingDays,
		JobID: personnel.JobID,
		HospitalPolyclinicID: personnel.HospitalPolyclinicID,
		TitleID: personnel.TitleID,
	}

	return c.JSON(fiber.Map {
		"Message": personnelDTO,
		"Success": true,
	})
}

func UpdatePersonnel(c * fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var personnel models.Personnel
	if result := database.DB.First( & personnel, id);
	result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	hospitalID, ok := c.Locals("HospitalID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if personnel.HospitalID != hospitalID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var input models.Personnel
	if err := c.BodyParser( & input);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if input.IdentityNumber == "" || input.Telephone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	if input.HospitalPolyclinicID != nil {
		var hp models.HospitalPolyclinic
		if err := database.DB.First( & hp, * input.HospitalPolyclinicID).Error;
		err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
				"Message": "",
				"Success": false,
			})
		}
	}

	if err := ensureUniqueTitle(personnel.TitleID);
	err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	personnel.IdentityNumber = input.IdentityNumber
	personnel.Name = input.Name
	personnel.Surname = input.Surname
	personnel.Telephone = input.Telephone
	personnel.WorkingDays = input.WorkingDays
	personnel.JobID = input.JobID
	personnel.HospitalPolyclinicID = input.HospitalPolyclinicID
	personnel.TitleID = input.TitleID

	if result := database.DB.Save( & personnel);
	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	personnelDTO := PersonnelDTO {
		IdentityNumber: personnel.IdentityNumber,
		Name: personnel.Name,
		Surname: personnel.Surname,
		Telephone: personnel.Telephone,
		WorkingDays: personnel.WorkingDays,
		JobID: personnel.JobID,
		HospitalPolyclinicID: personnel.HospitalPolyclinicID,
		TitleID: personnel.TitleID,
	}

	return c.JSON(fiber.Map {
		"Message": personnelDTO,
		"Success": true,
	})
}

func DeletePersonnel(c * fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	var personnel models.Personnel

	if result := database.DB.Delete( & personnel, id);
	result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"Message": "",
			"Success": false,
		})
	}

	return c.JSON(fiber.Map {
		"Message": "",
		"Success": true,
	})
}
