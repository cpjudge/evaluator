package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/cpjudge/proto/submission"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	address            = "172.17.0.1:10000"
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

// Internal function to call the sandbox service and submit the code.
func submitCode(client pb.SandboxClient, submission *pb.Submission) *pb.CodeStatus {
	log.Printf("Submitting code for submission_id: %s", submission.GetSubmissionId())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	codeStatus, err := client.SubmitCode(ctx, submission)
	if err != nil {
		log.Fatalf("%v.SubmitCode(_) = _, %v: ", client, err)
	}
	log.Println("After sandbox execution: ", codeStatus)
	return codeStatus
}

// SubmitCode : Spawns a connection to sandbox service and calls submitCode()
func SubmitCode(submission *pb.Submission) *pb.CodeStatus {
	// Establish a connection with the sandbox service.
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewSandboxClient(conn)
	// Call internal function to submit code to sandbox service
	codeStatus := submitCode(client, submission)
	return codeStatus
}
