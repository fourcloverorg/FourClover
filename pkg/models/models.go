package models

// Result represents the result model
type Result struct {
	ID        int
	Name      string
	Data      []byte `json:"-"` // don't show data in json response
	CreatedAt string
}

// NewResult creates a new result instance
func NewResult(id int, name string, data []byte, createdAt string) *Result {
	return &Result{
		ID:        id,
		Name:      name,
		Data:      data,
		CreatedAt: createdAt,
	}
}

// GetID returns the ID of the result
func (r *Result) GetID() int {
	return r.ID
}

// GetName returns the name of the result
func (r *Result) GetName() string {
	return r.Name
}

// GetData returns the data of the result
func (r *Result) GetData() string {
	return string(r.Data)
}

// GetCreatedAt returns the created at of the result
func (r *Result) GetCreatedAt() string {
	return r.CreatedAt
}

// SetResult sets the result
func (r *Result) SetResult(result Result) {
	r.ID = result.ID
	r.Name = result.Name
	r.Data = result.Data
	r.CreatedAt = result.CreatedAt
}
