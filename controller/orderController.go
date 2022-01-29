package controller

import (
	"TaskOne/main/model/dynamoDb"
	"TaskOne/main/model/mongoDb"
	"TaskOne/main/pkg/aws"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"time"
)

//function to add cors on headers
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// Order api to fetch all orders and  create one
func Order(w http.ResponseWriter, r *http.Request){
	// function to enable cors in headers of request
	enableCors(&w)
	//check about the path of request
	if r.URL.Path != "/api/v1/orders" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	//check about the method of request
	switch r.Method {
	// where the method is get
	case "GET":
		// fetch all orders
		GetAll(w, r)
	case "POST":
		// create order
		Create(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

// Create godoc
// @Summary Create a new order
// @Description Create a new order with the input payload
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body model.Order true "Create order"
// @Success 200 {string} string "Order Created !"
// @Failure 400,404 {string} string "Some Error Accrue !"
// @Router /orders [post]
func Create(w http.ResponseWriter, r *http.Request) {

	// get the content of the body
	body, err := ioutil.ReadAll(r.Body)
	// declare total amount by zero
	totalAmount := 0.0
	if err != nil {
		panic(err)
	}
	//convert to json format
	var orderDate mongoDb.Order
	// store the request data in json format
	json.Unmarshal(body, &orderDate)
	// loop on item list to declare the primitive Id and total price
	for j:=0 ; j< len(orderDate.ItemList) ; j++{
		orderDate.ItemList[j].ID =primitive.NewObjectID()
		orderDate.ItemList[j].TotalPrice = float64(orderDate.ItemList[j].Qty) * orderDate.ItemList[j].UnitPrice
		totalAmount += orderDate.ItemList[j].TotalPrice
	}
	orderDate.TotalAmount = totalAmount
	// create order
	order := mongoDb.Order{
		//get the desc from the body of req
		Desc: orderDate.Desc,
		Status: "In Queue",
		Title: orderDate.Title,
		ID: primitive.NewObjectID(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ItemList: orderDate.ItemList,
		Comments: orderDate.Comments,
		ShippingLife: orderDate.ShippingLife,
		TotalAmount: orderDate.TotalAmount,

	}
	path := "orders/"+ order.ID.Hex() +"/"
	fmt.Println(path)
	aws.UploadFile("",path)
	//run function to create the order
	if dynamoDb.PutItem(order) == nil{
		fmt.Fprintf(w, "Order Created !")
	}else {
		fmt.Fprintf(w, "Some Error Accrue !")

	}

}

// Update godoc
// @Summary Update order
// @Description Update on order data
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body model.Order true "Update order"
// @Success 200 {string} string "Order Updated !"
// @Failure 400,404 {string} string "Some Error Accrue !"
//@Router /orders/id [put]
func Update(w http.ResponseWriter, r *http.Request) {

		// get the content of the body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		//convert to json format
		var orderDate mongoDb.Order
		// store the request data in json format
		json.Unmarshal(body, &orderDate)

		// create order
		order := mongoDb.Order{
			//get the desc from the body of req
			Desc: orderDate.Desc,
			Status: orderDate.Status,
			Title: orderDate.Title,
			ID:	orderDate.ID,
			CreatedAt:  orderDate.CreatedAt,
			UpdatedAt:  orderDate.UpdatedAt,
			ItemList: orderDate.ItemList,
			Comments: orderDate.Comments,
			ShippingLife: orderDate.ShippingLife,
			TotalAmount: orderDate.TotalAmount,

		}
		// run function to update the order
	//if mongoDb.UpdateOrder(order) == nil{
		if dynamoDb.UpdateOrder(order) == nil{
			fmt.Fprintf(w, "Order Updated !")
		}else {
			fmt.Println(err)
			fmt.Fprintf(w, "Some Error Accure !")

		}

}

// Delete godoc
// @Summary Delete order
// @Description Delete on order by order id
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body model.Order true "Delete order"
// @Success 200 {string} string "Order Deleted !"
// @Failure 400,404 {string} string "Some Error Accrue !"
//@Router /orders/id [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	//extract to id from url
	query := r.URL.Query()
	// initials id
	id := query.Get("id")
		// run function to create the order
	//if mongoDb.DeleteOne(id) == nil{
		if dynamoDb.DeleteItem(id) == nil{
			aws.DeleteItem("orders/"+id+"/")
			fmt.Fprintf(w, "Order Deleted !")
		}else {
			fmt.Fprintf(w, "Some Error Accure !")

		}

}

// GetAll godoc
// @Summary Get details of all orders
// @Description Get details of all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Order
// @Router /orders [get]
func GetAll(w http.ResponseWriter, r *http.Request) {

	// get object orders
	var orders []*mongoDb.Order
	// get all orders from DB
	//orders,_ = mongoDb.GetAllOrders()
	orders = dynamoDb.GetItems()
	// send all in json format
	json.NewEncoder(w).Encode(orders)
}


// GetOne godoc
// @Summary Get details of one order
// @Description Get details of one order
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Order
// @Failure 400,404 {string} string "Some Error Accrue !"
// @Router /orders/id [get]
func GetOne(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	//extract to id from url
	query := r.URL.Query()
	// initials id
	id := query.Get("id")

	if r.URL.Path != "/api/v1/orders/id" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
		var order mongoDb.Order
		// run function to create the order
	//order, err :=  mongoDb.GetOneOrder(id)
		order, err :=  dynamoDb.GetOneOrder(id)
		if err == nil{
			json.NewEncoder(w).Encode(order)
		}else {
			fmt.Fprintf(w, "Some Error Accure !")
		}
}

// UploadFile godoc
// @Summary Get details of upload file
// @Description upload file on folder of order_id
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Successfully Uploaded"
// @Router /orders/upload/id [post]
func UploadFile(w http.ResponseWriter, r *http.Request){
	//extract to id from url
	query := r.URL.Query()
	// initials id
	id := query.Get("id")
	fmt.Fprintf(w,"Uploading File\n")
	// 1. parse input
	r.ParseMultipartForm(10<<30)
	// 2. retrieve file from posted form data
	file, handler, err := r.FormFile("myFile")
	defer file.Close()

	if err != nil{
		fmt.Println("Some Error Accrue")
		fmt.Println(err)
		return
	}

	fmt.Printf("upload File : %+v\n\n", handler.Filename)
	fmt.Printf("File Size : %+v\n\n", handler.Size)
	fmt.Printf("MIME Header : %+v\n\n", handler.Header)

	// 3. write temp file on out server
	tempFile, err := ioutil.TempFile("upload","upload-*")
	defer tempFile.Close()

	if err != nil {
		fmt.Println("Some Error Accrue ")
		fmt.Println(err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil{
		fmt.Println("Some Error Accrue!")
		fmt.Println(err)
		return
	}
	tempFile.Write(fileBytes)
	path := "orders/"+id
	aws.UploadFile(tempFile.Name(), path)
	// 4. return whether or not his been success
	fmt.Fprintf(w, "Successfully Uploaded" )


}


// DownloadFile godoc
// @Summary Get details of download file
// @Description download file on download folder
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Downloaded File Done  !"
// @Router /orders/download [post]
func DownloadFile(w http.ResponseWriter, r *http.Request)  {
	enableCors(&w)

	// get the content of the body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}
	//convert to json format
	var file mongoDb.Data
	// store the request data in json format
	json.Unmarshal(body, &file)
	aws.DownloadFile(file.Path)
	fmt.Fprintf(w,"Downloaded File Done  !")
}

// GetAllFiles godoc
// @Summary Get details of all files of one order
// @Description Get details of all files of one order
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Data
// @Router /orders/files [get]
func GetAllFiles(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	//extract to id from url
	query := r.URL.Query()
	// initials id
	id := query.Get("id")

	if r.URL.Path != "/api/v1/orders/files" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, aws.ListItems(id))

}
