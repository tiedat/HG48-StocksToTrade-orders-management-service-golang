package main

import (
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type User struct {
}

func (s *server) avatarHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		logrus.WithError(err).Error("Can't get file from request")
		s.writeResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		logrus.WithError(err)
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	out, err := os.Create("uploads/user/avatar/" + "name" + ".jpg")
	if err != nil {
		logrus.WithError(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		logrus.WithError(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
