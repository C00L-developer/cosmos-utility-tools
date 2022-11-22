# Chain Upgrades

## Cosmos Hub

### 1. Delta Upgrade (2021-07-12)

#### Context

    In July 2020, the Iqlusion team developed ATOM 2021 to drive the direction of the Cosmos Hub after the completion of the Cosmos whitepaper and IBC. It became clear that providing liquidity to new IBC connected zones was core the Hub's mission.

    Tendermint and B-Harvest joined forces to produce and develop a Liquidity Module. In 2021 March, they submitted a signal governance proposal to ask the Atom delegator community about Gravity DEX (Liquidity Module) adoption on the Cosmos Hub. Prop38 was very well approved by the community.

    This proposal completes the first leg of ATOM 2021 and achieves the goals of the signaling proposal by bringing an IBC compatible DEX to the Hub.

#### Updates

- Gaia v4.2.1 -> v5.0.0
- Gravity DEX:
    - A scalable AMM model for token swaps
    - Drives liquidity for tokens on the Cosmos Hub
    - Delivers price consistency and order execution

#### Metadata Changes

There are no metadata changes. The chain-id remain the same, `cosmoshb-4`.

#### [Code Changes](https://github.com/cosmos/gaia/releases/tag/v5.0.0)

- (golang) Bump golang prerequisite from 1.15 to 1.16.
- (gaia) Add Liquidity module v1.2.9.
- (sdk) Bump SDK version to v0.42.6.
- (tendermint) Bump Tendermint version to v0.34.11.

#### Upgrade steps

1. Prior to the upgrade, operators MUST be running Gaia v4.2.1.
2. At the upgrade block height of 6910000, the Gaia software will panic.
3. Important note to all validators: Although the upgrade path is essentially to replace the binary when the software panics and halts at the upgrade height, an important disaster recovery operation is to take a snapshot of your state after the halt and before starting v5.0.0.
4. Replace the Gaia v4.2.1 binary with the Gaia v5.0.0 binary
5. Start the Gaia v5.0.0 binary using the following command

    ```cmd
    gaiad start --x-crisis-skip-assert-invariants
    ```
6. Wait until 2/3+ of voting power has upgraded for the network to start producing blocks


### 2. Vega Upgrade (2021-12-14)

#### Context

    Bump `Cosmos-SDK` to `v0.44.3` which notably includes fixes for the vesting accounts and two new modules outlined below. 
    Add the `authz` module to the Cosmos Hub. x/authz is an implementation of a Cosmos SDK module, per `ADR30` , that allows granting arbitrary privileges from one account (the granter) to another account (the grantee). Authorizations must be granted for a particular Msg service method one by one using an implementation of the Authorization interface.
    Add the  `feegrant` module to the Cosmos Hub. This module allows accounts to grant fee allowances and to use fees from their accounts. Grantees can execute any transaction without the need to maintain sufficient fees.

    Add `IBC` as a standalone module from the Cosmos SDK using version `v2.0.0.
    Please note that the governance parameter for `MaxExpectedBlockDelay` is set to 30 seconds. This means if a connection is opened with a packet delay of 1 minute, it requires 2 blocks to be committed after the consensus state is submitted before the packet can be processed.
    
    Add `packet-forward-middleware v1.0.1`  as per (https://www.mintscan.io/cosmos/proposals/56) prepared. This feature allows multi-hop IBC transfer messages so that a user can send tokens from chain A to chain C via chain B.
    
    Bump `Liquidity` module to `v1.4.2` to ensure compatibility with Cosmos-SDK.

#### Updates

- Gaia v6.0.x
- Cosmos SDK v0.44
    - Fee grant module:
        - Allows paying fees on behalf of another account
    - Authz module:
        - Provide governance functions to execute transactions on behalf of another account
- Liquidity Module v1.4.2
    - The Gravity DEX with updates for dependencies
- IBC v2.0.0
- Tendermint v0.34.14
- Cosmosvisor v0.1.0
- IBC packet forward middleware v1.0.1
    - Cosmos Hub as a router

- External chain launch: Gravity Bridge
    - Transfer ATOM, ETH, ERC-20, and other Cosmos tokens between Ethereum and the Gravity Bridge Chain and by extension all IBC connected chains.
    - Fee and reward model hosted across Cosmos and Ethereum

#### Metadata Changes

There are no metadata changes. The chain-id remain the same, `cosmoshb-4`.

#### [Code Changes](https://github.com/cosmos/gaia/releases/tag/v6.0.0)

- (gaia) Add NewSetUpContextDecorator to anteDecorators
- (gaia) Reconfigure SetUpgradeHandler to ensure vesting is configured after auth and new modules have InitGenesis run.
- (golang) Bump golang prerequisite to 1.17. 
- (gaia) Bump [Liquidity](https://github.com/gravity-devs/liquidity) module to [v1.4.2](https://github.com/Gravity-Devs/liquidity/releases/tag/v1.4.2).
- (gaia) Bump [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) to [v0.44.3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.3). See the [CHANGELOG.md](https://github.com/cosmos/cosmos-sdk/blob/release/v0.44.x/CHANGELOG.md#v0443---2021-10-21) for details.
- (gaia) Add [IBC](https://github.com/cosmos/ibc-go) as a standalone module from the Cosmos SDK using version [v2.0.0](https://github.com/cosmos/ibc-go/releases/tag/v2.0.0). See the [CHANGELOG.md](https://github.com/cosmos/ibc-go/blob/v2.0.0/CHANGELOG.md) for details.
- (gaia) Add [packet-forward-middleware](https://github.com/strangelove-ventures/packet-forward-middleware) [v1.0.1](https://github.com/strangelove-ventures/packet-forward-middleware/releases/tag/v1.0.1).
- (gaia) [#969](https://github.com/cosmos/gaia/issues/969) Remove legacy migration code.

#### Upgrade steps

- Manual Upgrade 

    Run Gaia v5.0.x till upgrade height, the node will panic:
    ```
    ERR UPGRADE "Vega" NEEDED at height: 8695000
    panic: UPGRADE "Vega" NEEDED at height: 8695000
    ```
    Stop the node, and install Gaia v6.0.4 and re-start by gaiad start.

    It may take 20 min to a few hours until validators with a total sum voting power > 2/3 to complete their nodes upgrades. After that, the chain can continue to produce blocks.

- Upgrade using Cosmovisor by manually preparing the Gaia v6.0.4 binary

    1. Install the latest version of Cosmovisor
    2. Create a cosmovisor folder, create a Cosmovisor folder inside $GAIA_HOME and move Gaia v5.0.0 into $GAIA_HOME/cosmovisor/genesis/bin, build Gaia v6.0.4, and move gaiad v6.0.4 to $GAIA_HOME/cosmovisor/upgrades/Vega/bin
    3. Export the environmental variables, and start the node:
        ```cmd
        export DAEMON_NAME=gaiad
        # please change to your own gaia home dir
        export DAEMON_HOME= $GAIA_HOME
        export DAEMON_RESTART_AFTER_UPGRADE=true

        cosmovisor start --x-crisis-skip-assert-invariants
        ```

### 3. Theta Upgrade (2021-04-12)

#### Updates

- Gaia v7.0.x
- Cosmos SDK v0.45
    - Minimal update with small fixes
- Gravity DEX: Liquidity v1.4.5
    - Adds a circuit breaker governance proposal type to disable adding new liquidity in order to make a migration possible.
- IBC 3.0.0
    - Interchain Account Module
        - Allows the creation of accounts on a "Host" blockchain which are controlled by an authentication module on a "Controller" blockchain.
        - Arbitrary messages are able to be submitted from the "Controller" blockchain to the "Host" blockchain to be executed on behalf of the Interchain Account.
        - Uses ordered IBC channels, one per account.

#### Metadata Changes

There are no metadata changes. The chain-id remain the same, `cosmoshb-4`.

#### [Code Changes](https://github.com/cosmos/gaia/releases/tag/v6.0.0)

- (gaia) Add NewSetUpContextDecorator to anteDecorators
- (gaia) Reconfigure SetUpgradeHandler to ensure vesting is configured after auth and new modules have InitGenesis run.
- (golang) Bump golang prerequisite to 1.17. 
- (gaia) Bump [Liquidity](https://github.com/gravity-devs/liquidity) module to [v1.4.2](https://github.com/Gravity-Devs/liquidity/releases/tag/v1.4.2).
- (gaia) Bump [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) to [v0.44.3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.3). See the [CHANGELOG.md](https://github.com/cosmos/cosmos-sdk/blob/release/v0.44.x/CHANGELOG.md#v0443---2021-10-21) for details.
- (gaia) Add [IBC](https://github.com/cosmos/ibc-go) as a standalone module from the Cosmos SDK using version [v2.0.0](https://github.com/cosmos/ibc-go/releases/tag/v2.0.0). See the [CHANGELOG.md](https://github.com/cosmos/ibc-go/blob/v2.0.0/CHANGELOG.md) for details.
- (gaia) Add [packet-forward-middleware](https://github.com/strangelove-ventures/packet-forward-middleware) [v1.0.1](https://github.com/strangelove-ventures/packet-forward-middleware/releases/tag/v1.0.1).
- (gaia) [#969](https://github.com/cosmos/gaia/issues/969) Remove legacy migration code.

#### Upgrade steps

- Manual Upgrade 

    Run Gaia v6.0.x till upgrade height, the node will panic:
    ```
    ERR UPGRADE "v7-Theta" NEEDED at height: 10085397
    panic: UPGRADE "v7-Theta" NEEDED at height: 10085397
    ```
    Stop the node, and install Gaia v7.0.0 and re-start by gaiad start.

    It may take 7 minutes to a few hours until validators with a total sum voting power > 2/3 to complete their nodes upgrades. After that, the chain can continue to produce blocks.

- Upgrade using Cosmovisor by manually preparing the Gaia v7.0.0 binary

    1. Install the latest version of Cosmovisor
    2. Create a cosmovisor folder, create a Cosmovisor folder inside $GAIA_HOME and move Gaia v5.0.0 into $GAIA_HOME/cosmovisor/genesis/bin, build Gaia v6.0.4, and move gaiad v7.0.0 to `$GAIA_HOME/cosmovisor/upgrades/v7-Theta/bin`
    3. Export the environmental variables, and start the node:
        ```cmd
        export DAEMON_NAME=gaiad
        # please change to your own gaia home dir
        export DAEMON_HOME= $GAIA_HOME
        export DAEMON_RESTART_AFTER_UPGRADE=true

        cosmovisor start --x-crisis-skip-assert-invariants
        ```
### Links

- https://github.com/cosmos/gaia/blob/main/docs/roadmap/cosmos-hub-roadmap-2.0.md
- https://github.com/cosmos/gaia/blob/main/CHANGELOG.md
- https://github.com/cosmos/gaia/tree/main/docs/migration
- https://www.mintscan.io/cosmos/proposals/51
- https://www.mintscan.io/cosmos/proposals/59
- https://www.mintscan.io/cosmos/proposals/65

## Osmosis

### v10 Upgrade (2022-06-12)

#### Context 

    The v10.0.0 release fixes the JoinPool bug present in v9 of the Osmosis code.

#### Upgrade & Code Changes

- Breaking Chainges

    * [#1699](https://github.com/osmosis-labs/osmosis/pull/1699) Fixes bug in sig fig rounding on spot price queries for small values
    * [#1671](https://github.com/osmosis-labs/osmosis/pull/1671) Remove methods that constitute AppModuleSimulation APIs for several modules' AppModules, which implemented no-ops
    * [#1671](https://github.com/osmosis-labs/osmosis/pull/1671) Add hourly epochs to `x/epochs` DefaultGenesis.
    * [#1665](https://github.com/osmosis-labs/osmosis/pull/1665) Delete app/App interface, instead use simapp.App
    * [#1630](https://github.com/osmosis-labs/osmosis/pull/1630) Delete the v043_temp module, now that we're on an updated SDK version.

- Bug Fixes

    * [1700](https://github.com/osmosis-labs/osmosis/pull/1700) Upgrade sdk fork with missing snapshot manager fix.
    * [1716](https://github.com/osmosis-labs/osmosis/pull/1716) Fix secondary over-LP shares bug with uneven swap amounts in `CalcJoinPoolShares`.
    * [1759](https://github.com/osmosis-labs/osmosis/pull/1759) Fix pagination filter in incentives query.
    * [1698](https://github.com/osmosis-labs/osmosis/pull/1698) Register wasm snapshotter extension.
    * [1931](https://github.com/osmosis-labs/osmosis/pull/1931) Add explicit check for input denoms to `CalcJoinPoolShares`

#### Upgrade Steps

All validator nodes should upgrade to v10 prior to the network restarting. The v10 binary is state machine compatible with v9 until block 4713065. At 4:00PM UTC on June 12th, 2022, we will have a coordinated re-start of the network. The sequence of events will look like the following:

- All validator nodes upgrade to v10 now, but keep their nodes offline. Even if your node is further behind (i.e. you stopped your node first early within the shutdown and still have blocks ahead of you before reaching the halt height), you still must upgrade to v10 now. v10 will run the v9 state machine until the predefined "fork block height"
- At exactly 4:00PM UTC on June 12th, 2022, all validators start their nodes at the same time
- Once 67% or more of the voting power gets online, block 4713065 will be reached, along with the upgrade at this height. Prior to 67 percent of validator power getting online, you will only see p2p logs. This is also an epoch block, so it will take some time to process
- After block 4713065, three more epochs will happen back to back, one per block.
- If the June 12th epoch time has not occured yet, blocks will be produced until the epoch time. If the epoch time has occured, the June 12th epoch will occur in conjunction with the four other epochs above.

The coordination of restart will happen over Discord. In the event Discord is down, validators should form a Telegram group to further coordinate the network restart.

- Cosmovisor: Manual Method

    1. Set the following env variables:

        ```cmd
        echo "# Setup Cosmovisor" >> ~/.profile
        echo "export DAEMON_NAME=osmosisd" >> ~/.profile
        echo "export DAEMON_HOME=$HOME/.osmosisd" >> ~/.profile
        echo "export DAEMON_ALLOW_DOWNLOAD_BINARIES=false" >> ~/.profile
        echo "export DAEMON_LOG_BUFFER_SIZE=512" >> ~/.profile
        echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.profile
        echo "export UNSAFE_SKIP_BACKUP=true" >> ~/.profile
        source ~/.profile
        ```

    2. Create the required folder, make the build, and copy the daemon over to that folder. NOTE, you must put the v10 binary in the v9 folder as shown below 
    since this is a fork.

        ```cmd
        mkdir -p ~/.osmosisd/cosmovisor/upgrades/v9/bin
        cd $HOME/osmosis
        git pull
        git checkout v10.0.0
        make build
        cp build/osmosisd ~/.osmosisd/cosmovisor/upgrades/v9/bin
        ```

- Completely Manual Option

    ```cmd
    cd $HOME/osmosis
    git pull
    git checkout v10.0.0
    make install
    ```

### v11 Scambuster Upgrade (2022-08-03)

#### Context

    This upgrade contains only two minor changes aimed at reducing on-chain spam to protect Osmosis users. 

    On July 22, an attacker attempted to halt the Osmosis chain by creating over 20,000 external incentives gauges in various pools in an attempt to overload Validator nodes at the Osmosis epoch block.

    The attack was ultimately unsuccessful and the epoch block processed as normal (subject to a < 4 minute delay). Though unsuccessful, this attack did result in most Osmosis users being dusted with small amounts of tokens that they could not easily dispose of. From a user perspective, this made for a rather annoying and less than ideal user experience.

    To prevent this from happening in the future, this upgrade allows for messages to be assigned a fee which, if not paid, will cause the message to fail. Using this new feature, the upgrade assigns fees to creating external incentives gauges and adding tokens to them as follows:
    ```cmd
    CreateGauge: 50 OSMO
    AddToGauge: 25 OSMO
    ```
    These fees will be treated as the minimum gas fees required for these transactions to process, meaning that the fees will be distributed to OSMO stakers when paid by the gauge creator. After this upgrade is implemented, gauge creation will likely be considered too expensive for any potential attacker to create thousands of them for the sole purpose of spamming the network.

    Over the past weeks, numerous chains in the Cosmos ecosystem have fallen victim to scam proposals that were placed in the deposit period and contained links to malicious sites that have attempted to steal a user’s seed phrase. 

    To prevent further proposal spam, the v11 upgrade will require wallets that submit a governance proposal to also pay 25% of the deposit at the time of submission.
    

#### Updates

- Minimum fee for messages

    Implements the ability for messages to be assigned a minimum fee, without which the transaction will fail.
    This minimum fee is initially implemented on:
    ```cmd
    CreateGauge - 50 OSMO
    AddToGauge - 25 OSMO
    ```

- Governance proposal submission deposit

    Submitting a governance proposal now requires a percentage of the deposit to be filled by the proposing wallet.
    This is set with the governance parameter min_initial_deposit_percent and is initially implemented as 25%.

#### Code Changes

- Improvements
    * [#2237](https://github.com/osmosis-labs/osmosis/pull/2237) Enable charging fee in base denom for `CreateGauge` and `AddToGauge`.
    * [#2214](https://github.com/osmosis-labs/osmosis/pull/2214) Speedup epoch distribution, superfluid component
    * [#2130](https://github.com/osmosis-labs/osmosis/pull/2130) Introduce errors in mint types.
    * [#2000](https://github.com/osmosis-labs/osmosis/pull/2000) Update import paths from v9 to v10.
- Bug Fixes & API Breaks
    * [2011](https://github.com/osmosis-labs/osmosis/pull/2011) Fix bug in TokenFactory initGenesis, relating to denom creation fee param.
    * Restores vesting by duration command
    * Fixes pagination in x/incentives module queries
    * [#1937](https://github.com/osmosis-labs/osmosis/pull/1937) Change `lockupKeeper.ExtendLock` to take in lockID instead of the direct lock struct.
    * [#2030](https://github.com/osmosis-labs/osmosis/pull/2030) Rename lockup keeper `ResetAllLocks` to `InitializeAllLocks` and `ResetAllSyntheticLocks` to `InitializeAllSyntheticLocks`.
- SDK Upgrades
    * [#2245](https://github.com/osmosis-labs/osmosis/pull/2245) Upgrade SDK for to v0.45.0x-osmo-v9.2. Major changes:

#### Upgrade Steps

- Install and setup Cosmovisor

    1. install Cosmovisor and create the necessary folders for cosmosvisor in your daemon home directory (~/.osmosisd).

        ```cmd
        mkdir -p ~/.osmosisd
        mkdir -p ~/.osmosisd/cosmovisor
        mkdir -p ~/.osmosisd/cosmovisor/genesis
        mkdir -p ~/.osmosisd/cosmovisor/genesis/bin
        mkdir -p ~/.osmosisd/cosmovisor/upgrades
        ```
    
    2. Copy the current osmosisd binary into the cosmovisor/genesis folder and v9 folder.

        ```cmd
        cp $GOPATH/bin/osmosisd ~/.osmosisd/cosmovisor/genesis/bin
        mkdir -p ~/.osmosisd/cosmovisor/upgrades/v9/bin
        cp $GOPATH/bin/osmosisd ~/.osmosisd/cosmovisor/upgrades/v9/bin
        ```

    3. Set these environment variables:

        ```cmd
        echo "# Setup Cosmovisor" >> ~/.profile
        echo "export DAEMON_NAME=osmosisd" >> ~/.profile
        echo "export DAEMON_HOME=$HOME/.osmosisd" >> ~/.profile
        echo "export DAEMON_ALLOW_DOWNLOAD_BINARIES=false" >> ~/.profile
        echo "export DAEMON_LOG_BUFFER_SIZE=512" >> ~/.profile
        echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.profile
        echo "export UNSAFE_SKIP_BACKUP=true" >> ~/.profile
        source ~/.profile
        ```
    4. Create the required folder, make the build, and copy the daemon over to that folder

        ```cmd
        mkdir -p ~/.osmosisd/cosmovisor/upgrades/v11/bin
        cd $HOME/osmosis
        git pull
        git checkout v11.0.0
        make build
        cp build/osmosisd ~/.osmosisd/cosmovisor/upgrades/v11/bin
        ```

- Manual Option

    1. Wait for Osmosis to reach the upgrade height (5432450)
    2. Look for a panic message, followed by endless peer logs. Stop the daemon
    3. Run the following commands:
        ```cmd
        cd $HOME/osmosis
        git pull
        git checkout v11.0.0
        make install
        ```
    4. Start the osmosis daemon again, watch the upgrade happen, and then continue to hit blocks

### v12 Oxygen Upgrade (2022-09-30)

#### Context

    Between time weighted average pricing, CosmWasm development tooling, interchain accounts, and governance changes, the v12.0.0 upgrade may not be the “sexiest” upgrade Osmosis has seen to date, but the impact of v12 on the chain’s ability to serve as the base-layer for a massive DeFi ecosystem will be felt for years to come. Read on to learn more.

    NOTE — Due to variances in block time the upgrade height may occur sooner or later than the estimated time. You can follow a countdown of the time remaining until the upgrade height here. During the upgrade, you will be unable to use the Osmosis DEX. Staking and governance features will also be unusable until the upgrade is complete. A full changelog of upgrade features can be found here.

    - Time Weighted Average Pricing

    Following this upgrade, time weighted average pricing (TWAP) will be utilized to calculate asset prices in all new and existing liquidity pools on Osmosis moving forward.

    At a high level, TWAP involves storing cumulative pricing data for a given trading pair on chain, which can then be used to derive an average price over a selected time interval for that pair. To learn more about how exactly this is done, take a look at this insightful article by Osmosis contributor StevieWoofWoof.

    The impact of TWAP implementation is a pricing mechanism that is far less susceptible to manipulation on a block-by-block basis. This is an essential prerequisite for the launch of lending and margin trading protocols on Osmosis. If these protocols were to launch on Osmosis without TWAP, actors with substantial financial means could manipulate asset pricing at will to profit from liquidations or short interests.

    - Stargate Queries

    In v12.0.0, select CosmWasm queries, known as “Stargate Queries” will be enabled on-chain. Stargate queries will enable CosmWasm contract developers on Osmosis to write and deploy contracts that are capable of fetching data from chains connected to Osmosis via IBC. While the impact of this feature will not be felt immediately, enabling Stargate queries is the first step in allowing for cross-chain smart contracting between Osmosis and other chains.

    - Interchain Accounts Module

    If this section sounds familiar, it’s because the Interchain Accounts module has been covered in a previous upgrade article. This module was implemented with the Osmosis v9.0.0 Nitrogen upgrade, but was never enabled. With the v12.0.0 upgrade, this module will finally be initialized, enabling a host of new use cases for inter-blockchain communication. For more details on Interchain Accounts and how they will revolutionize IBC, see the v9.0.0 upgrade article.

    Many chains and d’apps have been waiting for Osmosis to enable this module to interact with or launch on Osmosis and, with this upgrade, all of these teams will finally be able to take the final steps to do so. Protocols like Stride and Quicksilver, for example, have been waiting for this module to be enabled in order to bring their liquid staking services to Osmosis. Yes, you heard that correctly. Shortly after this upgrade we will have liquid staked OSMO from multiple liquid staking providers! stOSMO and qOSMO will soon be used to provide liquidity on Osmosis.

#### Updates

- TWAP

    Time weighted average prices for all AMM pools have been implemented.
    These expose an on-chain oracle for each pool allowing CosmWasm apps to use this data.

- CosmWasm contract developer features
    * Enables select queries for cosmwasm contracts.
    * Adds message responses to gamm messages, to remove the neccessity of bindings.
    * Allow specifying denom metadata from tokenfactory.
    * Upgrade to wasmd v0.28.0

- Interchain Accounts Enabled

    Allows accounts on another ICA enabled chain to carry out transactions on Osmosis and vice versa.

- Expedited Proposals

    Allows governance to have proposals that execute in faster time windows given higher quorums (2/3rds).

#### Code Changes

- TWAP - Time weighted average prices for all AMM pools
- Cosmwasm contract developer facing features
  * Enabling select queries for cosmwasm contracts
  * Add message responses to gamm messages, to remove the neccessity of bindings
  * Allow specifying denom metadata from tokenfactory
- Enabling Interchain accounts (for real this time)
- Upgrading IBC to v3.3.0
- Consistently makes authz work with ledger for all messages

    The release also contains the following changes affecting Osmosis users and node operators

- Fixing State Sync
- Enabling expedited proposals

    This upgrade also adds a number of safety and API boundary improving changes to the codebase.
    While not state machine breaking, this release also includes the revamped Osmosis simulator,
    which acts as a fuzz testing tool tailored for the SDK state machine.

- Breaking Changes

    * [#2477](https://github.com/osmosis-labs/osmosis/pull/2477) Tokenfactory burn msg clash with sdk
    * TypeMsgBurn: from "burn" to "tf_burn"
    * TypeMsgMint: from "mint" to "tf_mint"
    * [#2222](https://github.com/osmosis-labs/osmosis/pull/2222) Add scaling factors to MsgCreateStableswapPool
    * [#1889](https://github.com/osmosis-labs/osmosis/pull/1825) Add proto responses to gamm LP messages:
    * MsgJoinPoolResponse: share_out_amount and token_in fields 
    * MsgExitPoolResponse: token_out field 
    * [#1825](https://github.com/osmosis-labs/osmosis/pull/1825) Fixes Interchain Accounts (host side) by adding it to AppModuleBasics
    * [#1994](https://github.com/osmosis-labs/osmosis/pull/1994) Removed bech32ibc module
    * [#2016](https://github.com/osmosis-labs/osmosis/pull/2016) Add fixed 10000 gas cost for each Balancer swap
    * [#2193](https://github.com/osmosis-labs/osmosis/pull/2193) Add TwapKeeper to the Osmosis app
    * [#2227](https://github.com/osmosis-labs/osmosis/pull/2227) Enable charging fee in base denom for `CreateGauge` and `AddToGauge`.
    * [#2283](https://github.com/osmosis-labs/osmosis/pull/2283) x/incentives: refactor `CreateGauge` and `AddToGauge` fees to use txfees denom
    * [#2206](https://github.com/osmosis-labs/osmosis/pull/2283) Register all Amino interfaces and concrete types on the authz Amino codec. This will allow the authz module to properly serialize and de-serializes instances using Amino.
    * [#2405](https://github.com/osmosis-labs/osmosis/pull/2405) Make SpotPrice have a max value of 2^160, and no longer be able to panic
    * [#2473](https://github.com/osmosis-labs/osmosis/pull/2473) x/superfluid `AddNewSuperfluidAsset` now returns error, if any occurs instead of ignoring it.
    * [#2714](https://github.com/osmosis-labs/osmosis/pull/2714) Upgrade wasmd to v0.28.0.
    * Remove x/Bech32IBC

- Golang API breaks

    * [#2160](https://github.com/osmosis-labs/osmosis/pull/2160) Clean up GAMM keeper (move `x/gamm/keeper/params.go` contents into `x/gamm/keeper/keeper.go`, replace all uses of `PoolNumber` with `PoolId`, move `SetStableSwapScalingFactors` to stableswap package, and delete marshal_bench_test.go and grpc_query_internal_test.go)
    * [#1987](https://github.com/osmosis-labs/osmosis/pull/1987) Remove `GammKeeper.GetNextPoolNumberAndIncrement` in favor of the non-mutative `GammKeeper.GetNextPoolNumber`.
    * [#1667](https://github.com/osmosis-labs/osmosis/pull/1673) Move wasm-bindings code out of app package into its own root level package.
    * [#2013](https://github.com/osmosis-labs/osmosis/pull/2013) Make `SetParams`, `SetPool`, `SetTotalLiquidity`, and `SetDenomLiquidity` GAMM APIs private
    * [#1857](https://github.com/osmosis-labs/osmosis/pull/1857) x/mint rename GetLastHalvenEpochNum to GetLastReductionEpochNum
    * [#2133](https://github.com/osmosis-labs/osmosis/pull/2133) Add `JoinPoolNoSwap` and `CalcJoinPoolNoSwapShares` to GAMM pool interface and route `JoinPoolNoSwap` in pool_service.go to new method in pool interface
    * [#2353](https://github.com/osmosis-labs/osmosis/pull/2353) Re-enable stargate query via whitelsit
    * [#2394](https://github.com/osmosis-labs/osmosis/pull/2394) Remove unused interface methods from expected keepers of each module
    * [#2390](https://github.com/osmosis-labs/osmosis/pull/2390) x/mint remove unused mintCoins parameter from AfterDistributeMintedCoin
    * [#2418](https://github.com/osmosis-labs/osmosis/pull/2418) x/mint remove SetInitialSupplyOffsetDuringMigration from keeper
    * [#2417](https://github.com/osmosis-labs/osmosis/pull/2417) x/mint unexport keeper `SetLastReductionEpochNum`, `getLastReductionEpochNum`, `CreateDeveloperVestingModuleAccount`, and `MintCoins`
    * [#2587](https://github.com/osmosis-labs/osmosis/pull/2587) remove encoding config argument from NewOsmosisApp

- Features

    * [#2387](https://github.com/osmosis-labs/osmosis/pull/2387) Upgrade to IBC v3.2.0, which allows for sending/receiving IBC tokens with slashes.
    * [#1312] Stableswap: Createpool logic 
    * [#1230] Stableswap CFMM equations
    * [#1429] solver for multi-asset CFMM
    * [#1539] Superfluid: Combine superfluid and staking query on querying delegation by delegator
    * [#2223] Tokenfactory: Add SetMetadata functionality

- Bug Fixes

    * [#2086](https://github.com/osmosis-labs/osmosis/pull/2086) `ReplacePoolIncentivesProposal` ProposalType() returns correct value of `ProposalTypeReplacePoolIncentives` instead of `ProposalTypeUpdatePoolIncentives`
    * [1930](https://github.com/osmosis-labs/osmosis/pull/1930) Ensure you can't `JoinPoolNoSwap` tokens that are not in the pool
    * [2186](https://github.com/osmosis-labs/osmosis/pull/2186) Remove liquidity event that was emitted twice per message.

- Improvements

    * [#2515](https://github.com/osmosis-labs/osmosis/pull/2515) Emit events from functions implementing epoch hooks' `panicCatchingEpochHook` cacheCtx
    * [#2526](https://github.com/osmosis-labs/osmosis/pull/2526) EpochHooks interface methods (and hence modules implementing the hooks) return error instead of panic