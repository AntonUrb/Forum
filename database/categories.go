package database

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"strconv"
)

// CreateCategoriesTable creates the Categories table in the database.
func CreateCategoriesTable(db *sql.DB) error {
	categoriesTable := `CREATE TABLE Categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,	
		name TEXT
	  );`

	log.Println("Creating Categories table...")
	statement, err := db.Prepare(categoriesTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return nil
}

// InsertCategory inserts predefined categories into the Categories table.
func InsertCategory(db *sql.DB) {
	cat := []string{"General", "Nature", "Sports", "Technology"}
	for i := 0; i < len(cat); i++ {
		fmt.Println("Adding to categories:", cat[i])
		categorySQL := "INSERT INTO Categories(name) VALUES (?)"
		statement, _ := db.Prepare(categorySQL)
		statement.Exec(cat[i])
	}
}

// GetCategoryPosts retrieves posts belonging to a specific category.
func GetCategoryPosts(db *sql.DB, category string) ([]models.Post, error) {
	// Get category ID from input category name
	q := "SELECT * FROM Categories WHERE name = ?"
	rows, err := db.Query(q, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryId int
	for rows.Next() {
		var name string
		err = rows.Scan(&categoryId, &name)
		if err != nil {
			return nil, err
		}
	}

	// Query PostCategories table to get post IDs associated with the category
	q = "SELECT * FROM PostCategories WHERE category_id = ?"
	rows, err = db.Query(q, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryPostIds []int
	for rows.Next() {
		var id, posts_id, category_id int
		err = rows.Scan(&id, &posts_id, &category_id)
		if err != nil {
			return nil, err
		}
		categoryPostIds = append(categoryPostIds, posts_id)
	}

	// Retrieve the posts using their IDs
	var categoryPosts []models.Post
	for _, id := range categoryPostIds {
		idAsStr := strconv.Itoa(id)

		var post models.Post
		if res, err := GetPost(db, idAsStr); err != nil {
			return nil, err
		} else {
			post = res
		}
		categoryPosts = append(categoryPosts, post)
	}
	return categoryPosts, err
}

// CreatePostCategoriesTable creates the PostCategories table in the database.
func CreatePostCategoriesTable(db *sql.DB) error {
	postCategoriesTable := `CREATE TABLE PostCategories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,	
		posts_id 	INTEGER references postsTable(id),
		category_id INTEGER references categoriesTable(id)
	  );`

	log.Println("Creating PostCategories table...")
	statement, err := db.Prepare(postCategoriesTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return nil
}

// InsertPostCategory inserts a mapping between a post and a category.
func InsertPostCategory(db *sql.DB, postId, categoryId int) error {
	categorySQL := "INSERT INTO PostCategories(posts_id, category_id) VALUES (?,?)"
	statement, err := db.Prepare(categorySQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(postId, categoryId)
	if err != nil {
		return err
	}
	return nil
}
