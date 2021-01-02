package user

var userManager *UserManager

// TODO: jones Test schreiben
func GetImageManager() *UserManager {
	if userManager == nil {
		userManager = NewUserManager()
		_ = userManager.LoadUsers()
	}
	return userManager
}
