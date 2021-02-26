package miscquery

// Iterator struct is a serializer for providing iteration support to get specific record
type Iterator struct {
	Next     bool `form:"next"`
	Previous bool `form:"prev"`
}
