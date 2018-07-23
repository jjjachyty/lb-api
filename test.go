package main

import (
	"fmt"
	"lb-api/models"
	"lb-api/models/order"

	"labix.org/v2/mgo/bson"
)

func main() {
	var payment = new(order.Payment)

	err := models.FindOne("payment", bson.M{"_id": bson.ObjectIdHex("5b53fbb2421aa9f349000001")}, payment)
	fmt.Println("payment", payment, err)
}
