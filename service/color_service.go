package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godemo/model"
	"github.com/godemo/repository"
)

type ColorService struct {}

func (colorService *ColorService) GetAllColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.GetAll()

	c.JSON(myOutput.StatusCode, myOutput)
}

func (colorService *ColorService) GetSingleColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var output model.Color
	c.ShouldBind(&output)
	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.GetByID(output)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (colorService *ColorService) InsertColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var output model.Color
	c.ShouldBind(&output)
	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.Insert(output)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (colorService *ColorService) UpdateColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var output model.Color
	c.ShouldBind(&output)
	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.Update(output)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (colorService *ColorService) DeleteColorService(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")
	var output model.Color
	c.ShouldBind(&output)
	colorRepository := new(repository.ColorRepository)
	myOutput := colorRepository.Delete(output)
	c.JSON(myOutput.StatusCode, myOutput)
}

func (colorService *ColorService) AddRouters(router *gin.Engine){
	router.GET("/getallcolor", colorService.GetAllColorService)
	router.POST("/getsinglecolor", colorService.GetSingleColorService)
	router.POST("/insertcolor", colorService.InsertColorService)
	router.PATCH("/updatecolor", colorService.UpdateColorService)
	router.DELETE("/deletecolor", colorService.DeleteColorService)
}