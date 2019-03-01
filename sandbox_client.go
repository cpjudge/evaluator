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
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

// Submits the code to sandbox service and gets submission status
func submitCode(client pb.SandboxClient, submission *pb.Submission) *pb.CodeStatus {
	log.Printf("Submitting code for submission_id: %s", submission.GetSubmissionId())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	codeStatus, err := client.SubmitCode(ctx, submission)
	if err != nil {
		log.Fatalf("%v.SubmitCode(_) = _, %v: ", client, err)
	}
	log.Println(codeStatus)

	expectedOutputPath := "/media/pvgupta24/MyZone/Projects/go/src/github.com/cpjudge/testcases/1/output/"

	if codeStatus.CodeStatus == pb.SubmissionStatus_TO_BE_EVALUATED {
		codeStatus.CodeStatus = pb.SubmissionStatus(
				EvaluateSubmission(submission.SubmissionPath, expectedOutputPath))
	}

	log.Println(codeStatus)

	return codeStatus
}


func main() {
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
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewSandboxClient(conn)


	submitCode(client, &pb.Submission{
		Language: "cpp",
		QuestionId: "1",
		SubmissionId: "1",
		SubmissionPath: "/media/pvgupta24/MyZone/Projects/go/src/github.com/cpjudge/submissions/1/",
		TestcasesPath: "/media/pvgupta24/MyZone/Projects/go/src/github.com/cpjudge/testcases/1/input/",
		UserId: "mahim23",
	})
}
