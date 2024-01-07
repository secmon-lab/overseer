package model

import "cloud.google.com/go/bigquery"

type BigQueryRow map[string]bigquery.Value
