package cloudscheduler

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugins/source/gcp/client"

	"google.golang.org/api/cloudscheduler/v1"
)

func fetchLocations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	nextPageToken := ""
	gcpClient, err := cloudscheduler.NewService(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}

	for {
		output, err := gcpClient.Projects.Locations.List("projects/" + c.ProjectId).PageToken(nextPageToken).Context(ctx).Do()
		if err != nil {
			return err
		}
		res <- output.Locations
		if output.NextPageToken == "" {
			break
		}
		nextPageToken = output.NextPageToken
	}
	return nil
}
