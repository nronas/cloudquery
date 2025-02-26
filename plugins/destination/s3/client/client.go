package client

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/cloudquery/filetypes/v2"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v2/plugins/destination"
	"github.com/rs/zerolog"
)

type Client struct {
	destination.UnimplementedUnmanagedWriter
	logger     zerolog.Logger
	spec       specs.Destination
	pluginSpec Spec

	s3Client   *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	*filetypes.Client
}

func New(ctx context.Context, logger zerolog.Logger, spec specs.Destination) (destination.Client, error) {
	if spec.WriteMode != specs.WriteModeAppend {
		return nil, fmt.Errorf("destination only supports append mode")
	}
	c := &Client{
		logger: logger.With().Str("module", "s3").Logger(),
		spec:   spec,
	}

	if err := spec.UnmarshalSpec(&c.pluginSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal s3 spec: %w", err)
	}
	if err := c.pluginSpec.Validate(); err != nil {
		return nil, err
	}
	c.pluginSpec.SetDefaults()

	filetypesClient, err := filetypes.NewClient(c.pluginSpec.FileSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to create filetypes client: %w", err)
	}
	c.Client = filetypesClient

	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	cfg.Region = c.pluginSpec.Region
	cfg.EndpointResolverWithOptions = c

	c.s3Client = s3.NewFromConfig(cfg)
	c.uploader = manager.NewUploader(c.s3Client)
	c.downloader = manager.NewDownloader(c.s3Client)

	if *c.pluginSpec.TestWrite {
		// we want to run this test because we want it to fail early if the bucket is not accessible
		timeNow := time.Now().UTC()
		if _, err := c.uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket: aws.String(c.pluginSpec.Bucket),
			Key:    aws.String(replacePathVariables(c.pluginSpec.Path, "TEST_TABLE", "TEST_UUID", c.pluginSpec.Format, timeNow)),
			Body:   bytes.NewReader([]byte("")),
		}); err != nil {
			return nil, fmt.Errorf("failed to write test file to S3: %w", err)
		}
	}

	return c, nil
}

func (*Client) Close(ctx context.Context) error {
	return nil
}

func (c *Client) ResolveEndpoint(service, region string, options ...any) (aws.Endpoint, error) {
	if c.pluginSpec.Endpoint == "" || service != s3.ServiceID {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	return aws.Endpoint{
		PartitionID:   "aws",
		URL:           c.pluginSpec.Endpoint,
		SigningRegion: region,
		Source:        aws.EndpointSourceCustom,
	}, nil
}
