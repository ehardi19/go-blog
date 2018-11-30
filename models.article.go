package main

import (
	"database/sql"
)

type article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *article) getArticle(db *sql.DB) error {
	return db.QueryRow("SELECT title, content FROM article WHERE id=$1", a.ID).Scan(&a.Title, &a.Content)
}

func (a *article) updateArticle(db *sql.DB) error {
	_, err := db.Exec("UPDATE article SET title=$1, content=$2 WHERE id=$3", a.Title, a.Content, a.ID)

	return err
}

func (a *article) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM article WHERE id=$1", a.ID)

	return err
}

func (a *article) createArticle(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO article(title, content) VALUES ($1, $2) RETURNING id", a.Title, a.Content).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

