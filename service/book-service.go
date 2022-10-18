package service

import (
	"fmt"
	"github.com/ByteAKA8bit/golang_api/dto"
	"github.com/ByteAKA8bit/golang_api/entity"
	"github.com/ByteAKA8bit/golang_api/repository"
	"github.com/mashingan/smapping"
	"log"
)

type BookService interface {
	Insert(b dto.BookCreateDTO) entity.Book
	Update(b dto.BookUpdateDTO) entity.Book
	Delete(book entity.Book)
	All() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{
		bookRepository,
	}
}

func (service *bookService) Insert(bookDTO dto.BookCreateDTO) entity.Book {
	var book entity.Book
	err := smapping.FillStruct(&book, smapping.MapFields(&bookDTO))
	if err != nil {
		log.Fatalf("Failed to map %v", err)
	}
	res := service.bookRepository.InsertBook(book)
	return res
}

func (service *bookService) Update(bookDTO dto.BookUpdateDTO) entity.Book {
	var book entity.Book
	err := smapping.FillStruct(&book, smapping.MapFields(&bookDTO))
	if err != nil {
		log.Fatalf("Failed to map %v", err)
	}
	res := service.bookRepository.UpdateBook(book)
	return res
}

func (service *bookService) Delete(book entity.Book) {
	service.bookRepository.DeleteBook(book)
}

func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBook()
}

func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
