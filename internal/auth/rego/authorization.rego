package role_validation

role_user := "USER"

role_admin := "ADMIN"

role_all := {role_admin, role_user}

default rule_any := false

rule_any if {
	# create a set from all input roles
	claim_roles := {role | some role in input.roles}

	# get the intersection between 2 sets
	matched_roles := role_all & claim_roles
	count(matched_roles) > 0
}

default rule_admin_only := false

rule_admin_only if {
	claim_roles := {role | some role in input.roles}
	matched_roles := {role_admin} & claim_roles
	count(matched_roles) > 0
}

default rule_user_only := false

rule_user_only if {
	claim_roles := {role | some role in input.roles}
	matched_roles := {role_user} & claim_roles
	count(matched_roles) > 0
}

default rule_admin_or_owner := false

rule_admin_or_owner if {
	claim_roles := {role | some role in input.roles}
	matched_roles := {role_admin} & claim_roles
	count(matched_roles) > 0
} else if {
	claim_roles := {role | some role in input.roles}
	matched_roles := {role_user} & claim_roles
	count(matched_roles) > 0

	# also must be the same userId as claim subject
	input.userId == input.subject
}
