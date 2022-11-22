package chainquery

import (
	"context"
	"log"
	"math/big"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
	query "github.com/strangelove-ventures/lens/client/query"
	"go.uber.org/zap"
)

type Client struct {
	chainName string
	client    *lens.ChainClient
}

type Validator struct {
	Moniker         string
	ValAddr         string
	VotingPower     int64
	SelfDelegation  *big.Int
	TotalDelegation *big.Int
}

type Delegator struct {
	ValAddr     string
	DelAddr     string
	Token       *big.Int
	VotingPower float64
}

func getChainClient(chainName, accountPrefix string) *lens.ChainClient {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	// Fetches chain info from chain registry
	chainInfo, err := registry.DefaultChainRegistry(logger).GetChain(ctx, chainName)
	if err != nil {
		log.Fatalf("failed to get chain info. err: %v", err)
	}

	// Use Chain info to select random endpoint
	rpc, err := chainInfo.GetRandomRPCEndpoint(ctx)
	if err != nil {
		log.Fatalf("failed to get RPC endpoints on chain %s. err: %v", chainInfo.ChainName, err)
	}

	// Creates client object to pull chain info
	chainClient, err := lens.NewChainClient(logger, &lens.ChainClientConfig{
		RPCAddr:        rpc,
		KeyringBackend: "test",
		Timeout:        "300s",
		AccountPrefix:  accountPrefix,
	}, os.Getenv("HOME"), os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("failed to build new chain client for %s. err: %v", chainInfo.ChainID, err)
	}

	return chainClient
}

func NewClient(chainName string, accountPrefix string) *Client {
	return &Client{
		chainName: chainName,
		client:    getChainClient(chainName, accountPrefix),
	}
}

func (c *Client) getNewQuerier() *query.Query {
	return &query.Query{
		Client:  c.client,
		Options: query.DefaultOptions(),
	}
}

func (c *Client) GetValidators() ([]Validator, error) {
	vq := c.getNewQuerier()
	validators := make([]Validator, 0)

	for {
		res, err := vq.Staking_Validators("BOND_STATUS_BONDED") // get only bonded validators
		if err != nil {
			return nil, err
		}

		for _, val := range res.Validators {
			validators = append(validators, Validator{
				Moniker:         val.Description.Moniker,
				ValAddr:         val.OperatorAddress,
				VotingPower:     sdk.TokensToConsensusPower(val.Tokens, sdk.DefaultPowerReduction), // PowerReduction should come from Genesis
				TotalDelegation: val.Tokens.BigInt(),
			})
		}
		if res.Pagination.Total >= vq.Options.Pagination.Limit {
			vq.Options.Pagination.Key = res.Pagination.NextKey
		} else {
			break
		}
	}

	return validators, nil
}

func (c *Client) GetAccAddress(valAddr string) (string, error) {
	valAddress, err := c.client.DecodeBech32ValAddr(valAddr)
	if err != nil {
		return "", err
	}
	return c.client.EncodeBech32AccAddr(sdk.AccAddress(valAddress))
}

func (c *Client) GetDelegators(valAddr string) ([]Delegator, error) {
	dq := c.getNewQuerier()
	delegators := make([]Delegator, 0)
	fDefaultPowerReduction, _ := sdk.DefaultPowerReduction.ToDec().Float64() // PowerReduction should come from Genesis

	for {
		dres, err := dq.Staking_ValidatorDelegations(valAddr)
		if err != nil {
			return nil, err
		}

		for _, delRes := range dres.DelegationResponses {
			fDelToken, _ := delRes.Balance.Amount.ToDec().Float64()
			delegators = append(delegators, Delegator{
				ValAddr:     delRes.Delegation.ValidatorAddress,
				DelAddr:     delRes.Delegation.DelegatorAddress,
				Token:       delRes.Balance.Amount.BigInt(),
				VotingPower: fDelToken / fDefaultPowerReduction,
			})
		}
		if dres.Pagination.Total >= dq.Options.Pagination.Limit {
			dq.Options.Pagination.Key = dres.Pagination.NextKey
		} else {
			break
		}
	}

	return delegators, nil
}
