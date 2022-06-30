package subspaces

import (
	"context"

	subspacestypes "github.com/desmos-labs/desmos/v3/x/subspaces/types"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/desmos-labs/djuno/v2/types"
)

// updateSubspace updates the stored data for the given subspace at the specified height
func (m *Module) updateSubspace(height int64, subspaceID uint64) error {
	// Get the subspace
	res, err := m.subspacesClient.Subspace(
		remote.GetHeightRequestContext(context.Background(), height),
		subspacestypes.NewQuerySubspaceRequest(subspaceID),
	)
	if err != nil {
		return err
	}

	// Save the subspace
	return m.db.SaveSubspace(types.NewSubspace(res.Subspace, height))
}

// updateUserGroup updates the stored data for the given user group at the specified height
func (m *Module) updateUserGroup(height int64, subspaceID uint64, groupID uint32) error {
	// Get the user group
	res, err := m.subspacesClient.UserGroup(
		remote.GetHeightRequestContext(context.Background(), height),
		subspacestypes.NewQueryUserGroupRequest(subspaceID, groupID),
	)
	if err != nil {
		return err
	}

	// Save the user group
	return m.db.SaveUserGroup(types.NewUserGroup(res.Group, height))
}

// updateUserPermissions updates the stored permissions for the given user at the specified height
func (m *Module) updateUserPermissions(height int64, subspaceID uint64, user string) error {
	// Get the permissions
	res, err := m.subspacesClient.UserPermissions(
		remote.GetHeightRequestContext(context.Background(), height),
		subspacestypes.NewQueryUserPermissionsRequest(subspaceID, user),
	)
	if err != nil {
		return err
	}

	// Save the user permissions
	return m.db.SaveUserPermission(types.NewUserPermission(
		subspacestypes.NewACLEntry(subspaceID, user, res.Permissions),
		height,
	))
}