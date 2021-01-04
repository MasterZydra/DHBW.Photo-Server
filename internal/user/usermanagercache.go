package user

var userManager *UserManager

func GetImageManager() *UserManager {
	if userManager == nil {
		userManager = NewUserManager()
		_ = userManager.LoadUsers()
	}
	return userManager
}
