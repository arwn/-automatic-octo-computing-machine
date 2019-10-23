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

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: ./automatic-octo-computing-machine user")
		fmt.Println("environment variables: user_id, user_secret")
		os.Exit(1)
	}

	uid := os.Getenv("user_id")
	secret := os.Getenv("user_secret")
	username := os.Args[3]
	client := apiwrap.NewClient(uid, secret)
	fmtstr := "/v2/users/%d/projects_users?filter[project_id]=%d&page[size]=100"

	userID, err := client.GetUserID(username)
	if err != nil {
		fmt.Println(username + ": user not found")
		os.Exit(1)
	}

	counter := makeProjectCounter()
	for {
		projectID, projectName := counter()
		content, err := client.GetJSON(fmt.Sprintf(fmtstr, userID, projectID))
		if err != nil {
			log.Fatal(err)
		}
		var dat []map[string]interface{}
		err = json.Unmarshal(content, &dat)
		if err != nil {
			log.Fatal(err)
		}
		if len(dat) > 0 && dat[0]["final_mark"] != nil {

			fmt.Println(fmt.Sprintf("%s: %f", projectName, dat[0]["final_mark"]))
		}
		if projectID == 5 { // id of ft_printf
			os.Exit(0)
		}
		time.Sleep(time.Second) // sleep for api rate limit
	}

}
