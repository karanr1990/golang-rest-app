package handlers

import (
	"github.ibm.com/Quest-CIO/go-micro-app/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
return &Products{l}
}

func (p *Products)ServeHTTP(rw http.ResponseWriter,r *http.Request)  {
	if r.Method == http.MethodGet{
		p.getProducts(rw,r)
		return
	}

	if r.Method == http.MethodPost{
		p.addProduct(rw,r)
		return
	}

	if r.Method == http.MethodPut{
		p.l.Println("PUT",r.URL.Path)

		//expect the ID in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path,-1)

		if len(g) != 1 {
			p.l.Println("invalid URI more than one ID")
			http.Error(rw,"Invalid URI",http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2{
			p.l.Println("invalid URI more than one capture group")
			http.Error(rw,"Invalid URI",http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id,err := strconv.Atoi(idString)

		if err !=nil {
			p.l.Println("invalid URI unable to convert to number",idString)
			http.Error(rw,"Invalid URI",http.StatusBadRequest)
		}
       p.updateProducts(id,rw,r)

	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request)  {
	lp :=data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil{
		http.Error(rw,"Unable to marshal json",http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request)  {
	p.l.Println("Handle post product")

	prod := &data.Product{}

    err:=prod.FromJSON(r.Body)

    if err !=nil {
    	http.Error(rw,"unable to unmarshal json", http.StatusBadRequest)
	}
    data.AddProduct(prod)
}

func (p Products)updateProducts(id int,rw http.ResponseWriter, r *http.Request)  {
	prod := &data.Product	{}

	err:=prod.FromJSON(r.Body)

	if err !=nil {
		http.Error(rw,"unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)

    if err == data.ErrorProductNotFound{
    	http.Error(rw,"product not found",http.StatusNotFound)
		return
	}
	if err !=nil{
		http.Error(rw,"product not found",http.StatusNotFound)
		return
	}
	
}
