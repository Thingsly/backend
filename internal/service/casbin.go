package service

import global "github.com/HustIoTPlatform/backend/pkg/global"

type Casbin struct {
}


func (*Casbin) AddFunctionToRole(role string, functions []string) bool {
	var rules [][]string
	for _, function := range functions {
		rule := []string{role, function, "allow"}
		rules = append(rules, rule)
	}
	isSuccess, _ := global.CasbinEnforcer.AddNamedPolicies("p", rules)
	return isSuccess
}


func (*Casbin) GetFunctionFromRole(role string) ([]string, bool) {
	policys := global.CasbinEnforcer.GetFilteredPolicy(0, role)
	var functions []string
	for _, policy := range policys {
		functions = append(functions, policy[1])
	}
	return functions, true
}


func (*Casbin) RemoveRoleAndFunction(role string) bool {
	isSuccess, _ := global.CasbinEnforcer.RemoveFilteredPolicy(0, role)
	return isSuccess

}


func (*Casbin) AddRolesToUser(user string, roles []string) bool {
	var rules [][]string
	for _, role := range roles {
		rule := []string{user, role}
		rules = append(rules, rule)
	}
	isSuccess, _ := global.CasbinEnforcer.AddNamedGroupingPolicies("g", rules)
	return isSuccess
}


func (*Casbin) GetRoleFromUser(user string) ([]string, bool) {
	policys := global.CasbinEnforcer.GetFilteredNamedGroupingPolicy("g", 0, user)
	var roles []string
	for _, policy := range policys {
		roles = append(roles, policy[1])
	}
	return roles, true
}


func (*Casbin) RemoveUserAndRole(user string) bool {
	isSuccess, _ := global.CasbinEnforcer.RemoveFilteredNamedGroupingPolicy("g", 0, user)
	return isSuccess
}


func (*Casbin) GetUrl(url string) bool {
	stringList := global.CasbinEnforcer.GetFilteredNamedGroupingPolicy("g2", 0, url)
	return len(stringList) != 0
}


func (*Casbin) HasRole(role string) bool {
	stringList := global.CasbinEnforcer.GetFilteredNamedGroupingPolicy("g", 1, role)
	return len(stringList) != 0
}


func (*Casbin) Verify(user string, url string) bool {
	isTrue, _ := global.CasbinEnforcer.Enforce(user, url, "allow")
	return isTrue
}
