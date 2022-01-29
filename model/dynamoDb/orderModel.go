package dynamoDb

import (
	"TaskOne/main/model/mongoDb"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
)

// ItemList item - struct to map with mongodb documents
type ItemList struct {
	ID       		  string `bson:"ItemId"`
	ItemName 			  string `bson:"name"`
	Qty 			  int  `bson:"qty"`
	UnitPrice 		  float64 `bson:"unit-price"`
	TotalPrice 		  float64 `bson:"total-price"`
	OrderId       	  string `bson:"OrderId"`


}
type ShippingLife struct {
	ID       		   	    string `bson:"ID"`
	TrackingNumber 			string `bson:"tracking_number"`
	ShippingMethod 			string `bson:"shipping-method"`
	Order       		    string `bson:"OrderId"`

}
type Comments struct {
	ID       		      string `bson:"ID"`
	OrderId       		  string `bson:"OrderId"`
	Body string 		  `bson:"body"`
}
//Order - struct to map with mongodb documents
type Order struct {
	ID                string `bson:"_id"`
	CreatedAt         time.Time          `bson:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at"`
	Description     	 	  string             `bson:"desc"`
	StatusNw          string             `bson:"status"`
	Title             string             `bson:"title"`
	TotalAmount       float64            `bson:"total_amount"`
	ShippingLifeID 	  string 			 `bson:"shipping_life_id"`

}
func connection() *session.Session {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess, err := session.NewSession(&aws.Config{
		// add region
		Region:     aws.String(os.Getenv("Reg")),
		// add access key and secret access key to the to the credential
		Credentials: credentials.NewStaticCredentials(os.Getenv("Access_key_ID"),os.Getenv("Secret_access_key"),""),
	})

	if err != nil{
		exitError("Some Error Accrue when open aws session ", err)
	}
	return sess
}

// ListTables function to get list of all tables
func ListTables(){

	sess  := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

}

// CreateOrderTable function to create a table
func CreateOrderTable (){

	// create a session
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// Create table Movies
	tableName := "ShippingLife"
	// add table content

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},

			{
				AttributeName: aws.String("ShippingLifeID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},

			{
				AttributeName: aws.String("ShippingLifeID"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
}

// CreateTable function to create a table
func CreateTable (tableName string){

	// create a session
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// add table content

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
}

//DeleteTable - function to delete table
func DeleteTable(tableName string) {

	sess := connection()

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	resp, err := svc.DeleteTable(params)

	if err!=nil {
		exitError("can't delete a table ", err)

	}


	log.Printf("Table %s deleted:\n%s\n", tableName, resp.String())
}

//AddOrder - function to add order in order table
func AddOrder(order Order)error{
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {

		return err
	}
	// Create item in table order
	tableName := "Order"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}
	fmt.Println("Successfully added  to table " + tableName)
	return nil
}

//AddShippingInfo - function to add shipping Life in shipping Life table
func AddShippingInfo(order ShippingLife)error{
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatalf("Got error marshalling new ShippingLife item: %s", err)
		return err

	}
	// Create item in table order
	tableName := "ShippingLife"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err

	}
	fmt.Println("Successfully added  to table " + tableName)
	return nil
}

//AddItemList - function to add Add Item List in Item List table
func AddItemList(order ItemList)error{
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
		return err

	}
	// Create item in table order
	tableName := "ItemLists"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err

	}
	fmt.Println("Successfully added  to table " + tableName)
	return nil
}

//AddComment - function to add Add Comment in Comments table
func AddComment(order Comments) error{
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
		return err

	}
	// Create item in table order
	tableName := "Comments"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err

	}
	fmt.Println("Successfully added  to table " + tableName)
	return nil
}

//PutItem - function to put item
func PutItem(order mongoDb.Order)error{
	orderId := primitive.NewObjectID().Hex()
	shippinglifeId := primitive.NewObjectID().Hex()

	// Add order data
	orderItem := Order{
		ID:             orderId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Description:           order.Desc,
		StatusNw:         order.Status,
		Title:          order.Title,
		TotalAmount:    order.TotalAmount,
		ShippingLifeID: shippinglifeId,
	}
	// Call function to add order in order table
	err := AddOrder(orderItem)
	if err != nil {
		return err
	}
	// Add shipping info data
	shippingLife := ShippingLife{
		ID:             shippinglifeId,
		ShippingMethod: order.ShippingLife.ShippingMethod,
		TrackingNumber: order.ShippingLife.TrackingNumber,
		Order:        orderId,
	}
	// Call function to add shippingLife in shippingLife table
	err =  AddShippingInfo(shippingLife)
	if err != nil {
		return err
	}
	//loop on items list to store each item
	for i:=0;i<len(order.ItemList);i++{
		// add one itemList data
		orderItemList := ItemList{
			ID:         primitive.NewObjectID().Hex(),
			ItemName:       order.ItemList[i].Name,
			Qty:        order.ItemList[i].Qty,
			UnitPrice:  order.ItemList[i].UnitPrice,
			TotalPrice: order.ItemList[i].TotalPrice,
			OrderId:    orderId,
		}
		// Call function to add orderItemList in ItemList table
		err = AddItemList(orderItemList)
		if err != nil {
			return err
		}
	}

	// loop on comments to store each comment
	for i:=0;i<len(order.Comments);i++{
		comment := Comments{
			ID:      primitive.NewObjectID().Hex(),
			OrderId: orderId,
			Body:    order.Comments[i].Body,
		}
		// Call function to add comment in comment table
		err = AddComment(comment)
		if err != nil {
			return err
		}
	}
	return nil

}

//GetTables - function to execute to get table content
func GetTables(tableName string) (*dynamodb.ScanOutput, error) {
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName:  aws.String(tableName),
	}
	// Make the DynamoDB Query API call
	return svc.Scan(params)
}

//GetOrder - function to get one item of shipping with the order id
func GetOrder(id string) []*Order{
	tableName:="Order"
	var allOrders []*Order
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// create a filter with id
	filt := expression.Name("ID").Equal(expression.Value(id))
	// Get back the id , CreatedAt,UpdatedAt,Desc,Status,Title and TotalAmount
	proj := expression.NamesList(expression.Name("ID"),expression.Name("CreatedAt"),
		expression.Name("UpdatedAt"),expression.Name("Description"),expression.Name("StatusNw"),
		expression.Name("Title"),expression.Name("TotalAmount"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	for _, i := range result.Items {
		item := Order{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		allOrders = append(allOrders, &item)
	}
	return allOrders
}
//GetShipping - function to get one item of shipping with the order id
func GetShipping(id string) []*ShippingLife{
	tableName:="ShippingLife"
	var allShipping []*ShippingLife
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// create a filter with id
	filt := expression.Name("Order").Equal(expression.Value(id))
	// Get back the id
	proj := expression.NamesList(expression.Name("ID"),expression.Name("ShippingMethod"),
		expression.Name("TrackingNumber"),expression.Name("Order"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	for _, i := range result.Items {
		item := ShippingLife{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		allShipping = append(allShipping, &item)
	}
	return allShipping
}

//GetItemsList - function to get all items with the order id
func GetItemsList(id string)[]*ItemList{
	tableName:="ItemLists"
	var allItems []*ItemList
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// create a filter with id
	filt := expression.Name("OrderId").Equal(expression.Value(id))
	// Get back the id
	proj := expression.NamesList(expression.Name("OrderId"),expression.Name("ItemName"),expression.Name("Qty"),
		expression.Name("UnitPrice"), expression.Name("TotalPrice"), expression.Name("ID"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	for _, i := range result.Items {
		item := ItemList{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		allItems = append(allItems, &item)
	}
	return allItems
}

//GetComments - function to get all comments with the order id
func GetComments(id string)[]*Comments{
	tableName:="Comments"
	var allComments []*Comments
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// create a filter with id
	filt := expression.Name("OrderId").Equal(expression.Value(id))
	// Get back the id
	proj := expression.NamesList(expression.Name("OrderId"),expression.Name("Body"),expression.Name("ID"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	for _, i := range result.Items {
		item := Comments{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		allComments = append(allComments, &item)
	}
	return allComments
}

//CompressOrder - function to add all structs in one struct
func CompressOrder(order []*Order) []*mongoDb.Order {
	var allData []*mongoDb.Order

	for i:=0;i<len(order);i++{
		id, err :=  primitive.ObjectIDFromHex(order[i].ID)
		if err!=nil{
			exitError("Error in Convert Id ", err)
		}
		shippingLife := GetShipping(order[i].ID)
		itemList := GetItemsList(order[i].ID)
		comments := GetComments(order[i].ID)
		item := mongoDb.Order{
			ID:             id ,
			CreatedAt:      order[i].CreatedAt,
			UpdatedAt:      order[i].UpdatedAt,
			Desc:           order[i].Description,
			Status:        order[i].StatusNw,
			Title:         order[i].Title,
			TotalAmount:    order[i].TotalAmount,
		}
		if len(shippingLife) > 0 {
			item.ShippingLife.TrackingNumber = shippingLife[0].TrackingNumber
			item.ShippingLife.ShippingMethod = shippingLife[0].ShippingMethod
		}
		for j:=0 ;j<len(comments);j++ {
			item.Comments = append(item.Comments,mongoDb.Comments{Body: comments[j].Body})
		}

		for k:=0 ;k<len(itemList);k++ {
			id , _ = primitive.ObjectIDFromHex(itemList[k].ID)
			item.ItemList = append(item.ItemList, mongoDb.ItemList{
				ID:         id,
				Name:       itemList[k].ItemName,
				Qty:        itemList[k].Qty,
				UnitPrice:  itemList[k].UnitPrice,
				TotalPrice: itemList[k].TotalPrice,
			})

		}

		allData = append(allData, &item)
	}
	return allData


}

//GetItems - function to get all data and return it
func GetItems()[]*mongoDb.Order {
	// create array from each struct to store all tables data
	var allOrders []*Order
	orderLen :=0

	//  get order data
	result, err := GetTables("Order")
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	for _, i := range result.Items {
		orderLen++
		//fmt.Println("Load New item !")
		item := Order{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		allOrders = append(allOrders, &item)
	}


	return CompressOrder(allOrders)

}

//DeleteItem - function to delete item
func DeleteItem(id string)error{
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := "Order"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		return err
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}
	fmt.Println("Deleted  from table " + tableName)
	// delete from comments table
	tableName = "Comments"
	allComments := GetComments(id)
	for i:=0;i<len(allComments);i++ {
		input = &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(allComments[i].ID),
				},
			},
			TableName: aws.String(tableName),
		}
	}

		_, err = svc.DeleteItem(input)
		if err != nil {
			log.Fatalf("Got error calling DeleteItem: %s", err)
			return err
		}
		fmt.Println("Deleted  from table " + tableName)

	// delete from ItemList table
	tableName = "ItemLists"
	allItems := GetItemsList(id)
	for i:=0;i<len(allItems);i++ {
		input := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(allItems[i].ID),
				},
			},
			TableName: aws.String(tableName),
		}

		_, err := svc.DeleteItem(input)
		if err != nil {
			log.Fatalf("Got error calling DeleteItem: %s", err)
			return err
		}
		fmt.Println("Deleted  from table " + tableName)
	}
	//delete from ItemList table
	tableName = "ShippingLife"
	allShipping := GetShipping(id)

	for i:=0;i<len(allShipping);i++ {
		fmt.Println(allShipping[i].ID)
		input := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(allShipping[i].ID),
				},
			},
			TableName: aws.String(tableName),
		}

		_, err := svc.DeleteItem(input)
		if err != nil {
			log.Fatalf("Got error calling DeleteItem: %s", err)
			return err
		}
		fmt.Println("Deleted  from table " + tableName)
	}


	return nil
}

//GetOneOrder - function to return one order
func GetOneOrder(id string)(mongoDb.Order, error){
	Orders   :=  GetOrder(id)
	shippingLife :=  GetShipping(Orders[0].ID)
	comments :=  GetComments(Orders[0].ID)
	itemList :=  GetItemsList(Orders[0].ID)

	objectId, err :=  primitive.ObjectIDFromHex(id)
	if err!=nil{
		exitError("Error in Convert Id ", err)
	}
	order := mongoDb.Order{
		ID:           objectId,
		CreatedAt:    Orders[0].CreatedAt,
		UpdatedAt:    Orders[0].UpdatedAt,
		Desc:         Orders[0].Description,
		Status:       Orders[0].StatusNw,
		Title:        Orders[0].Title,
		TotalAmount:  Orders[0].TotalAmount,
	}
	if len(shippingLife) > 0 {
		order.ShippingLife.TrackingNumber = shippingLife[0].TrackingNumber
		order.ShippingLife.ShippingMethod = shippingLife[0].ShippingMethod
	}
	for j:=0 ;j<len(comments);j++ {
		order.Comments = append(order.Comments,mongoDb.Comments{Body: comments[j].Body})
	}

	for k:=0 ;k<len(itemList);k++ {
		objectId , _ = primitive.ObjectIDFromHex(itemList[k].ID)
		order.ItemList = append(order.ItemList, mongoDb.ItemList{
			ID:         objectId,
			Name:       itemList[k].ItemName,
			Qty:        itemList[k].Qty,
			UnitPrice:  itemList[k].UnitPrice,
			TotalPrice: itemList[k].TotalPrice,
		})

	}

	return order, nil
}

//UpdateOrder - function to update all order data
func UpdateOrder(order mongoDb.Order)error{
	//fmt.Println(order)
	sess := connection()
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	TotalAmount := fmt.Sprintf("%v", order.TotalAmount)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {
				S: aws.String(order.Desc),
			},
			":s": {
				S: aws.String(order.Status),
			},
			":t": {
				S: aws.String(order.Title),
			},
			":m": {
				N: aws.String(TotalAmount),
			},
			":u": {
				S: aws.String(time.Now().Format("2006-01-02T15:04:05.999+05:05")),
			},
		},
		TableName: aws.String("Order"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(order.ID.Hex()),
			},

		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression:  aws.String("set Description = :d , StatusNw = :s, Title = :t, TotalAmount= :m, UpdatedAt =:u"),
	}
	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem on order table: %s", err)
	}

	// get the shipping id
	shipping := GetShipping(order.ID.Hex())
	input = &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{

			":s": {
				S: aws.String(order.ShippingLife.ShippingMethod),
			},
			":t": {
				S: aws.String(order.ShippingLife.TrackingNumber),
			},

		},
		TableName: aws.String("ShippingLife"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(shipping[0].ID),
			},

		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression:  aws.String("set ShippingMethod = :s , TrackingNumber = :t"),
	}
	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem on shipping table: %s", err)
	}

	// get the items lists ids
	itemLists := GetItemsList(order.ID.Hex())

	for i:=0;i<len(itemLists);i++ {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{

				":n": {
					S: aws.String(order.ItemList[i].Name),
				},
				":q": {
					N: aws.String(fmt.Sprintf("%v", order.ItemList[i].Qty)),
				},
				":t": {
					N: aws.String(fmt.Sprintf("%v", order.ItemList[i].TotalPrice)),
				},
				":u": {
					N: aws.String(fmt.Sprintf("%v", order.ItemList[i].UnitPrice)),
				},

			},
			TableName: aws.String("ItemLists"),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(itemLists[i].ID),
				},

			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression:  aws.String("set ItemName = :n , Qty = :q, TotalPrice =:t, UnitPrice =:u"),
		}
		_, err = svc.UpdateItem(input)
		if err != nil {
			log.Fatalf("Got error calling UpdateItem on shipping table: %s", err)
		}

	}
	fmt.Println("Successfully updated Item Lists Table")
	// get the items lists ids
	comments := GetComments(order.ID.Hex())

	fmt.Println("lOAD OUT OF LOOP ")
	for i:=0;i<len(comments);i++ {
		fmt.Println("lOAD IN OF LOOP ")

		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{

				":d": {
					S: aws.String(order.Comments[i].Body),
				},

			},
			TableName: aws.String("Comments"),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(comments[i].ID),
				},

			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression:  aws.String("set Body = :d"),
		}
		_, err = svc.UpdateItem(input)
		if err != nil {
			log.Fatalf("Got error calling UpdateItem on shipping table: %s", err)
		}

	}
	return nil
}

func exitError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
