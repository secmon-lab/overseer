# METADATA
# title: Test Policy
# custom:
#  tags: [ "red" ]
#  input: [ "stone" ]

package mytest.blue

import rego.v1

alert contains {
    "title": "Blue Alert",
    "attrs": {
        "color": "blue"
    },
} if {
    r := input.stone == "blue"
}
