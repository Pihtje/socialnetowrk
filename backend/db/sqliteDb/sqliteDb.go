package sqliteDb

import (
	"database/sql"
	"fmt"
	"forum/exec"
	"log"
	"os"
	"strconv"

	migrate "github.com/rubenv/sql-migrate"
)

// creates the database
func Initialise() {
	fmt.Println("creating database...")
	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal("CreateDatabase: ", err)
	}
	file.Close()
	exec.Db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database created successfully!")
}

// applies all Up migrations
func ApplyMigrationsUp() {
	fmt.Println("applying up migrations...")
	Db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal("1: ", err)
	}
	defer Db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migrations/sqlite",
	}

	_, err = migrate.Exec(Db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Fatal("ApplyMigrationsUp: ", err)
	}

	fmt.Println("up migrations completed successfully!")
}

// Populates the database with already existing users.
func PopulateUsers() {
	fmt.Println("populating users...")
	for i := 1; i < 6; i++ {
		num := strconv.Itoa(i)
		exec.Signup("user"+num+"@mail.com", "password", "user"+num, "Smith", "01/01/1970", "placeholder", "user"+num, "", "private")
	}
	fmt.Println("done populating users!")

	exec.UpdateUserVisibility(3, "public")
	exec.UpdateUserVisibility(5, "public")

	exec.UpdateOptionalUserdata("/images/user_images/user_pic_2_73cf3c9d.jpg", "TheDude", "there is nothing here for you, get lost!", "2")

	_, err := exec.Db.Exec("INSERT INTO images (user_id, image_URL) VALUES (?,?);", 2, "/images/user_images/user_pic_2_73cf3c9d.jpg")
	if err != nil {
		fmt.Println("PopulateUsers err 1: ", err)
		return
	}

	exec.UpdateOptionalUserdata("", "fdsaguy", "", "3")
	exec.UpdateOptionalUserdata("/images/user_images/user_pic_5_84f20d8d.jpg", "asdfguy", "", "5")
	_, err = exec.Db.Exec("INSERT INTO images (user_id, image_URL) VALUES (?,?);", 5, "/images/user_images/user_pic_5_84f20d8d.jpg")
	if err != nil {
		fmt.Println("PopulateUsers err 1: ", err)
		return
	}
}

// Populates the database with already existing posts.
func InitToPosts() {
	exec.Post(1, 1, "What's the deal with airline food?", "Like fr!1", []int{}, 0)
	exec.SaveImagePath(1, "/images/post_images/post_pic_1_c49e6e2b.jpg")

	exec.Post(2, 1, "IKR", "xddxd2", []int{}, 0)
	exec.Post(3, 1, "Please post in Off-topic", "3", []int{}, 0)
	exec.SaveImagePath(3, "/images/post_images/post_pic_3_95edbecd.jpg")

	exec.Post(2, 2, "Guys check out my jokes in Off-topic!!!!", "They are about chickens ;)", []int{}, 0)
	exec.Post(5, 2, "Can we get some moderators in this forum?", "People keep posting randomly in every category", []int{}, 0)
	exec.SaveImagePath(5, "/images/post_images/post_pic_5_c8d59a4c.jpg")

	exec.Post(2, 3, "Why did the chicken cross the road?", "To get to the other side", []int{1, 4}, 0)
	exec.SaveImagePath(9, "/images/post_images/post_pic_6_d67c9ec8.jpg")

	exec.Post(2, 3, "Why did the chicken cross the playground?", "To get to the other side", []int{1}, 0)
	exec.Post(4, 3, "Why did the chicken cross the playground?", "To get to the other side", []int{}, 0)

	exec.Post(1, 1, "test post 1 with image", "test post 1 with image body", []int{}, 0)
	exec.SaveImagePath(9, "/images/post_images/post_pic_9_5a894206.png")

	exec.Post(1, 1, "test post in group 1", "test post in group 1 body", []int{}, 1)
}

// Adds images to some already existing comments.
func InitToCommentImages() {
	exec.SaveCommentImagePath(1, "/images/comment_images/comment_pic_1_0eb1e3ec.jpg")
	exec.SaveCommentImagePath(3, "/images/comment_images/comment_pic_3_1e87f8ce.jpg")
	exec.SaveCommentImagePath(5, "/images/comment_images/comment_pic_5_d67c9ec8.jpg")
	exec.SaveCommentImagePath(6, "/images/comment_images/comment_pic_6_c49e6e2b.jpg")
}

func CreateGroups() {
	groupId1, _ := exec.CreateNewGroup(1, "test group 1", "group for testing")
	exec.AddGroupMember(fmt.Sprint(groupId1[0]), "2")
	exec.AddGroupMember(fmt.Sprint(groupId1[0]), "4")

	groupId2, _ := exec.CreateNewGroup(1, "test group 2", "group for more testing")
	exec.AddGroupMember(fmt.Sprint(groupId2[0]), "5")

	exec.CreateNewGroup(1, "test group 3", "group for testing if no other members besides creator")
}
