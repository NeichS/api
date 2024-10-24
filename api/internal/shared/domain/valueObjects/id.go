package valueobjects

type (
	ID struct {
		Id uint
	}
)

func NewID(id uint) (*ID, error) {
	return &ID{
		Id: id,
	}, nil
}
