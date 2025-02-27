package keeper

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/regen-network/gocuke"
	"gotest.tools/v3/assert"

	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/orm/model/ormdb"
	"github.com/cosmos/cosmos-sdk/orm/model/ormtable"
	"github.com/cosmos/cosmos-sdk/orm/testing/ormtest"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	api "github.com/regen-network/regen-ledger/api/v2/regen/ecocredit/basket/v1"
	baseapi "github.com/regen-network/regen-ledger/api/v2/regen/ecocredit/v1"
	"github.com/regen-network/regen-ledger/x/ecocredit/v3"
	"github.com/regen-network/regen-ledger/x/ecocredit/v3/mocks"
)

const (
	testClassID     = "C01"
	testBasketDenom = "eco.uC.NCT"
)

type baseSuite struct {
	t          gocuke.TestingT
	db         ormdb.ModuleDB
	ctx        context.Context
	k          Keeper
	ctrl       *gomock.Controller
	addrs      []sdk.AccAddress
	stateStore api.StateStore
	baseStore  baseapi.StateStore
	bankKeeper *mocks.MockBankKeeper
	storeKey   *storetypes.KVStoreKey
	sdkCtx     sdk.Context
}

func setupBase(t gocuke.TestingT) *baseSuite {
	// prepare database
	s := &baseSuite{t: t}
	var err error
	s.db, err = ormdb.NewModuleDB(&ecocredit.ModuleSchema, ormdb.ModuleDBOptions{})
	assert.NilError(t, err)
	s.stateStore, err = api.NewStateStore(s.db)
	assert.NilError(t, err)
	s.baseStore, err = baseapi.NewStateStore(s.db)
	assert.NilError(t, err)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	s.storeKey = sdk.NewKVStoreKey("test")
	cms.MountStoreWithDB(s.storeKey, storetypes.StoreTypeIAVL, db)
	assert.NilError(t, cms.LoadLatestVersion())
	ormCtx := ormtable.WrapContextDefault(ormtest.NewMemoryBackend())
	s.sdkCtx = sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger()).WithContext(ormCtx)
	s.ctx = sdk.WrapSDKContext(s.sdkCtx)

	// setup test keeper
	s.ctrl = gomock.NewController(t)
	assert.NilError(t, err)
	s.bankKeeper = mocks.NewMockBankKeeper(s.ctrl)

	_, _, moduleAddress := testdata.KeyTestPubAddr()
	authority, err := sdk.AccAddressFromBech32("regen1nzh226hxrsvf4k69sa8v0nfuzx5vgwkczk8j68")
	assert.NilError(t, err)

	s.k = NewKeeper(s.stateStore, s.baseStore, s.bankKeeper, moduleAddress, authority)
	s.baseStore, err = baseapi.NewStateStore(s.db)
	assert.NilError(t, err)

	// add test addresses
	_, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	s.addrs = append(s.addrs, addr1, addr2)

	return s
}
