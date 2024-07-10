package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nguyenthenguyen/docx"
)

func generateDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только метод POST поддерживается", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	templatePath := "template.docx"
	outputPath := "output.docx"

	docTemplate, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		http.Error(w, "Ошибка при открытии шаблона", http.StatusInternalServerError)
		return
	}
	doc := docTemplate.Editable()
	doc.Replace("{{Name}}", name, -1)

	if err := doc.WriteToFile(outputPath); err != nil {
		http.Error(w, "Ошибка при сохранении документа", http.StatusInternalServerError)
		return
	}
	docTemplate.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=output.docx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	http.ServeFile(w, r, outputPath)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
        <body>
            <form action="/generate" method="post">
                <label for="name">Введите имя:</label><br>
                <input type="text" id="name" name="name" required><br>
                <input type="submit" value="Создать документ">
            </form>
        </body>
    </html>`)
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/generate", generateDoc)
	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
