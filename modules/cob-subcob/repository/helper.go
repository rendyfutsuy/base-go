package repository

import "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"

func fromFilterToWhere(filter []request.Filter) (where string) {
	if len(filter) == 0 {
		return ""
	}

	where = ""
	for i, f := range filter {
		// validate filter operator

		if i > 0 {
			where += " AND "
		}

		where += f.Option + " " + f.Operator + " '" + f.Value + "'"
	}

	return
}


