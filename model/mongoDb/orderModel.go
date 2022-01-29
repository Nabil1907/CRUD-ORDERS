package mongoDb

import (
	"TaskOne/main/pkg/connection"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Data struct {
	Path string
	Size int
}

// ItemList item - struct to map with mongodb documents
type ItemList struct {
	ID        primitive.ObjectID `bson:"ItemId"`
	Name 		string `bson:"name"`
	Qty 		int  `bson:"qty"`
	UnitPrice   float64 `bson:"unit-price"`
	TotalPrice  float64 `bson:"total-price"`

}
type ShippingLife struct {
	TrackingNumber string `bson:"tracking_number"`
	ShippingMethod string `bson:"shipping-method"`
}
type Comments struct {
	Body string `bson:"body"`
}
//Order - struct to map with mongodb documents
type Order struct {
	ID           primitive.ObjectID `bson:"_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	Desc         string             `bson:"desc"`
	Status       string             `bson:"status"`
	Title        string             `bson:"title"`
	ItemList     []ItemList         `bson:"item_list"`
	ShippingLife ShippingLife       `bson:"shipping_life"`
	Comments     []Comments         `bson:"comments"`
	TotalAmount  float64            `bson:"total_amount"`
}

// CreateOrder CreateOrder - Insert a new document in the collection.
func CreateOrder(task Order) error {

	//Get MongoDB connection using connection.
	client, err := connection.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connection.DB).Collection(connection.ORDERS)
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(context.TODO(), task)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//GetAllOrders - Get All Orders for collection
func GetAllOrders() ([]Order, error) {
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	var orders []Order
	//Get MongoDB connection using Connection.
	client, err := connection.GetMongoClient()
	if err != nil {
		return orders, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connection.DB).Collection(connection.ORDERS)
	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return orders, findError
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := Order{}
		err := cur.Decode(&t)
		if err != nil {
			return orders, err
		}
		orders = append(orders, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(orders) == 0 {
		return orders, mongo.ErrNoDocuments
	}
	return orders, nil
}

// GetOneOrder - get one order by id
func GetOneOrder(id string) (Order, error) {
	// get the id
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "_id", Value: idPrimitive}}
	var order Order
	//Get MongoDB connection using Connection.
	client, err := connection.GetMongoClient()
	if err != nil {
		return order, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connection.DB).Collection(connection.ORDERS)
	//Perform Find operation & validate against the error.
	findError := collection.FindOne(context.TODO(), filter).Decode(&order)
	if findError != nil {
		return order, findError
	}
	return order, nil
}

//DeleteOne -delete one order from the  collection
func DeleteOne(id string) error {
	// get the id
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "_id", Value: idPrimitive}}
	//Get MongoDB connection using model.
	client, err := connection.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connection.DB).Collection(connection.ORDERS)
	//Perform DeleteOne operation & validate against the error.
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

// UpdateOrder update Date - update the data
func UpdateOrder(order Order) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "_id", Value: order.ID}}
	//Get MongoDB connection using connection.
	client, err := connection.GetMongoClient()
	if err != nil {
		fmt.Println(err)
		return err
	}
	//get the collection from DB
	collection := client.Database(connection.DB).Collection(connection.ORDERS)
	//Perform UpdateOne operation & validate against the error.
	_, err = collection.UpdateOne(context.TODO(), filter,
		bson.D{
		// update all struct
			{"$set", order},
		},
	)
	//check if any error accrue !
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}


