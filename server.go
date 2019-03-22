package main

import (
	"context"
	"flag"
	"log"
	"net"

	epb "github.com/cpjudge/proto/evaluator"
	spb "github.com/cpjudge/proto/submission"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	serverAddr = flag.String("server_addr", "172.17.0.1:12000", "The server address in the format of host:port")
)

type evaluatorServer struct{}

func (e *evaluatorServer) EvaluateCode(ctx context.Context, submission *spb.Submission) (*epb.CodeStatus, error) {
	codeStatus := SubmitCode(submission)
	evaluateCodeStatus := &epb.CodeStatus{}
	switch codeStatus.CodeStatus {
	case spb.SubmissionStatus_TO_BE_EVALUATED:
		evaluateStatus := EvaluateSubmission(submission.SubmissionId, submission.QuestionId)
		switch evaluateStatus {
		case 0:
			evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_CORRECT_ANSWER
		case 2:
			evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_TIME_LIMIT_EXCEEDED
		case 3:
			evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_COMPILATION_ERROR
		case 4:
			evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_RUNTIME_ERROR
		default:
			evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_WRONG_ANSWER
		}
	case spb.SubmissionStatus_TIME_LIMIT_EXCEEDED:
		evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_TIME_LIMIT_EXCEEDED
	case spb.SubmissionStatus_COMPILATION_ERROR:
		evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_COMPILATION_ERROR
	default:
		evaluateCodeStatus.CodeStatus = epb.EvaluationStatus_WRONG_ANSWER
	}
	return evaluateCodeStatus, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	epb.RegisterEvaluatorServer(grpcServer, &evaluatorServer{})
	grpcServer.Serve(lis)
}
