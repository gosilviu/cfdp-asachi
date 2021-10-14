package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var tpl *template.Template
var db *sql.DB
var store = sessions.NewCookieStore([]byte("super-secret"))

func main() {

	tpl, _ = template.ParseGlob("templates/*.html")
	var err error
	db, err = sql.Open("mysql", "root:Federerthebestcasa2007!!!@tcp(127.0.0.1:3306)/mygodatabase")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/addbridge", addBridge)
	http.HandleFunc("/fisa_stare_tehnica", fisaStareTehnica)
	http.HandleFunc("/addbridgehandler", addBridgeHandler)
	http.ListenAndServe(":8080", nil)
}

func fisaStareTehnica(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****fisa stare tehnica running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	tpl.ExecuteTemplate(w, "fisa_stare_tehnica.html", session.Values["user"])
}
func addBridge(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****add bridge running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	tpl.ExecuteTemplate(w, "addbridge.html", session.Values["user"])
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****logoutHandler running*****")
	session, _ := store.Get(r, "session")
	// The delete built-in function deletes the element with the specified key (m[key]) from the map.
	// If m is nil or there is no such element, delete is a no-op.
	delete(session.Values, "userID")
	session.Save(r, w)
	tpl.ExecuteTemplate(w, "login.html", "Logged Out")
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****indexHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	tpl.ExecuteTemplate(w, "index.html", session.Values["user"])
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****aboutHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	tpl.ExecuteTemplate(w, "about.html", session.Values["user"])
}
func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****contactHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	tpl.ExecuteTemplate(w, "contact.html", session.Values["user"])
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "register.html", nil)
}
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	var nameAlphaNumeric = true
	for _, char := range username {
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			nameAlphaNumeric = false
		}
	}
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	password := r.FormValue("password")
	fmt.Println("password:", password, "\npswdLength:", len(password))
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			pswdLowercase = true
		case unicode.IsUpper(char):
			pswdUppercase = true
		case unicode.IsNumber(char):
			pswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdLength || !pswdNoSpaces || !pswdNumber || !pswdUppercase || !pswdSpecial || !nameAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}
	stmt := "SELECT id FROM register WHERE user = ?"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username aleready exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem resiterging account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO register (user, pass) VALUES (?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()
	fmt.Println("--------------------")
	var result sql.Result

	result, err = insertStmt.Exec(username, hash)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rowsAff, _ := result.RowsAffected()
	fmt.Println("--------------------")
	lastIns, _ := result.LastInsertId()

	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Fprint(w, "congrats, your account has been successfully created")
}
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password: ", password)
	var userID, hash string
	stmt := "SELECT id, pass FROM register WHERE user = ?"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&userID, &hash)
	fmt.Println("hash from db: ", hash)
	if err != nil {
		fmt.Println("error selecting hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		session, _ := store.Get(r, "session")
		session.Values["userID"] = userID
		session.Values["user"] = username
		session.Save(r, w)
		tpl.ExecuteTemplate(w, "index.html", "Logged In")
		return
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
}

func addBridgeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****addBridgeHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "fisa_stare_tehnica.html", nil)
		return
	}

	r.ParseForm()
	tip_lucrare := r.FormValue("tip_lucrare")
	obstacol_traversat := r.FormValue("obstacol_traversat")
	localitate := r.FormValue("localitate")
	categoria_drum := r.FormValue("categoria_drum")
	poz_km := r.FormValue("poz_km")
	an_consolidat := r.FormValue("an_consolidat")
	tip_pod := r.FormValue("tip_pod")
	material := r.FormValue("material")
	lungime := r.FormValue("lungime")
	latime := r.FormValue("latime")
	reazem := r.FormValue("reazem")
	infra := r.FormValue("infra")
	fundatii := r.FormValue("fundatii")
	imbracaminte := r.FormValue("imbracaminte")
	rosturi := r.FormValue("rosturi")
	pozitie := r.FormValue("pozitie")
	parapeti_pietonali := r.FormValue("parapeti_pietonali")
	parapeti_siguranta := r.FormValue("parapeti_siguranta")
	racordari := r.FormValue("racordari")
	aparari := r.FormValue("aparari")

	_1c1 := r.FormValue("1c1")
	_1c2 := r.FormValue("1c2")
	_2c4 := r.FormValue("2c4")
	_3c5 := r.FormValue("3c5")
	_4c3 := r.FormValue("4c3")
	_5c3 := r.FormValue("5c3")
	_6c1 := r.FormValue("6c1")
	_6c2 := r.FormValue("6c2")
	_6c3 := r.FormValue("6c3")
	_7c1 := r.FormValue("7c1")
	_7c2 := r.FormValue("7c2")
	_7c3 := r.FormValue("7c3")
	_8c1 := r.FormValue("8c1")
	_8c2 := r.FormValue("8c2")
	_8c3 := r.FormValue("8c3")
	_9c1 := r.FormValue("9c1")
	_9c2 := r.FormValue("9c2")
	_9c3 := r.FormValue("9c3")
	_10c1 := r.FormValue("10c1")
	_11c5 := r.FormValue("11c5")
	_12c1 := r.FormValue("12c1")
	_12c2 := r.FormValue("12c2")
	_12c3 := r.FormValue("12c3")
	_13c5 := r.FormValue("13c5")
	_14c1 := r.FormValue("14c1")
	_14c2 := r.FormValue("14c2")
	_14c3 := r.FormValue("14c3")
	_15c1 := r.FormValue("15c1")
	_15c2 := r.FormValue("15c2")
	_16c1 := r.FormValue("16c1")
	_16c2 := r.FormValue("16c2")
	_16c3 := r.FormValue("16c3")
	_17c1 := r.FormValue("17c1")
	_17c2 := r.FormValue("17c2")
	_17c3 := r.FormValue("17c3")
	_18c1 := r.FormValue("18c1")
	_18c2 := r.FormValue("18c2")
	_19c1 := r.FormValue("19c1")
	_20c5 := r.FormValue("20c5")
	_21c5 := r.FormValue("21c5")
	_22c4 := r.FormValue("22c4")
	_23c4 := r.FormValue("23c4")
	_24c5 := r.FormValue("24c5")
	_25c3 := r.FormValue("25c3")
	_26c2 := r.FormValue("26c2")
	_27c1 := r.FormValue("27c1")
	_28c1 := r.FormValue("28c1")
	_29c3 := r.FormValue("29c3")
	_30c3 := r.FormValue("30c3")
	_31c2 := r.FormValue("31c2")
	_33c3 := r.FormValue("33c3")
	_34c1 := r.FormValue("34c1")
	_34c2 := r.FormValue("34c2")
	_35c1 := r.FormValue("35c1")
	_35c2 := r.FormValue("35c2")
	_35c3 := r.FormValue("35c3")
	_36c1 := r.FormValue("36c1")
	_36c2 := r.FormValue("36c2")
	_36c3 := r.FormValue("36c3")
	_37c1 := r.FormValue("37c1")
	_37c2 := r.FormValue("37c2")
	_37c3 := r.FormValue("37c3")
	_38c5 := r.FormValue("38c5")
	_39c1 := r.FormValue("39c1")
	_40c1 := r.FormValue("40c1")
	_40c2 := r.FormValue("40c2")
	_41c1 := r.FormValue("41c1")
	_41c2 := r.FormValue("41c2")
	_42c5 := r.FormValue("42c5")
	_43c3 := r.FormValue("43c3")
	_44c1 := r.FormValue("44c1")
	_44c2 := r.FormValue("44c2")
	_44c3 := r.FormValue("44c3")
	_45c1 := r.FormValue("45c1")
	_46c5 := r.FormValue("46c5")
	_47c4 := r.FormValue("47c4")
	_48c5 := r.FormValue("48c5")
	_49c1 := r.FormValue("49c1")
	_49c2 := r.FormValue("49c2")
	_50c5 := r.FormValue("50c5")
	_51c5 := r.FormValue("51c5")
	_52c3 := r.FormValue("52c3")
	_53c4 := r.FormValue("53c4")
	_54c1 := r.FormValue("54c1")
	_54c3 := r.FormValue("54c3")
	_55c4 := r.FormValue("55c4")
	_56c1 := r.FormValue("56c1")
	_57c1 := r.FormValue("57c1")
	_57c2 := r.FormValue("57c2")
	_58c3 := r.FormValue("58c3")
	_59c3 := r.FormValue("59c3")
	_60c1 := r.FormValue("60c1")
	_60c2 := r.FormValue("60c2")
	_61c5 := r.FormValue("61c5")
	_62c1 := r.FormValue("62c1")
	_62c2 := r.FormValue("62c2")
	_63c5 := r.FormValue("63c5")
	_64c1 := r.FormValue("64c1")
	_64c3 := r.FormValue("64c3")
	_65c5 := r.FormValue("65c5")
	_66c5 := r.FormValue("66c5")
	_67c1 := r.FormValue("67c1")
	_67c2 := r.FormValue("67c2")
	_67c3 := r.FormValue("67c3")
	_68c1 := r.FormValue("68c1")
	_68c2 := r.FormValue("68c2")
	_68c3 := r.FormValue("68c3")
	_69c4 := r.FormValue("69c4")
	_70c1 := r.FormValue("70c1")
	_70c2 := r.FormValue("70c2")
	_71c1 := r.FormValue("71c1")
	_71c3 := r.FormValue("71c3")
	_72c1 := r.FormValue("72c1")
	_72c3 := r.FormValue("72c3")
	_73c3 := r.FormValue("73c3")
	_74c1 := r.FormValue("74c1")
	_74c2 := r.FormValue("74c2")
	_74c3 := r.FormValue("74c3")
	_75c1 := r.FormValue("75c1")
	_76c1 := r.FormValue("76c1")
	_77c1 := r.FormValue("77c1")
	_78c1 := r.FormValue("78c1")
	_79c1 := r.FormValue("79c1")
	_80c1 := r.FormValue("80c1")
	_81c2 := r.FormValue("81c2")
	_82c2 := r.FormValue("82c2")
	_83c2 := r.FormValue("83c2")
	_84c3 := r.FormValue("84c3")
	_85c3 := r.FormValue("85c3")
	_86c3 := r.FormValue("86c3")
	_87c3 := r.FormValue("87c3")
	_88c3 := r.FormValue("88c3")
	_89c3 := r.FormValue("89c3")
	_90c3 := r.FormValue("90c3")
	_91c5 := r.FormValue("91c5")
	_92c5 := r.FormValue("92c5")
	_93c5 := r.FormValue("93c5")
	_94c5 := r.FormValue("94c5")
	_95c5 := r.FormValue("95c5")
	_96c5 := r.FormValue("96c5")
	_97c5 := r.FormValue("97c5")

	stmt := "SELECT id FROM bridges WHERE user = ?"
	row := db.QueryRow(stmt, tip_lucrare)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		//do nothing
	}

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO bridges (user, tip_lucrare, obstacol, localitatea, categorie, pozitie_km, an_constructie, tip_pod, materialul, lung_pod, latime_pod, reazem, infrastructura, tip_fundatii, imbracaminte, rosturi, pozitie, parapet_pietonal, parapeti_siguranta, racordari_terasamente, aparari_mal, 1c1, 1c2, 2c4, 3c5, 4c3, 5c3, 6c1, 6c2, 6c3, 7c1, 7c2, 7c3, 8c1, 8c2, 8c3, 9c1, 9c2, 9c3, 10c1, 11c5, 12c1,12c2, 12c3, 13c5, 14c1, 14c2, 14c3, 15c1, 15c2, 16c1, 16c2, 16c3, 17c1, 17c2, 17c3, 18c1, 18c2, 19c1, 20c5, 21c5, 22c4, 23c4, 24c5, 25c3, 26c2, 27c1, 28c1, 29c3, 30c3, 31c2, 33c3, 34c1, 34c2, 35c1, 35c2, 35c3, 36c1,36c2, 36c3, 37c1, 37c2, 37c3, 38c5, 39c1, 40c1, 40c2, 41c1, 41c2, 42c5,43c3, 44c1, 44c2, 44c3, 45c1, 46c5, 47c4, 48c5, 49c1, 49c2, 50c5, 51c5,52c3, 53c4, 54c1, 54c3, 55c4, 56c1, 57c1, 57c2, 58c3, 59c3, 60c1, 60c2, 61c5, 62c1, 62c2, 63c5, 64c1, 64c3, 65c5, 66c5, 67c1, 67c2, 67c3, 68c1, 68c2, 68c3, 69c4, 70c1, 70c2, 71c1, 71c3, 72c1, 72c3, 73c3, 74c1, 74c2, 74c3, 75c1, 76c1, 77c1, 78c1, 79c1, 80c1, 81c2, 82c2, 83c2, 84c3, 85c3, 86c3, 87c3, 88c3, 89c3, 90c3, 91c5, 92c5, 93c5, 94c5, 95c5, 96c5, 97c5) VALUES (?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?, ?, ?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "fisa_stare_tehnica.html", "there was a problem registering bridge")
		return
	}
	defer insertStmt.Close()
	fmt.Println("--------------------")
	var result sql.Result

	result, err = insertStmt.Exec(
		session.Values["user"],
		tip_lucrare,
		obstacol_traversat,
		localitate,
		categoria_drum,
		poz_km,
		an_consolidat,
		tip_pod,
		material,
		lungime,
		latime,
		reazem,
		infra,
		fundatii,
		imbracaminte,
		rosturi,
		pozitie,
		parapeti_pietonali,
		parapeti_siguranta,
		racordari,
		aparari,
		_1c1,
		_1c2,
		_2c4,
		_3c5,
		_4c3,
		_5c3,
		_6c1,
		_6c2,
		_6c3,
		_7c1,
		_7c2,
		_7c3,
		_8c1,
		_8c2,
		_8c3,
		_9c1,
		_9c2,
		_9c3,
		_10c1,
		_11c5,
		_12c1,
		_12c2,
		_12c3,
		_13c5,
		_14c1,
		_14c2,
		_14c3,
		_15c1,
		_15c2,
		_16c1,
		_16c2,
		_16c3,
		_17c1,
		_17c2,
		_17c3,
		_18c1,
		_18c2,
		_19c1,
		_20c5,
		_21c5,
		_22c4,
		_23c4,
		_24c5,
		_25c3,
		_26c2,
		_27c1,
		_28c1,
		_29c3,
		_30c3,
		_31c2,
		_33c3,
		_34c1,
		_34c2,
		_35c1,
		_35c2,
		_35c3,
		_36c1,
		_36c2,
		_36c3,
		_37c1,
		_37c2,
		_37c3,
		_38c5,
		_39c1,
		_40c1,
		_40c2,
		_41c1,
		_41c2,
		_42c5,
		_43c3,
		_44c1,
		_44c2,
		_44c3,
		_45c1,
		_46c5,
		_47c4,
		_48c5,
		_49c1,
		_49c2,
		_50c5,
		_51c5,
		_52c3,
		_53c4,
		_54c1,
		_54c3,
		_55c4,
		_56c1,
		_57c1,
		_57c2,
		_58c3,
		_59c3,
		_60c1,
		_60c2,
		_61c5,
		_62c1,
		_62c2,
		_63c5,
		_64c1,
		_64c3,
		_65c5,
		_66c5,
		_67c1,
		_67c2,
		_67c3,
		_68c1,
		_68c2,
		_68c3,
		_69c4,
		_70c1,
		_70c2,
		_71c1,
		_71c3,
		_72c1,
		_72c3,
		_73c3,
		_74c1,
		_74c2,
		_74c3,
		_75c1,
		_76c1,
		_77c1,
		_78c1,
		_79c1,
		_80c1,
		_81c2,
		_82c2,
		_83c2,
		_84c3,
		_85c3,
		_86c3,
		_87c3,
		_88c3,
		_89c3,
		_90c3,
		_91c5,
		_92c5,
		_93c5,
		_94c5,
		_95c5,
		_96c5,
		_97c5)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rowsAff, _ := result.RowsAffected()
	fmt.Println("--------------------")
	lastIns, _ := result.LastInsertId()

	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "fisa_stare_tehnica.html", "there was a problem registering bridge")
		return
	}
	fmt.Fprint(w, "congrats, your bridge has been successfully created")

}
