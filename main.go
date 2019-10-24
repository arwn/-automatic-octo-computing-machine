package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arwn/apiwrap"
)

func makeProjectCounter() func() (int, string) {
	projects := []string{"ft_ls", "fdf", "ft_printf"}
	count := 2 // start at 2, first project is at 3
	return func() (int, string) {
		count++
		return count, projects[count-3] // 3 is the offset from start
	}
}

func getUsers(client apiwrap.WrapperClient) []int64 {
	type User struct {
		ID    int64
		Login string
		URL   string
	}

	var userIDs []int64
	page := 0

	for { // loop over all users in campus 7 (ours ;)
		body, err := client.GetBody(fmt.Sprintf("/v2/campus/7/users?page=%d", page))
		if err != nil {
			log.Println(err)
			log.Println(body)
			break
		}
		if len(body) == 0 {
			break
		}

		var users []User
		err = json.Unmarshal(body, &users)
		if err != nil {
			log.Println(err)
			log.Println(body)
			break
		}

		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}

		page++
		time.Sleep(time.Second) // don't dos the api
	}
	return userIDs
}

func main() {
	uid := os.Getenv("user_id")
	secret := os.Getenv("user_secret")
	client := apiwrap.NewClient(uid, secret)
	fmtstr := "/v2/users/%d/projects_users?filter[project_id]=%d&page[size]=100"

	users := getUsers(client)

	fmt.Println(users)

	for _, user := range users {
		fmt.Printf("checking user %d\n", user)
		counter := makeProjectCounter()
		for { // each project {ls, fdf, printf}
			projectID, projectName := counter()
			content, err := client.GetJSON(fmt.Sprintf(fmtstr, user, projectID))
			if err != nil {
				log.Fatal(err)
			}
			var dat []map[string]interface{}
			err = json.Unmarshal(content, &dat)
			if err != nil {
				log.Fatal(err)
			}
			if len(dat) > 0 && dat[0]["final_mark"] != nil {
				fmt.Println(fmt.Sprintf("%d %s: %1f", user, projectName, dat[0]["final_mark"]))
			}
			if projectID == 5 { // id of ft_printf
				break
			}
			time.Sleep(time.Second) // sleep for api rate limit
		}
	}

}
