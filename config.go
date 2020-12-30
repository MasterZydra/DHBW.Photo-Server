package DHBW_Photo_Server

const (
	// userfiles
	TestUserFile = "usersFile_test.csv"
	ProdUserFile = "usersFile.csv"

	// user1
	User1Name = "Max"
	Pw1Clear  = "test123"
	//Pw1Hash = "6dfbf8730f569dba965ead781f536f7b5ccc2f6b9824f0e49e6878b349a94bc9186c7d7145df80e841def14f3dd70791"
	Pw1Hash = "73a64b63aeb9e71d4e10df824ab4a9d32ce1911bf343d085fb67dec7aba0fb711bda08780efc5d9291df3e8e1c7a66b2"

	// user2
	User2Name = "Ana"
	Pw2Clear  = "123test"
	Pw2Hash   = "e0f3ae3d616e121df29464a191b1e5cec18c84190490550230b3a8f93930b71e46a90e876b6896839996675259096fd4"
	//Pw2Hash = "e9fa8567977ba0db64bc5d5f18118d377032a4820c38ed404400b52bdb6751b9e27c0beb37e35f2bf75608c634a28990"

	UsernameRegexBlacklist = `[^a-z^A-Z^0-9\^_\^.\^-]`
	UsernameRegexWhitelist = `[a-zA-Z0-9\_\.\-]+`
)
