package database_test

import (
	"time"

	"github.com/desmos-labs/djuno/x/profiles/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

func (suite *DbTestSuite) TestDesmosDb_SaveUserIfNotExisting() {
	err := suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", 1)
	suite.Require().NoError(err, "storing of address should return no error")

	err = suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", 1)
	suite.Require().NoError(err, "storing address second time should return no error")

	user, err := suite.database.GetUserByAddress("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	suite.Require().Equal("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", user.GetAddress().String())
}

func (suite *DbTestSuite) verifyEqual(expected, actual *profilestypes.Profile) {
	suite.Require().Equal(expected.Account, actual.Account)
	suite.Require().Equal(expected.DTag, actual.DTag)
	suite.Require().Equal(expected.Moniker, actual.Moniker)
	suite.Require().Equal(expected.Bio, actual.Bio)
	suite.Require().Equal(expected.Pictures, actual.Pictures)
	suite.Require().True(expected.CreationDate.Equal(actual.CreationDate))
}

func (suite *DbTestSuite) TestDesmosDb_SaveProfile() {
	addr, err := sdk.AccAddressFromBech32("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)

	original, err := profilestypes.NewProfile(
		"original-moniker",
		"",
		"",
		profilestypes.NewPictures("", ""),
		time.Time{},
		authtypes.NewBaseAccountWithAddress(addr),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(original, 10))
	suite.Require().NoError(err)

	// Verify the storing
	stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(original, stored)

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	updated, err := original.Update(profilestypes.NewProfileUpdate(
		"new-dtag",
		"new-moniker",
		"new-bio",
		profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify)),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(updated, 9))
	suite.Require().NoError(err)

	// Verify the data
	stored, err = suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(original, stored)

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	updated, err = original.Update(profilestypes.NewProfileUpdate(
		"new-dtag",
		"new-moniker",
		"new-bio",
		profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify)),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(updated, 10))
	suite.Require().NoError(err)

	// Verify the data
	stored, err = suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(updated, stored)

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	updated, err = original.Update(profilestypes.NewProfileUpdate(
		"new-dtag-2",
		"new-moniker-2",
		"new-bio-2",
		profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify)),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(updated, 11))
	suite.Require().NoError(err)

	// Verify the data
	stored, err = suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(updated, stored)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveRelationship() types.Relationship {
	relationship := types.NewRelationship(
		profilestypes.NewRelationship(
			"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
			"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
			"mooncake",
		),
		10,
	)

	// Save the relationship
	err := suite.database.SaveRelationship(relationship)
	suite.Require().NoError(err)

	return relationship
}

func (suite *DbTestSuite) TestDesmosDb_SaveRelationship() {
	relationship := suite.saveRelationship()

	err := suite.database.SaveRelationship(relationship)
	suite.Require().NoError(err, "double inserting the same relationship should return no error")

	var rows []dbtypes.RelationshipRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM relationship")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewRelationshipRow(
		relationship.Creator,
		relationship.Recipient,
		relationship.Subspace,
		relationship.Height,
	)))
}

func (suite *DbTestSuite) TestDesmosDb_DeleteRelationship() {
	relationship := suite.saveRelationship()

	err := suite.database.DeleteRelationship(relationship)
	suite.Require().NoError(err, "removing existing relationship should return no error")

	var rows []dbtypes.RelationshipRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM relationship")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 0)

	err = suite.database.DeleteRelationship(relationship)
	suite.Require().NoError(err, "deleting non existent relationship should return no error")
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveBlockage() types.Blockage {
	suite.Require().NoError(suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1))
	suite.Require().NoError(suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1))

	blockage := types.NewBlockage(
		profilestypes.NewUserBlock(
			"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
			"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
			"this is my blocking reason",
			"mooncake"),
		1,
	)

	// Save the blockage
	err := suite.database.SaveBlockage(blockage)
	suite.Require().NoError(err)

	return blockage
}

func (suite *DbTestSuite) TestDesmosDB_SaveUserBlockage() {
	blockage := suite.saveBlockage()

	err := suite.database.SaveBlockage(blockage)
	suite.Require().NoError(err, "double inserting blockage should return no error")

	var rows []dbtypes.BlockageRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_block")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewBlockageRow(
		blockage.Blocker,
		blockage.Blocked,
		blockage.Reason,
		blockage.Subspace,
		blockage.Height,
	)))
}

func (suite *DbTestSuite) TestDesmosDB_RemoveUserBlockage() {
	blockage := suite.saveBlockage()

	err := suite.database.RemoveBlockage(blockage)
	suite.Require().NoError(err)

	var rows []dbtypes.BlockageRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_block")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 0)

	err = suite.database.RemoveBlockage(blockage)
	suite.Require().NoError(err, "deleting non existing blockage should return no error")
}