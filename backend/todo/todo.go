package todo

import (
	"context"
	"time"
)

// Item is an item contained within a List.
type Item struct {
	// ID represents a globally unique identifier and is set by the TodoService.
	ID int `json:"id"`
	// UserID represents the user who created the Item
	UserID int `json:"userId"`
	// ListID represents the list to which this Item belongs
	ListID int `json:"listId"`
	// Name is a used defined identifier for the Item.
	Name string `json:"name"`
	// Completed indicates whether this Item is completed or not.
	Completed bool `json:"completed"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (i *Item) Validate() error {
	if i.Name == "" {
		return Err(EINVALID, "name required")
	}

	if i.UserID <= 0 {
		return Err(EINVALID, "user id required")
	}

	if i.ListID <= 0 {
		return Err(EINVALID, "list id required")
	}

	return nil
}

// List represents a collection of Items
type List struct {
	// ID represents the unique identifier for this List.
	ID int `json:"id"`
	// UserID represents the ID of the user who created this List.
	UserID int `json:"userId"`
	// Name represents a user assigned identifier for the List.
	Name string `json:"name"`
	// Completed indicates if all items are completed.
	Completed bool `json:"completed"`
	// Items represents all the items contained within this List.
	Items []*Item `json:"items"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (l *List) Validate() error {
	if l.UserID <= 0 {
		return Err(EINVALID, "user id required")
	}

	if l.Name == "" {
		return Err(EINVALID, "name required")
	}

	return nil
}

type ListFilter struct {
	// Filter fields
	ID        *int
	UserID    *int
	Name      *string
	Completed *bool

	// Range restrictions
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ListUpdate struct {
	Name      *string `json:"name,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

type ItemFilter struct {
	// Filter fields
	ID        *int
	UserID    *int
	ListID    *int
	Name      *string
	Completed *bool

	// Range restrictions
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ItemUpdate struct {
	Name      *string `json:"name,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

// ItemListService provides functionality for manipulating Lists and Items.
type ItemListService interface {
	// CreateItem creates a Todo. The ID property of a Todo is ignored if specified.
	// Errors returned:
	//	invalid: the todo specified failed to validate.
	CreateItem(ctx context.Context, i *Item) error
	// CreateList creates a List. The ID property of a List is ignored if specified.
	// Errors returned:
	//	invalid: the todo specified failed to validate.
	CreateList(ctx context.Context, l *List) error
	// DeleteItem deletes a Item by ID.
	// Errors returned:
	//	invalid: an invalid ID was specified
	//	not_found: no matching Todo was found
	DeleteItem(ctx context.Context, id int) error
	// DeleteList deletes a List by ID.
	DeleteList(ctx context.Context, id int) error
	// FindItemByID returns a Todo with the matching ID.
	// Errors returned:
	//	invalid: an invalid ID was specified
	//	not_found: no matching Todo could be found
	FindItemByID(ctx context.Context, id int) (*Item, error)
	// FindItems finds the Items with the matching filters applied as a logical AND.
	// Errors returned:
	//	invalid: an invalid filter was specified
	//	not_found: no matching Items could be found.
	FindItems(ctx context.Context, f ItemFilter) ([]*Item, error)
	// FindListByID returns a List with the matching ID.
	// Errors returned:
	//	invalid: an invalid ID was specified
	//	not_found: no matching Todo could be found
	FindListByID(ctx context.Context, id int) (*List, error)
	// FindLists finds the Items with the matching filters applied as a logical AND.
	// Errors returned:
	//	invalid: an invalid filter was specified
	//	not_found: no matching Items could be found.
	FindLists(ctx context.Context, f ListFilter) ([]*List, error)
	// UpdateItem updates the Name and/or Completed state of a Todo.
	// Errors returned:
	//	invalid: an invalid if no updates were specified.
	//	not_found: no matching Item was found
	UpdateItem(ctx context.Context, id int, upd ItemUpdate) (*Item, error)
	// UpdateList updates the Title and/or Completed state of a Todo.
	// Errors returned:
	//	invalid: an invalid if no updates were specified.
	//	not_found: no matching Todo was found
	UpdateList(ctx context.Context, id int, upd ListUpdate) (*List, error)
}
