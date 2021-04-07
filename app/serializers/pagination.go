package serializers

// Pagination is the serializer for binding limit and offset
type Pagination struct {
	Limit  uint64 `form:"limit,default=20"`
	Offset uint64 `form:"offset,default=0"`
}
