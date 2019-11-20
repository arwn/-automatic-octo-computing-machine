package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arwn/apiwrap"
)

type user struct {
	ID    int
	Login string
}

type project struct {
	ID   int
	Name string
}

func projectIterator() func() project {
	projects := []string{"ft_ls", "fdf", "ft_printf"}
	counter := 2 // init to two, first project has an ID of 3
	return func() project {
		var p project
		if counter == 5 {
			counter = 2
		}
		counter++
		p.ID = counter
		p.Name = projects[counter-3]
		return p
	}
}

func userIterator(c *apiwrap.WrapperClient) func() (bool, user) {
	var page = 0
	var users []user

	return func() (bool, user) {
		if len(users) == 0 {
			log.Println("getting new users")
			body, err := c.GetBody(fmt.Sprintf("/v2/campus/7/users?page=%d", page))
			if err != nil {
				log.Println(err)
			}
			// no more users, return that we are done
			if len(body) == 0 {
				return true, user{}
			}

			err = json.Unmarshal(body, &users)
			if err != nil {
				log.Println(err)
				users = nil
				return false, user{}
			}

			for _, user := range users {
				users = append(users, user)
			}
			page++
		}

		// drop the first user and return it
		var u user
		u, users = users[0], users[1:]
		return false, u
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	uid := os.Getenv("user_id")
	secret := os.Getenv("user_secret")
	client := apiwrap.NewClient(uid, secret)
	fmtstr := "/v2/users/%d/projects_users?filter[project_id]=%d&page[size]=100"
	client.Timeout = time.Second / 2

	// for every user
	user := userIterator(&client)
	project := projectIterator()
	for {
		done, u := user()
		if done {
			break
		}
		log.Printf("checking user %v\n", u)

		for {
			p := project()
			content, err := client.GetJSON(fmt.Sprintf(fmtstr, u.ID, p.ID))
			if err != nil {
				log.Println(err)
			}
			var dat []map[string]interface{}
			err = json.Unmarshal(content, &dat)
			if err != nil {
				log.Println(err)
			}
			if len(dat) > 0 && dat[0]["final_mark"] != nil {
				log.Printf("%s %s: %1f", u.Login, p.Name, dat[0]["final_mark"])
			}
			if p.ID == 5 { // ID of printf, the last project we check for each user
				break
			}
		}
	}

}
