package main

import (
	"context"
	"flag"
	"fmt"
	"full-project/data"
	"full-project/models"
	pb "full-project/proto"
	"log"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr = flag.String("addr", "localhost:2003", "the address to connect to")
func CallSignUp(client pb.CustomerServiceClient, req *models.Customer) {
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cust := &pb.Customer{
		Id:              req.ID,
		Name:            req.Name,
		Password:        req.Password,
		Email:           req.Email,
		Address:         &pb.Address{
			Country: req.Address.Country,
			Street1: req.Address.Street1,
			Street2: req.Address.Street2,
			City:    req.Address.City,
			State:   req.Address.State,
			Zip:     req.Address.Zip,
		},
		Shippingaddress: &pb.Address{
			Country: req.ShippingAddress.Country,
			Street1: req.ShippingAddress.Street1,
			Street2: req.ShippingAddress.Street2,
			City:    req.ShippingAddress.City,
			State:   req.ShippingAddress.State,
			Zip:     req.ShippingAddress.Zip,
		},
	}
	res, err := client.Signup(ctx, cust)
	fmt.Println(cust)
	if err != nil {
		fmt.Println("err")
		log.Fatalf("could not get user: %v", err)
	}
	fmt.Println(res)
}

func CallSignIn(client pb.CustomerServiceClient, req *models.Customer) {
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := client.Signin(ctx, &pb.SignInRequest{
		Id: 	 req.ID,
		Name:   req.Name,
	})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	fmt.Println(res)
}

func main() {
	flag.Parse()

	// Set up the credentials for the connection.
	
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	rgc := pb.NewCustomerServiceClient(conn)

	CallSignUp(rgc, &models.Customer{
		ID:              "1",
		Name:            "sona",
		Password:        "sona",
		Email:           "sona@gmail.com",
		Address:         models.AddressStruct{
			Country: "india",
			Street1: "xyz",
			Street2: "xtr",
			City:    "chennai",
			State:   "tamilnadu",
			Zip:     "123456",
		},
		ShippingAddress: models.AddressStruct{
			Country: "india",
			Street1: "xyz",
			Street2: "xtr",
			City:    "chennai",
			State:   "tamilnadu",
			Zip:     "123456",
		},
	})

	CallSignIn(rgc, &models.Customer{
		ID:              "1",
		Name:            "sona",
	})
}

