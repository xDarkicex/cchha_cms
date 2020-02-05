package controllers

// // GenerateURL ...
// func generateURL(host, encodedUUID string) string {
// 	return fmt.Sprintf("%s/bd-admin/account/reset?token=%s", host, encodedUUID)
// }

// // NewReset wraps model function create
// func NewReset(r *http.Request, user models.User, token string) models.Reset {
// 	reset := models.Reset{
// 		UserID:  user.ID,
// 		Created: time.Now(),
// 		Expires: time.Now().Add(24 * time.Hour),
// 		Active:  true,
// 		Token:   token,
// 	}
// 	nr := models.CreateReset(reset)
// 	fmt.Println(nr)
// 	return nr
// }
