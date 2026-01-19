package kube_client

import "time"

const POD_PAGINATION_LIMIT = int64(500)
const POD_PAGINATION_TIMEOUT = 30 * time.Second
