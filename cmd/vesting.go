package main

import (
	"encoding/csv"
	"math/big"
	"os"
	"time"

	genesis "github.com/C00L-developer/cosmos-utility-tools/pkg/genesis_query"
	"github.com/urfave/cli/v2"
)

func vestingAnalyze(ctx *cli.Context) error {
	s, v, seg, err := genesis.GenesisAnalyze()
	if err != nil {
		return err
	}

	// output
	f, err := os.OpenFile("./Result/umee-vesting.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	err = w.Write([]string{"time", "supply increase", "unlocked vesting", "total supply"})
	if err != nil {
		return err
	}
	length := len(s) - 1
	period := (seg[1] - seg[0]) / int64(length)
	for i, supply := range s {
		t := time.Unix(seg[0]+int64(i)*period, 0)
		if err := w.Write([]string{t.String(), supply.String(), v[i].String(), big.NewInt(0).Add(supply, v[i]).String()}); err != nil {
			return err
		}
	}
	w.Flush() //nolint:errcheck
	f.Close() //nolint:errcheck
	return nil
}
