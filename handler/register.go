package handler

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"lheinrich.de/secpass/conf"
	"lheinrich.de/secpass/shorts"
)

// Register function
func Register(w http.ResponseWriter, r *http.Request) {
	// check logged in and redirect to index if so
	if checkSession(w, r) != "" {
		redirect(w, "/")
		return
	}

	// output text if equals special
	special := 0

	// define values
	name, password, repeat := r.PostFormValue("name"), r.PostFormValue("password"), r.PostFormValue("repeat")

	// check for input
	if name != "" && password != "" && repeat != "" && len(password) >= 8 && len(repeat) >= 8 {
		// check whether passwords match
		if password == repeat {
			// check whether name already exists
			errQuery := conf.DB.QueryRow("SELECT password FROM users WHERE name = $1", name).Scan(nil)

			// name does not exist
			if errQuery == sql.ErrNoRows {
				// generate random key
				key := shorts.Encrypt(shorts.UUID(), shorts.GenerateKey(shorts.UUID()))

				// encrypt key
				encryptedKey := shorts.Encrypt(key, shorts.GenerateKey(password))

				// hash password and insert user
				passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+1)
				_, errExec := conf.DB.Exec(conf.GetSQL("register"), name, string(passwordHash), "", encryptedKey)
				shorts.Check(errExec)

				// redirect and return
				redirect(w, "/login")
				return
			}

			// name exists, print
			special = 1
		} else {
			// passwords does not match, print
			special = 2
		}
	}

	// execute template
	shorts.Check(tpl.ExecuteTemplate(w, "register.html", Data{User: "", Lang: getLang(r), Special: special, LoggedOut: true}))
}
