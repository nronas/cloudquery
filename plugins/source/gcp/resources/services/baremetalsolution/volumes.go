package baremetalsolution

import (
	"context"

	"google.golang.org/api/iterator"

	pb "cloud.google.com/go/baremetalsolution/apiv2/baremetalsolutionpb"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/transformers"
	"github.com/cloudquery/plugins/source/gcp/client"

	baremetalsolution "cloud.google.com/go/baremetalsolution/apiv2"
)

func Volumes() *schema.Table {
	return &schema.Table{
		Name:        "gcp_baremetalsolution_volumes",
		Description: `https://cloud.google.com/bare-metal/docs/reference/rest/v2/projects.locations.volumes#Volume`,
		Resolver:    fetchVolumes,
		Multiplex:   client.ProjectMultiplexEnabledServices("baremetalsolution.googleapis.com"),
		Transform:   client.TransformWithStruct(&pb.Volume{}, transformers.WithPrimaryKeys("Name")),
		Columns: []schema.Column{
			{
				Name:     "project_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveProject,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
		Relations: []*schema.Table{
			VolumeLuns(),
		},
	}
}

func fetchVolumes(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	req := &pb.ListVolumesRequest{
		Parent: "projects/" + c.ProjectId + "/locations/-",
	}
	gcpClient, err := baremetalsolution.NewClient(ctx, c.ClientOptions...)
	if err != nil {
		return err
	}
	it := gcpClient.ListVolumes(ctx, req, c.CallOptions...)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		res <- resp
	}
	return nil
}
