package security

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/transformers"
)

func Locations() *schema.Table {
	return &schema.Table{
		Name:        "azure_security_locations",
		Resolver:    fetchLocations,
		Description: "https://learn.microsoft.com/en-us/rest/api/defenderforcloud/locations/get?tabs=HTTP#asclocation",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_security_locations", client.Namespacemicrosoft_security),
		Transform:   transformers.TransformWithStruct(&armsecurity.AscLocation{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{client.SubscriptionID},
	}
}

func fetchLocations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armsecurity.NewLocationsClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	pager := svc.NewListPager(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
