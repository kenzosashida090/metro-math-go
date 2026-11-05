package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kenzosashida090/metro-math-go/store"
)

type SectionPayload struct {
	Groups []int `json:"groups" validate:"required,min=1,max=16,dive,min=1,max=16"`
	Bars   int   `json:"bars" validate:"required,min=1,max=64"`
}

type CreateSongPayload struct {
	Name     string           `json:"name" validate:"required,min=1,max=80"`
	BPM      int              `json:"bpm" validate:"required,min=30,max=240"`
	Sections []SectionPayload `json:"sections" validate:"required,min=1,max=50,dive"`
}

type UpdateSongPayload struct {
	Name     string           `json:"name" validate:"omitempty,min=1,max=80"`
	BPM      int              `json:"bpm" validate:"omitempty,min=30,max=240"`
	Sections []SectionPayload `json:"sections" validate:"omitempty,min=1,max=50,dive"`
}

// Songs OLANGdoc's
//
// @Summary		Create song
// @Description	Create a song
// @Tags		songs
// @Accept			json
// @Produce		json
// @Param			body	body CreateSongPayload true "Song payload"
// @Success		200		{object}	store.Song
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Router			/songs [POST]
func (app *application) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateSongPayload

	err := ReadJSON(w, r, &payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "missing data")
		return
	}

	err = Validate.Struct(payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "missing data")
		return
	}
	sections := make([]store.Section, len(payload.Sections))
	for i, s := range payload.Sections {
		sections[i] = store.Section{Groups: s.Groups, Bars: s.Bars}
	}
	song := &store.Song{
		Name:     payload.Name,
		BPM:      payload.BPM,
		Sections: sections,
	}
	err = app.store.Songs.CreateSong(r.Context(), song)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if err := app.responseJSON(w, http.StatusOK, song); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Something wnt wribg")
	}
}

// Songs OLANGdoc's
//
// @Summary		Update song
// @Description	Update a song
// @Tags		songs
// @Accept			json
// @Produce		json
// @Param id path int true "songId"
// @Param			body	body UpdateSongPayload false "Song payload"
// @Success		200		{object}	store.Song
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Router			/songs/{id} [PATCH]
func (app *application) UpdateSongByIdHandler(w http.ResponseWriter, r *http.Request) {
	var updatePayload UpdateSongPayload
	songId := chi.URLParam(r, "songId")
	fmt.Println(songId, "songid")
	id, err := strconv.ParseInt(songId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = ReadJSON(w, r, &updatePayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(updatePayload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if updatePayload.Name == "" && updatePayload.BPM == 0 && len(updatePayload.Sections) == 0 {
		WriteJSONError(w, http.StatusBadRequest, "nothing to update")
		return
	}
	sections := make([]store.Section, len(updatePayload.Sections))
	for i, s := range updatePayload.Sections {
		sections[i] = store.Section{Groups: s.Groups, Bars: s.Bars}
	}
	song := &store.Song{
		Name:     updatePayload.Name,
		BPM:      updatePayload.BPM,
		Sections: sections,
	}
	err = app.store.Songs.UpdateSong(r.Context(), id, song)
	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}
	if err := app.responseJSON(w, http.StatusOK, song); err != nil {
		app.intertalServerError(w, r, err)
	}

}

// Songs OLANGdoc's
//
// @Summary		Get song by id
// @Description	Song bY ID
// @Tags		songs
// @Accept			json
// @Produce		json
// @Param id path int true "songId"
// @Success		200		{object}	store.Song
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Router			/songs/{id} [GET]
func (app *application) GetSongByIdHandler(w http.ResponseWriter, r *http.Request) {
	songId := chi.URLParam(r, "songId")
	id, err := strconv.ParseInt(songId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	song, err := app.store.Songs.GetSongById(r.Context(), id)
	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}
	if err := app.responseJSON(w, http.StatusOK, song); err != nil {
		app.intertalServerError(w, r, err)
		return
	}
}

// Songs OLANGdoc's
//
// @Summary		GET ALL songs
// @Description	GET all  songs
// @Tags		songs
// @Accept			json
// @Produce		json
// @Success		200		{object}	store.Song
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Router			/songs [GET]
func (app *application) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := app.store.Songs.GetSongsByMetadata(r.Context())
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if err := app.responseJSON(w, http.StatusOK, songs); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Something went wrong")
	}
}

// Songs OLANGdoc's
//
// @Summary		DELETE A SONG
// @Description	Delete a song
// @Tags		songs
// @Accept			json
// @Produce		json
// @Param id path int true "songId"
// @Success		200		{object}	store.Song
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Router			/songs/{id} [DELETE]
func (app *application) DeleteSongByIdHandler(w http.ResponseWriter, r *http.Request) {
	songId := chi.URLParam(r, "songId")
	id, err := strconv.ParseInt(songId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = app.store.Songs.DeleteSong(r.Context(), id)
	if err != nil {
		app.intertalServerError(w, r, err)
		return
	}

	if err := app.responseJSON(w, http.StatusOK, "ok"); err != nil {
		app.intertalServerError(w, r, err)
	}
}
