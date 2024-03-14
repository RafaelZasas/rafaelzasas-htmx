package apiserver

import (
	"database/sql"
	"fmt"
	"html/template"
	"htmx-go/database"
	"htmx-go/models"
	"htmx-go/utils"
	"log"
	"net/http"
	"os"
	"time"
)

func (api *api) bindAuthRoutes() {
	// Authentication
	api.router.POST("/login", api.handleLogin)
	api.router.POST("/register", api.handleRegister)
	api.router.POST("/api/auth/refresh-token", api.handleTokenRefresh)
	api.router.POST("/api/auth/logout", api.handleLogout)
}

func (api *api) handleLogin(w http.ResponseWriter, r *http.Request) {
	data := api.basePageData(r.Context())

	db, err := api.DbFromContext(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl := api.newTemplate("views/login.html")

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := db.GetUserByEmail(email)

	data.Title = "Login | Go-Htmx Example"

	if err == sql.ErrNoRows {
		errHtml := template.HTML(`
        Email Not Found.&nbsp;
        <a 
        class="hover:text-sky-500 text-sky-600 font-semibold"
        href='/register' hx-boost="false"
        >Register</a>&nbsp; instead?`)

		data.Extra["ErrorMessage"] = errHtml
		w.WriteHeader(http.StatusUnauthorized)
		tmpl.Render(w, "errors", data)
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("(api *api) handleLogin: %s", err))
		data.Extra["ErrorMessage"] = "An error occurred. Please try again later."
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Render(w, "errors", data)
		return
	}

	if !utils.CheckPasswordHash(password, *user.Password) {
		errHtml := template.HTML(`
        Incorrect Password.&nbsp;
        <a 
        class="hover:text-sky-500 text-sky-600 font-semibold"
        href='/forgot-password' hx-boost="false"
        >Forgot Password</a>?`)
		data.Extra["ErrorMessage"] = errHtml
		w.WriteHeader(http.StatusUnauthorized)
		tmpl.Render(w, "errors", data)
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWTokens(user.UID)
	if err != nil {
		log.Println(fmt.Errorf("could not generate JWT: %s", err))
		data.Extra["ErrorMessage"] = "An error occurred. Please try again later."
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Render(w, "errors", data)
		return
	}

	// Update users refresh token
	err = db.UpdateUserRefreshToken(user.UID, refreshToken)
	if err != nil {
		log.Println(fmt.Errorf("could not update refresh token: %s", err))
		data.Extra["ErrorMessage"] = "An error occurred. Please try again later."
		w.WriteHeader(http.StatusInternalServerError)
		tmpl.Render(w, "errors", data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
	})

	// add redirect header and return successful response
	redirect := r.URL.Query().Get("redirect")
	w.Header().Add("HX-Redirect", redirect)
	w.WriteHeader(http.StatusOK)
	return

}

func (api *api) handleRegister(w http.ResponseWriter, r *http.Request) {
	data := api.basePageData(r.Context())

	db, err := api.DbFromContext(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl := api.newTemplate("views/register.html")
	extraData := make(map[string]interface{})

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")

	// Check if user already exists
	_, err = db.GetUserByEmail(email)
	if err == nil {
		errHtml := template.HTML(`
        Email already exists.&nbsp;
        <a
        class="hover:text-sky-500 text-sky-600 font-semibold"
        href='/login' hx-boost="false"
        >Login</a>&nbsp; instead?`)
		data.Extra["ErrorMessage"] = errHtml
		w.WriteHeader(http.StatusConflict)
		tmpl.Render(w, "errors", data)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println(fmt.Errorf("could not hash password: %s", err))
		extraData["ErrorMessage"] = "An error occurred. Please try again later."
		data.Extra = extraData
		tmpl.Render(w, "errors", data)
		return
	}

	uid, err := utils.GenerateUID()
	if err != nil {
		log.Println(fmt.Errorf("could not generate UUID: %s", err))
		extraData["ErrorMessage"] = "An error occurred. Please try again later."
		data.Extra = extraData
		tmpl.Render(w, "errors", data)
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWTokens(uid)

	if err != nil {
		log.Println(fmt.Errorf("could not generate JWT: %s", err))
		extraData["ErrorMessage"] = "An error occurred. Please try again later."
		data.Extra = extraData
		tmpl.Render(w, "errors", data)
		return
	}

	user := &models.User{
		UID:           uid,
		Email:         email,
		Password:      &hashedPassword,
		Name:          name,
		Provider:      "email",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		RefreshToken:  &refreshToken,
	}

	err = db.CreateUser(user)
	if err != nil {
		log.Println(fmt.Errorf("could not create user: %s", err))
		extraData["ErrorMessage"] = "An error occurred. Please try again later."
		data.Extra = extraData
		tmpl.Render(w, "errors", data)
		return
	}

	utils.AddAuthCookies(&w, accessToken, refreshToken)

	w.Header().Add("HX-Redirect", "/profile")
	w.WriteHeader(http.StatusCreated)
	return
}

func (api *api) handleTokenRefresh(w http.ResponseWriter, r *http.Request) {

	db, err := api.DbFromContext(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	uid, err := utils.ValidateRefreshToken(r)

	if err != nil {
		log.Println(fmt.Errorf("could not validate refresh token: %s", err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userRefreshToken, err := db.GetUserRefreshToken(uid)
	if err != nil {
		log.Println(fmt.Errorf("could not get refresh token: %s", err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	providedRefreshToken, _ := r.Cookie("refresh_token")

	if providedRefreshToken.Value != userRefreshToken {
		log.Println(fmt.Errorf("refresh token does not match"))
		fmt.Println(userRefreshToken)
		fmt.Println(providedRefreshToken.Value)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWTokens(uid)
	if err != nil {
		log.Println(fmt.Errorf("could not generate JWT: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update users refresh token
	err = db.UpdateUserRefreshToken(uid, refreshToken)
	if err != nil {
		log.Println(fmt.Errorf("could not update refresh token: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.AddAuthCookies(&w, accessToken, refreshToken)

	w.WriteHeader(http.StatusOK)
}

func (api *api) handleLogout(w http.ResponseWriter, r *http.Request) {

	uid, _ := utils.ValidateRefreshToken(r)

	db := r.Context().Value("db").(database.Database)

	if uid != "" {
		err := db.UpdateUserRefreshToken(uid, "")
		if err != nil {
			log.Println(fmt.Errorf("(api.handleLogout: %s", err))
		}
	}

	// Expire the access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1, // Instructs the browser to delete the cookie immediately.
		HttpOnly: true,
		Secure:   true,
	})

	// Expire the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1, // Instructs the browser to delete the cookie immediately.
		HttpOnly: true,
		Secure:   true,
	})

	w.Header().Add("HX-Refresh", "true")
	// Optionally, return a response indicating success
	w.WriteHeader(http.StatusResetContent)
}
