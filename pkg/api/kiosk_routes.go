package api

import (
	"html/template"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (a *App) handleKioskSetup(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("kiosk test page accessed", zap.Time("at", time.Now()))

	files := []string{
		"./ui/html/base.html",
		"./ui/html/kiosk.page.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.logger.Error("failed parsing templates", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		a.logger.Error("failed executing template", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
