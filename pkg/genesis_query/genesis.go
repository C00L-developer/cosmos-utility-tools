package genesis

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vesttypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	tmjson "github.com/tendermint/tendermint/libs/json"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppState struct {
	AuthState authtypes.GenesisState `json:"auth"`
	MintState minttypes.GenesisState `json:"mint"`
	BankState banktypes.GenesisState `json:"bank"`
}

func GenesisAnalyze() ([]*big.Int, []*big.Int, []int64, error) {
	var genesisDoc *tmtypes.GenesisDoc
	bz, err := tmos.ReadFile("./Result/genesis.json")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read genesis file: %s", err)
	}

	if err = tmjson.Unmarshal(bz, &genesisDoc); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal genesis doc: %s", err)
	}

	appState := new(AppState)
	if err = json.Unmarshal(genesisDoc.AppState, appState); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal appstate: %v", err)
	}

	cvAccs := make([]*vesttypes.ContinuousVestingAccount, 0)
	dvAccs := make([]*vesttypes.DelayedVestingAccount, 0)
	st := genesisDoc.GenesisTime.Unix()
	en := st + 5*int64(appState.MintState.Params.BlocksPerYear)*15 // 15 years after
	for _, account := range appState.AuthState.Accounts {
		bz, err := account.MarshalJSON()
		if err != nil {
			return nil, nil, nil, err
		}
		var accMsg map[string]json.RawMessage
		if err = json.Unmarshal(bz, &accMsg); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal account: %v", err)
		}
		if strings.Contains(string(accMsg["@type"]), "ContinuousVestingAccount") {
			cvAcc := new(vesttypes.ContinuousVestingAccount)
			if err := json.Unmarshal(bz, cvAcc); err != nil {
				return nil, nil, nil, err
			}
			if cvAcc.EndTime > en {
				en = cvAcc.EndTime
			}
			cvAccs = append(cvAccs, cvAcc)
		} else if strings.Contains(string(accMsg["@type"]), "DelayedVestingAccount") {
			dvAcc := new(vesttypes.DelayedVestingAccount)
			if err := json.Unmarshal(bz, dvAcc); err != nil {
				return nil, nil, nil, err
			}
			if dvAcc.EndTime > en {
				en = dvAcc.EndTime
			}
			dvAccs = append(dvAccs, dvAcc)
		}
	}
	inflation := appState.MintState.Params.InflationMax.Add(appState.MintState.Params.InflationMin).BigInt()
	inflation.Quo(inflation, big.NewInt(2))
	s, v := analyzeSupply(st, en, appState.MintState.Params.BlocksPerYear, inflation, appState.BankState.Supply[0].Amount.BigInt(), cvAccs, dvAccs)
	return s, v, []int64{st, en}, nil
}

func getIndex(stTime, cur, period int64) int {
	return int((cur - stTime + period - 1) / period)
}

func analyzeSupply(stTime, enTime int64, blocksPerYear uint64, inflation, intialSupply *big.Int, cvs []*vesttypes.ContinuousVestingAccount, dvs []*vesttypes.DelayedVestingAccount) ([]*big.Int, []*big.Int) {
	// unlocked vesting account supply
	// Y = a * X + b, here X is a time axis
	fixedLength := 1000000
	period := (enTime - stTime) / int64(fixedLength)
	a := make([]*big.Float, fixedLength+1)
	b := make([]*big.Float, fixedLength+1)

	for i := 0; i <= fixedLength; i++ {
		a[i] = big.NewFloat(0)
		b[i] = big.NewFloat(0)
	}
	sum := big.NewFloat(0)
	for _, dv := range dvs {
		amount := big.NewFloat(0).SetInt(dv.BaseVestingAccount.OriginalVesting[0].Amount.BigInt())
		sum.Add(sum, amount)
		index := getIndex(stTime, dv.EndTime, period)
		b[index].Add(b[index], amount)
	}

	for _, cv := range cvs {
		amount := big.NewFloat(0).SetInt(cv.BaseVestingAccount.OriginalVesting[0].Amount.BigInt())
		sum.Add(sum, amount)
		st := getIndex(stTime, cv.GetStartTime(), period)
		en := getIndex(stTime, cv.EndTime-period+1, period)
		seg := cv.EndTime - cv.GetStartTime()
		// update the coefficient
		df := big.NewFloat(0).Quo(amount, big.NewFloat(0).SetInt64(seg))
		a[st].Add(a[st], df)
		a[en].Sub(a[en], df)
		// udate the constant
		df.Mul(df, big.NewFloat(0).SetInt64(cv.GetStartTime()))
		b[st].Sub(b[st], df)
		b[en].Add(b[en], df)
		b[en+1].Add(b[en+1], amount)
	}

	supply := make([]*big.Int, fixedLength+1)
	unlockedVesting := make([]*big.Int, fixedLength+1)
	supply[0] = intialSupply
	inflationRatio, _ := big.NewFloat(0).SetInt(inflation).Float64()
	inflationRatio /= 1e18
	incRatio := big.NewFloat(math.Pow(1.0+inflationRatio/float64(blocksPerYear), float64(period)/5.0))

	for i := 0; i <= fixedLength; i++ {
		x := int64(i)*period + stTime
		if i > 0 {
			supply[i], _ = big.NewFloat(0).Mul(big.NewFloat(0).SetInt(supply[i-1]), incRatio).Int(nil)
			b[i].Add(b[i], b[i-1])
			a[i].Add(a[i], a[i-1])
		}
		unlockedVesting[i], _ = big.NewFloat(0).Add(b[i], big.NewFloat(0).Mul(big.NewFloat(0).SetInt64(x), a[i])).Int(nil)
	}

	return supply, unlockedVesting
}
