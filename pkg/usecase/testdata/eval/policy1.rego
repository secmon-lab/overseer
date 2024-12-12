# METADATA
# title: Test Policy 1
# custom:
#   tags: [ "daily" ]
#   input: [ "cache1" ]

package test_policy1

import rego.v1

alert contains {
    "title": "Test Policy 1",
    "description": "Principal attempted to access data",
    "timestamp": r.latest,
    "attrs": {
        "id": r.id,
    }
} if {
    r := input.cache1[_]
    r.age > 20
}
