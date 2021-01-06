package user

var userManager *UserManager

// TODO: jones Documentation

func UserManagerCache() *UserManager {
	if userManager == nil {
		userManager = NewUserManager()
		_ = userManager.LoadUsers()
	}
	return userManager
}

func ResetUserManagerCache() {
	userManager = nil
}
