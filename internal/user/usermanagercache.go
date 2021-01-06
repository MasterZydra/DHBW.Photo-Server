package user

var userManager *UserManager

func ImageManagerCache() *UserManager {
	if userManager == nil {
		userManager = NewUserManager()
		_ = userManager.LoadUsers()
	}
	return userManager
}

func ResetImageManagerCache() {
	userManager = nil
}
