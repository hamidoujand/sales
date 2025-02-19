package token_validation

default valid := false

valid if {
	valid_time
	valid_issuer
	valid_roles_format
}

valid_time if {
	input.token.exp > input.now
}

valid_issuer if {
	input.token.iss == "auth-service"
}

valid_roles_format if {
	input.token.roles
	type_name(input.token.roles) == "array"
	count(input.token.roles) > 0
	every role in input.token.roles {
		type_name(role) == "string"
	}
}
