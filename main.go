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
		fmt.Println("usage: ./automatic-octo-computing-machine api-uid api-secret user")
		os.Exit(1)
	}

	uid := os.Args[1]
	secret := os.Args[2]
	username := os.Args[3]
	client := apiwrap.NewClient(uid, secret)
	fmtstr := "/v2/users/%d/projects_users?filter[project_id]=%d"

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
		if projectID == 5 {
			os.Exit(0)
		}
		time.Sleep(time.Second) // sleep for api rate limit
	}

}
