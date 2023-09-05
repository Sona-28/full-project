package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"full-project/data"
	"full-project/models"
	pb "full-project/proto"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var port = flag.Int("port", 2003, "the port to serve on")
var secretKey = []byte("your-secret-key")

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterCustomerServiceServer(s, &customerServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type customerServer struct {
	pb.UnimplementedCustomerServiceServer
}

type CustomClaims struct {
	ID   string    `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func (s *customerServer) SignUp(ctx context.Context, req *pb.Customer) (*pb.CustomerResponse, error) {
	fmt.Println(req)
	cust := &models.Customer{
		ID:              req.Id,
		Name:            req.Name,
		Password:        req.Password,
		Email:           req.Email,
		Address:         models.AddressStruct{
			Country: req.Address.Country,
			Street1: req.Address.Street1,
			Street2: req.Address.Street2,
			City:    req.Address.City,
			State:   req.Address.State,
			Zip:     req.Address.Zip,
		},
		ShippingAddress: models.AddressStruct{
			Country: req.Shippingaddress.Country,
			Street1: req.Shippingaddress.Street1,
			Street2: req.Shippingaddress.Street2,
			City:    req.Shippingaddress.City,
			State:   req.Shippingaddress.State,
			Zip:     req.Shippingaddress.Zip,
		},
	}
	fmt.Println(cust)
	return &pb.CustomerResponse{
		Message: "created customer",
	},nil
}

func (s *customerServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.Token, error) {
	token, err := generateJWT(req.Id, req.Name)
	if err != nil {
		log.Fatal("Error:", err)
	}
	return &pb.Token{
		Token: token,
	},nil
}

func generateJWT(id string, name string) (string, error) {
	// Create custom claims
	claims := CustomClaims{
		id,
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(), 
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// func verifyAndDecodeJWT(tokenString string) (*CustomClaims, error) {
// 	// Parse the token with custom claims
// 	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if the token is valid
// 	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, fmt.Errorf("invalid token")
// }
