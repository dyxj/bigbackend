package integration

import "fmt"

func buildUserProfileUrl(url string, userId string) string {
	return fmt.Sprintf("%s/user/%s/profile", url, userId)
}
