package model

const (
	TableNameCategoryMember = "category_members"
)

type CategoryMember struct {
	CategoryID int64
	MemberID   int64
}

func NewCategoryMember(categoryID, memberID int64) *CategoryMember {
	return &CategoryMember{
		CategoryID: categoryID,
		MemberID:   memberID,
	}
}
