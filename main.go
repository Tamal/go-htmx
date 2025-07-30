package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

type Product struct {
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Thumbnail   string  `json:"thumbnail"`
	Description string  `json:"description"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

var tmpl = template.Must(template.New("main").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8"/>
	<meta name="viewport" content="width=device-width,initial-scale=1.0"/>
	<title>Products</title>
	<script src="https://unpkg.com/htmx.org@1.9.6"></script>
	<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen p-4">
	<div class="max-w-4xl mx-auto">
		<h1 class="text-2xl font-bold mb-4">Product List</h1>
		<button 
			hx-get="/products-table" 
			hx-target="#products" 
			hx-swap="outerHTML"
			class="mb-4 px-4 py-2 bg-blue-600 text-white rounded"
		>
			Reload Products
		</button>
		<div id="products">
			{{template "table" .}}
		</div>
	</div>
</body>
</html>
`))

var tableTmpl = template.Must(template.New("table").Parse(`
{{define "table"}}
	<table class="table-auto w-full shadow rounded bg-white">
		<thead>
			<tr>
				<th class="px-4 py-2">Image</th>
				<th class="px-4 py-2">Title</th>
				<th class="px-4 py-2">Category</th>
				<th class="px-4 py-2">Price</th>
				<th class="px-4 py-2">Description</th>
			</tr>
		</thead>
		<tbody>
		{{range .}}
			<tr class="border-t">
				<td class="px-4 py-2">
					<img src="{{.Thumbnail}}" class="w-16 h-16 object-contain" alt="{{.Title}}">
				</td>
				<td class="px-4 py-2 font-semibold">{{.Title}}</td>
				<td class="px-4 py-2">{{.Category}}</td>
				<td class="px-4 py-2">${{printf "%.2f" .Price}}</td>
				<td class="px-4 py-2 text-sm">{{.Description}}</td>
			</tr>
		{{end}}
		</tbody>
	</table>
{{end}}
`))

func fetchProducts() ([]Product, error) {
	resp, err := http.Get("https://dummyjson.com/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ProductsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Products, err
}

func main() {
	tableTmpl, err := template.Must(tableTmpl, nil).Parse(``) // Avoid redefinition error
	if err != nil {
		panic(err)
	}
	tmpl = template.Must(tmpl.AddParseTree("table", tableTmpl.Lookup("table").Tree))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		products, err := fetchProducts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error fetching products: %v", err)
			return
		}
		tmpl.Execute(w, products)
	})

	http.HandleFunc("/products-table", func(w http.ResponseWriter, r *http.Request) {
		products, err := fetchProducts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error fetching products: %v", err)
			return
		}
		tableTmpl.ExecuteTemplate(w, "table", products)
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
