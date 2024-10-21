# How to write rules

Overseer provides a way to manage and maintain the security monitoring queries by using combination of SQL and Rego. You can write general SQL queries to extract, transform and correlation logs in BigQuery and Rego policies to manage and maintain the generic policies.

## SQL query

You can write general SQL queries to extract, transform and correlation logs in BigQuery. For example, the following SQL query extracts the logs from BigQuery:

```sql
/*
id: create_bucket_logs
*/
SELECT
  *
FROM
    `your_project.your_dataset.cloudaudit_googleapis_com_activity`
WHERE
    protoPayload.methodName = 'storage.buckets.create'
LIMIT 1000
```

You need to put metadata in the SQL query as follows:

- `id`: The unique identifier for the SQL query. The `id` is used to identify the SQL query in the Rego policy.

Currently Overseer supports only simple parser. The metadata should be:

- In the comment block at the top of the SQL query.
- Started with `/*` and ended with `*/` in each line.
- Written in the YAML format

## Rego policy

You can write Rego policies to manage and maintain the generic policies. For example, the following Rego policy detects the suspicious activities:

```rego
# METADATA
# title: Create bucket logs
# custom:
#   tags: [daily]
#   input: ["create_bucket_logs"]

package your_query.create_bucket_logs

import rego.v1

alert contains {
    "title": "Create suspicious Cloud Storage bucket",
    "timestamp": r.timestamp,
    "attrs": {
        "project_id": r.resource.labels.project_id,
    },
} if {
    r := input.create_bucket_logs[_]
}
```

The rules in the Rego policy are as follows:

- Policy MUST have [metadata](https://www.openpolicyagent.org/docs/latest/policy-language/#metadata) in the comment block.
  - Metadata MAY have `title` field to describe the policy.
  - Metadata `scope` MUST be `package` (default value).
  - Metadata MUST have `custom` field to describe the custom fields.
    - `tags` (array of string): The tags for the policy to control the execution.
    - `input` (array of string): The ID set of SQL queries for the policy.
- package name CAN be set arbitrary name.
- One or more `alert` rules can be defined in the package.
  - `alert` rule MUST be set of the following fields:
    - `title` (string): The title of the alert.
    - `timestamp` (string): The timestamp of the alert. The field will be parsed as `time.Time` in Overseer. It's compatible with [TIMESTAMP](https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#timestamp_type) type of BigQuery.
    - `attrs` (map): The attributes of the alert.
