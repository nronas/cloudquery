# Table: digitalocean_project_resources

This table shows data for DigitalOcean Project Resources.

https://docs.digitalocean.com/reference/api/api-reference/#tag/Project-Resources

The primary key for this table is **urn**.

## Relations

This table depends on [digitalocean_projects](digitalocean_projects).

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|urn (PK)|String|
|assigned_at|String|
|links|JSON|
|status|String|