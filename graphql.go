package main

import (
	"context"
	"github.com/machinebox/graphql"
	"log"
)


type GraphqlClient struct{
	Client *graphql.Client
}

func (g *GraphqlClient) Init(url string) {
	// create a client (safe to share across requests)
	g.Client = graphql.NewClient(url)
}

type ProposalRes struct{
	Proposal struct{
		ID string
		Title string
		Body string
		Choices []string
		Start int64
		End int64
		Snapshot string
		State string
		Author string
		Link string
		Type string
		Space struct{
			Id string
			Name string
		}
	}
}

func (g *GraphqlClient) GetSingleProposal(id string) ProposalRes {

	// make a request
	req := graphql.NewRequest(`
    query ($key: String!) {
        proposal (id:$key) {
			id
			title
			body
			choices
			start
			end
			snapshot
			state
			author
			created
			plugins
			network
			link
			type
			space {
				id
			  	name
			}
        }
    }
`)

	// set any variables
	req.Var("key", id)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var p ProposalRes
	if err := g.Client.Run(ctx, req, &p); err != nil {
		log.Println(err)
	}
	return p
}