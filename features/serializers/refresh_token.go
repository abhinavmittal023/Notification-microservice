package serializers

//RefreshToken struct is a serializer for getting in the refresh token
type RefreshToken struct {
	RefreshToken string `json:"refresh_token,omitempty" binding:"required"`
	AccessToken  string `json:"access_token"`
}
