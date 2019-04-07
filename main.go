package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	pb "github.com/jamesnaftel/learn-grpc/api"
	"google.golang.org/grpc"
)

func main() {
	//TODO: create subcommands to get full help (not critical for learning GRPC 😐)
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "3001", "Server port")

	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", *host, *port), grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing %s: %v\n", fmt.Sprintf("%s:%s", *host, *port), err)
		os.Exit(1)
	}
	defer conn.Close()

	c := pb.NewPodcastsClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "list":
		listPodcasts(c)
	case "query":
		queryPodcast(c, flag.Arg(1))
	case "add":
		addPodcast(c, flag.Arg(1))
	default:
		flag.Usage()
		os.Exit(0)
	}
}

func addPodcast(client pb.PodcastsClient, in string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	podcast := pb.Podcast{}
	err := json.Unmarshal([]byte(in), &podcast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal error: %v\n", err)
		return
	}

	resp, err := client.AddPodcast(ctx, &podcast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "add error: %v\n", err)
		return
	}

	printOutput([]*pb.Podcast{resp})
}

func listPodcasts(client pb.PodcastsClient) {
	empty := pb.Empty{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.GetPodcasts(ctx, &empty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	p := []*pb.Podcast{}
	for {
		podcast, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		p = append(p, podcast)
	}

	printOutput(p)
}
func queryPodcast(client pb.PodcastsClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := pb.PodcastRequest{Name: name}
	podcast, err := client.GetPodcast(ctx, &req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	printOutput([]*pb.Podcast{podcast.GetPodcast()})
}

func printOutput(podcasts []*pb.Podcast) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "Name\tAuthor\tLength\n")
	fmt.Fprintf(w, "------------------\t------------------\t----------\n")
	for _, val := range podcasts {
		fmt.Fprintf(w, "%s\t%s\t%d\n", val.GetName(), val.GetAuthor(), val.Length)
	}

	w.Flush()
}
