package main

import "snippetbox.purple-mountain.gg/internal/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
