package ios_app

const iOSAppPrefix = "https://itunes.apple.com/cn/app/id"

func PageURL(id string) string {
	return iOSAppPrefix + id
}
