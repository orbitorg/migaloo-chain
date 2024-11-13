package v4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	Denom        = "ibc/BC5C0BAFD19A5E4133FDA0F3E04AE1FBEE75A4A226554B2CBB021089FF2E1F8A"
	DeadContract = "migaloo1qelh4gv5drg3yhj282l6n84a6wrrz033kwyak3ee3syvqg3mu3msgphpk4"
	Foundation   = "migaloo10zqfqhw44e6gvu97frjzcghunndskhu40uyztwu00y6dr9qxrz6qcjfrf7"
)

// CreateUpgradeHandler that migrates the chain from v4.2.5 to v4.2.6
func CreateUpgradeHandler(
	mm *module.Manager,
	bankKeeper bankKeeper.Keeper,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// ignore the error if any
		if err := migrateFundFromDeadContacts(ctx, bankKeeper); err != nil {
			ctx.Logger().Error("migrateFundFromDeadContacts", "error", err)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// migrate fund from dead contacts to foundation
func migrateFundFromDeadContacts(
	ctx sdk.Context,
	bankKeeper bankKeeper.Keeper,
) error {
	deadContractAddr := sdk.MustAccAddressFromBech32(DeadContract)
	foundationAddr := sdk.MustAccAddressFromBech32(Foundation)

	// transfer token from dead contract to foundation

	// Get all balances from the dead contract
	allBalances := bankKeeper.GetAllBalances(ctx, deadContractAddr)
	if allBalances.IsZero() {
		return nil // Return early if no balances
	}

	return bankKeeper.SendCoins(ctx, deadContractAddr, foundationAddr, allBalances)
}
