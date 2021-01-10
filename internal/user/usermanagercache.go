/*
 * DHBW Mosbach project of subject "Programmieren 2" from:
 * 6439456
 * 8093702
 * 9752762
 */

package user

var userManager *UserManager

// The UserManagerCache is used to only load the users in go objects and store it in the variable userManager
// at runtime.
// This way the usersFile doesn't need to be loaded at every request.
// This results in less IO usage and faster response times.
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
