package cli_test

import (
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/overseer/pkg/cli"
)

func TestListQueryFiles(t *testing.T) {
	base := "testdata/queries"
	t.Run("without recursive", func(t *testing.T) {
		paths := gt.R1(cli.ListQueryFiles(base, false)).NoError(t)
		gt.A(t, paths).Length(2).
			Have(base + "/a.sql").
			Have(base + "/b.sql")
	})

	t.Run("with recursive", func(t *testing.T) {
		paths := gt.R1(cli.ListQueryFiles(base, true)).NoError(t)
		gt.A(t, paths).Length(3).
			Have(base + "/a.sql").
			Have(base + "/b.sql").
			Have(base + "/sub/c.sql")
	})
}
