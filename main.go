package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

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
	ctx := context.Background()

	switch cmd := flag.Arg(0); cmd {
	case "list":
		listPodcasts(ctx, c)
	case "query":
		queryPodcast(ctx, c, flag.Arg(1))
	case "add":
		fmt.Fprintf(os.Stdout, "TODO\n")

	default:
		flag.Usage()
		os.Exit(0)
	}
}

func listPodcasts(ctx context.Context, client pb.PodcastsClient) {
	empty := pb.Empty{}
	podcasts, err := client.GetPodcasts(ctx, &empty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	printOutput(podcasts.GetPodcasts())
}
func queryPodcast(ctx context.Context, client pb.PodcastsClient, name string) {
	//Request a specific podcast - return empty
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
