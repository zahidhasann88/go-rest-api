package repository

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/godemo/model"
	"github.com/godemo/model/custommodel"
	"github.com/godemo/util"
)

type ColorRepository struct{}

// GetALL
func (colorRepository *ColorRepository) GetAll() custommodel.ResponseDto {
	var output custommodel.ResponseDto

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var color []model.Color
	result := db.Order("color_id desc").Find(&color)

	if result.RowsAffected == 0 {
		output.Message = "Srver Error So no data found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}

	type tempOutput struct {
		Output      []model.Color `json:"output"`
		OutputCount int           `json:"outputCount"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	tOutput.OutputCount = len(color)
	output.Message = "Successfully Get All Colors"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

// GetByID
func (colorRepository *ColorRepository) GetByID(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_id == 0 {
		output.Message = "Color ID is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	result := db.Where(&model.Color{Color_id: color.Color_id}).First(&color)
	if result.RowsAffected == 0 {
		output.Message = "No data found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}
	type tempOutput struct {
		Output model.Color `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	output.Message = "Color info details found for given criteria"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

// Insert Color
func (colorRepository *ColorRepository) Insert(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_name == "" {
		output.Message = "Color name is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var output1 model.Color
	result1 := db.Where("lower(color_name) = ?", strings.ToLower(color.Color_name)).First(&output1)
	if result1.RowsAffected > 0 {
		output.Message = color.Color_name + " Color already exists"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusConflict
		return output
	}
	// ID Autoincrement
	_ = db.Raw("select coalesce ((max(color_id) + 1), 1) from public.color").First(&color.Color_id)

	result := db.Create(&color)
	if result.RowsAffected == 0 {
		output.Message = "Color not inserted for Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}
	type tempOutput struct {
		Output model.Color `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	output.Message = "Color inserted successfully"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

// update color
func (colorRepository *ColorRepository) Update(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_id == 0 {
		output.Message = "Color ID is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}
	if color.Color_name == "" {
		output.Message = "Color name is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	var output1 model.Color
	result0 := tx.Where("color_id = ?", color.Color_id).First(&output1)
	id := strconv.Itoa(color.Color_id)
	if result0.RowsAffected == 0 {
		output.Message = "this id " + id + " Not Found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}

	result1 := db.Where("lower(color_name) = ?", strings.ToLower(color.Color_name)).First(&output1)
	if result1.RowsAffected > 0 {
		output.Message = color.Color_name + " Color already exists"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusConflict
		return output
	}

	var archColor model.Color_archive
	var output2 model.Color_archive
	archColor.Color_id = color.Color_id
	archColor.Color_name = color.Color_name
	dt := time.Now()
	archColor.Changedate = dt.Format("2006-01-02 15:04:05")
	archColor.Changeflag = "Update"

	_ = db.Raw("select coalesce ((max(trackid) + 1), 1) from public.color_archive").First(&output2.Trackid)

	archColor.Trackid = output2.Trackid
	archColor.Changeuser = "Admin"
	result2 := tx.Create(&archColor)
	if result2.RowsAffected == 0 {
		output.Message = "Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return output
	}

	result3 := tx.Where("lower(color_name) = ?", strings.ToLower(color.Color_name)).First(&output1)
	if result3.RowsAffected > 0 {
		output.Message = color.Color_name + " Color already exists"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusConflict
		tx.RollbackTo("savepoint")
		return output
	}

	result := db.Model(&color).Where(&model.Color{Color_id: color.Color_id}).Updates(&color)
	if result.RowsAffected == 0 {
		output.Message = "Color is not updated for Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}

	tx.Commit()

	type tempOutput struct {
		Output model.Color `json:"output"`
	}
	var tOutput tempOutput
	tOutput.Output = color
	output.Message = "Color updated successfully"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}

// delete color
func (colorRepository *ColorRepository) Delete(color model.Color) custommodel.ResponseDto {
	var output custommodel.ResponseDto
	if color.Color_id == 0 {
		output.Message = "Color ID is required"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusBadRequest
		return output
	}

	db := util.CreateConnectionUsingGormToCommonSchema()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	tx := db.Begin()
	tx.SavePoint("savepoint")

	result0 := db.Where("color_id = ?", color.Color_id).First(&color)
	id := strconv.Itoa(color.Color_id)
	if result0.RowsAffected == 0 {
		output.Message = "this id " + id + " Not Found"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusNotFound
		return output
	}

	var archColor model.Color_archive
	var output2 model.Color_archive
	archColor.Color_id = color.Color_id
	archColor.Color_name = color.Color_name
	dt := time.Now()
	archColor.Changedate = dt.Format("2006-01-02 15:04:05")
	archColor.Changeflag = "Delete"

	_ = db.Raw("select coalesce ((max(trackid) + 1), 1) from public.color_archive").First(&output2.Trackid)

	archColor.Trackid = output2.Trackid
	archColor.Changeuser = "Admin"
	result2 := tx.Create(&archColor)
	if result2.RowsAffected == 0 {
		output.Message = "Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		tx.RollbackTo("savepoint")
		return output
	}

	result := db.Where("color_id = ?", color.Color_id).Delete(&color)
	if result.RowsAffected == 0 {
		output.Message = "Color is not deleted for Internal Server Error"
		output.IsSuccess = false
		output.Payload = nil
		output.StatusCode = http.StatusInternalServerError
		return output
	}

	tx.Commit()

	type tempOutput struct {
		Output model.Color `json:"output"`
	}

	var tOutput tempOutput
	tOutput.Output = color
	output.Message = "Color deleted successfully"
	output.IsSuccess = true
	output.Payload = tOutput
	output.StatusCode = http.StatusOK
	return output
}
