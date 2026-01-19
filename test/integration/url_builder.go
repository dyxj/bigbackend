package integration

import "fmt"

func buildUserProfileUrl(url string, userId string) string {
	return fmt.Sprintf("%s/api/v1/user/%s/profile", url, userId)
}
