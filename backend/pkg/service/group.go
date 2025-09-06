package service

import (
	"fmt"

	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/repository"
	"asynclab.club/asynx/backend/pkg/security"
	"github.com/sirupsen/logrus"
)

type ServiceGroup struct {
	repositoryGroup *repository.RepositoryGroup
}

func NewServiceGroup(repo *repository.RepositoryGroup) *ServiceGroup {
	return &ServiceGroup{
		repositoryGroup: repo,
	}
}

func (r *ServiceGroup) FindByOuAndCn(ou security.OuGroup, cn string) (group *entity.Group, err error) {
	return r.repositoryGroup.FindByOuAndCn(ou.String(), cn)
}

func (s *ServiceGroup) FindAll() ([]*entity.Group, error) {
	return s.repositoryGroup.FindAll()
}

func (s *ServiceGroup) FindAllByOu(ou security.OuGroup) ([]*entity.Group, error) {
	return s.repositoryGroup.FindAllByOu(ou.String())
}

func (s *ServiceGroup) FindAllByOuAndMemberUid(ou security.OuGroup, uid string) ([]*entity.Group, error) {
	return s.repositoryGroup.FindAllByOuAndMemberUid(ou.String(), uid)
}

func (s *ServiceGroup) GetRoleByUid(uid string) (security.Role, error) {
	groups, err := s.FindAllByOuAndMemberUid(security.OuGroupSupplementary, uid)
	if err != nil {
		logrus.Error("Failed to get groups for user ", uid, ": ", err)
		return security.RoleAnonymous, err
	}
	if len(groups) == 0 {
		return security.RoleAnonymous, nil
	}
	role, err := security.GetRoleFromLdapGroups(groups)
	if err != nil {
		return security.RoleAnonymous, fmt.Errorf("error getting role for user %s: %w", uid, err)
	}
	return role, nil
}

func (s *ServiceGroup) GetRole(user *entity.User) (security.Role, error) {
	return s.GetRoleByUid(user.Uid)
}

func (s *ServiceGroup) GrantRoleByUid(uid string, newRole security.Role) error {
	oldRole, err := s.GetRoleByUid(uid)
	if err != nil {
		return err
	}

	// 如果角色没有变化，直接返回
	if oldRole == newRole {
		return nil
	}

	attr := map[string][]string{"memberUid": {uid}}

	// 如果新角色是匿名（移除所有角色）
	if newRole == security.RoleAnonymous {
		if oldRole != security.RoleAnonymous {
			oldGroup, err := s.FindByOuAndCn(security.OuGroupSupplementary, oldRole.String())
			if err != nil {
				return err
			}
			return s.repositoryGroup.Modify(s.repositoryGroup.BuildDn(oldGroup), nil, attr, nil)
		}
		return nil
	}

	// 查找新角色组
	newGroup, err := s.FindByOuAndCn(security.OuGroupSupplementary, newRole.String())
	if err != nil {
		return err
	}
	if newGroup == nil {
		return fmt.Errorf("group %s not found", newRole)
	}

	// 如果用户之前没有角色（直接添加）
	if oldRole == security.RoleAnonymous {
		return s.repositoryGroup.Modify(s.repositoryGroup.BuildDn(newGroup), attr, nil, nil) // 添加用户
	}

	// 如果是角色切换：先从旧组移除，再添加到新组
	oldGroup, err := s.FindByOuAndCn(security.OuGroupSupplementary, oldRole.String())
	if err != nil {
		return err
	}
	if err := s.repositoryGroup.Modify(s.repositoryGroup.BuildDn(oldGroup), nil, attr, nil); err != nil {
		return err
	}
	if err := s.repositoryGroup.Modify(s.repositoryGroup.BuildDn(newGroup), attr, nil, nil); err != nil {
		// 回滚
		if err = s.repositoryGroup.Modify(s.repositoryGroup.BuildDn(oldGroup), attr, nil, nil); err != nil {
			logrus.Warningf("Failed to rollback group modification when grant role: %v", err)
		}
		return err
	}

	return nil
}

func (s *ServiceGroup) GrantRole(user *entity.User, role security.Role) error {
	return s.GrantRoleByUid(user.Uid, role)
}
