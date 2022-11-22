package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	chainquery "github.com/C00L-developer/cosmos-utility-tools/pkg/chain_query"
	"github.com/urfave/cli/v2"
)

func monitorValidator(ctx *cli.Context) error {
	chainName := ctx.String(flagChain)
	accountPrefix := ctx.String(flagPrefix)
	c := chainquery.NewClient(chainName, accountPrefix)
	vals, err := c.GetValidators()
	totalDels := make([]chainquery.Delegator, 0)
	if err != nil {
		return err
	}
	for i, val := range vals {
		dels, err := c.GetDelegators(val.ValAddr)
		if err != nil {
			return err
		}
		for _, del := range dels {
			valAddr, err := c.GetAccAddress(del.ValAddr)
			if err != nil {
				return err
			}
			if del.DelAddr == valAddr {
				val.SelfDelegation = del.Token
				vals[i] = val
			}
		}
		totalDels = append(totalDels, dels...)
	}

	// output 1
	f, err := os.OpenFile(fmt.Sprintf("./Result/%s-validator.csv", chainName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].VotingPower > vals[j].VotingPower
	})
	err = w.Write([]string{"moniker", "validator address", "self delegation", "total delegation", "voting power"})
	if err != nil {
		return err
	}
	for _, val := range vals {
		if err := w.Write([]string{val.Moniker, val.ValAddr, val.SelfDelegation.String(), val.TotalDelegation.String(), fmt.Sprintf("%d", val.VotingPower)}); err != nil {
			return err
		}
	}
	w.Flush() //nolint:errcheck
	f.Close() //nolint:errcheck

	// output 2
	f, err = os.OpenFile(fmt.Sprintf("./Result/%s-delegator.csv", chainName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	w = csv.NewWriter(f)
	sort.Slice(totalDels, func(i, j int) bool {
		return totalDels[i].VotingPower > totalDels[j].VotingPower
	})
	err = w.Write([]string{"delegator", "validator", "voting power"})
	if err != nil {
		return err
	}
	for _, del := range totalDels {
		if err := w.Write([]string{del.DelAddr, del.ValAddr, fmt.Sprintf("%.06f", del.VotingPower)}); err != nil {
			return err
		}
	}
	w.Flush() //nolint:errcheck
	f.Close() //nolint:errcheck

	// output 2
	f, err = os.OpenFile(fmt.Sprintf("./Result/%s-multival-delegator.csv", chainName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	w = csv.NewWriter(f)
	multiDels := make(map[string][]chainquery.Delegator)
	for _, del := range totalDels {
		if dels, found := multiDels[del.DelAddr]; found {
			multiDels[del.DelAddr] = append(dels, del)
		} else {
			multiDels[del.DelAddr] = []chainquery.Delegator{del}
		}
	}

	err = w.Write([]string{"delegator", "validator", "bonded token"})
	if err != nil {
		return err
	}
	for delAddr, dels := range multiDels {
		if len(dels) > 1 {
			for _, del := range dels {
				if err := w.Write([]string{delAddr, del.ValAddr, del.Token.String()}); err != nil {
					return err
				}
			}
		}
	}
	w.Flush() //nolint:errcheck
	f.Close() //nolint:errcheck

	return nil
}
