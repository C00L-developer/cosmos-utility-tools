# Cosmos Utility Tools

## Distribution Analyze Tool

It is a tool to analyzes the validators and delegators distribution. It just creates 3 csv reports to review it in `Result` subdirectory.

- `{ChainName}-validator.csv` : the list of validators with moniker, voting power, self and total delegation
- `{ChainName}-delegator.csv` : the list of delegators with delegator and validator addresses, voting power
- `{ChainName}-multival-delegator.csv` : the list of delegators which stake to multiple validators

- Execute the tool

    ```cmd
    make validator-status
    ```

## Vesting Analyze Tool

It is a tool to undertand the vesting accounts unlock schedule and review the total supply.
It just analyzes the `genesis.json` file from the local path of `./Result` to reduce the downloading time due to genesis file is too large.
It reports the vesting account circulation and total supply in the time series to the `csv` file in `Result` subdirectory.

File name: `umee-vesting.csv`

- Execute the tool

    ```cmd
    make vesting-analyze
    ```

## [Reports about the chain upgrades](./Result/upgrade_analysis.md)
