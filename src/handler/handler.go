package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"example.com/graph"
)

// menyimpan nama file secara global
var fileName = "a"

func Start() {
	// definisi rute web
	http.NewServeMux()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/polyline", polyLineHandler)
	http.HandleFunc("/map", indexHandler)
	http.HandleFunc("/Hello", helloHandler)
	http.HandleFunc("/input", inputHandler)

	// memfiksasi direktori sehingga dapat dipanggil file yang ada di dalamnya dengan /static/
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost:8080")
	// mengaktifkan server
	http.ListenAndServe(":8080", nil)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// mengoper data ke method get
	if r.Method == "POST" {
		var filepath = path.Join("views", "index.html")
		var tmpl = template.Must(template.New("result").ParseFiles(filepath))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileName = r.FormValue("fileName")
		fmt.Println(fileName)
		graf := graph.ReadFile(fileName)
		nodes := graf.GetNodes()

		if err := tmpl.Execute(w, nodes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bukan Disini alamatnyaaaa! Disini : http://localhost:8080/"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// mendapatkan input nama file di html
	if r.Method == "GET" {
		var filepath = path.Join("views", "home.html")
		var tmpl = template.Must(template.New("form").ParseFiles(filepath))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var filepath = path.Join("views", "home.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Error(w, "", http.StatusBadRequest)
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	// mendapatkan input node awal dan tujuan
	if r.Method == "GET" {
		graf := graph.ReadFile(fileName)
		nodes := graf.GetNodes()
		var filepath = path.Join("views", "input.html")
		var tmpl = template.Must(template.New("form").ParseFiles(filepath))
		var err = tmpl.Execute(w, nodes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var filepath = path.Join("views", "input.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Error(w, "", http.StatusBadRequest)
}

func polyLineHandler(w http.ResponseWriter, r *http.Request) {
	// mendapatkan input node awal dan tujuan
	if r.Method == "POST" {
		var filepath = path.Join("views", "polyline.html")
		var tmpl = template.Must(template.New("result").ParseFiles(filepath))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// graf = graph.ReadFile("tes1.txt")
		var simpulAwal = r.FormValue("simpulAwal")
		var simpulAkhir = r.FormValue("simpulAkhir")
		fmt.Println(simpulAwal)

		// mendefinisikan tipe graf
		graf := graph.ReadFile(fileName)

		// mendapatkan jarak dan rute
		distance, rute := graf.Astar(simpulAwal, simpulAkhir)

		fmt.Print("Jarak dari " + simpulAwal + " ke " + simpulAkhir + " :")
		fmt.Println(distance)
		ruteInfo := graf.GetNodeswithName(rute, distance)
		fmt.Println(ruteInfo)

		// mengoper data ke dalam server
		if err := tmpl.Execute(w, ruteInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
