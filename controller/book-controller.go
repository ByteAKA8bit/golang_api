package controller

import (
	"fmt"
	"github.com/ByteAKA8bit/golang_api/dto"
	"github.com/ByteAKA8bit/golang_api/entity"
	"github.com/ByteAKA8bit/golang_api/helper"
	"github.com/ByteAKA8bit/golang_api/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BookController interface {
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	FindByID(context *gin.Context)
	All(context *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtSevice   service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService,
		jwtService,
	}
}

func (controller *bookController) Insert(context *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := context.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("authorization")
		userID := controller.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			bookCreateDTO.UserID = convertedUserID
		}
		result := controller.bookService.Insert(bookCreateDTO)
		res := helper.BuildResponse(true, "Ok", result)
		context.JSON(http.StatusOK, res)
	}
}

func (controller *bookController) Update(context *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := context.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	authHeader := context.GetHeader("authorization")
	userID := controller.getUserIDByToken(authHeader)
	if controller.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		result := controller.bookService.Update(bookUpdateDTO)
		res := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You don't have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (controller *bookController) Delete(context *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	book.ID = id
	authHeader := context.GetHeader("authorization")
	userID := controller.getUserIDByToken(authHeader)
	if controller.bookService.IsAllowedToEdit(userID, book.ID) {
		controller.bookService.Delete(book)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You don't have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (controller *bookController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), nil)
		context.JSON(http.StatusBadRequest, res)
	}

	var book entity.Book = controller.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "Ok", book)
		context.JSON(http.StatusOK, res)
	}
}

func (controller *bookController) All(context *gin.Context) {
	var books []entity.Book = controller.bookService.All()
	res := helper.BuildResponse(true, "OK", books)
	context.JSON(http.StatusOK, res)
}

func (controller *bookController) getUserIDByToken(token string) string {
	aToken, err := controller.jwtSevice.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
