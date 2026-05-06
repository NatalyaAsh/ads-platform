package model

import "errors"

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrAdNotFound            = errors.New("ad not found")
	ErrCategoryAlreadyExists = errors.New("category alredy exists")
	ErrCategoryHasAds        = errors.New("cannot delete category with existing ads")
	ErrInvalidPrice          = errors.New("price must be greater than zero")
	ErrInvalidTitle          = errors.New("title is required")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrForbidden             = errors.New("you don't have permission to edit this ad")
)
