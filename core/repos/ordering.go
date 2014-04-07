package repos

import (
	r "github.com/dancannon/gorethink"
)

// var hoursSince = function(row) {
//   return r.epochTime(r.now().sub(row("Created"))).hours()
// };

// var score = function(row) {
//   return row("Likes").sub(row("Dislikes"))
// };

// r.db("news").table("posts").map(function(row) {
//   return score(row).div(hoursSince(row).add(2))
// })

func hoursSince(row r.RqlTerm) r.RqlTerm {
	return r.EpochTime(r.Now().Sub(row.Field("Created"))).Hours()
}

func score(row r.RqlTerm) r.RqlTerm {
	return row.Field("Likes").Sub(row.Field("Dislikes"))
}

// Generic sorting functions

func orderByPopular(row r.RqlTerm) r.RqlTerm {
	hours := hoursSince(row).Add(2)

	return score(row).Div(hours.Mul(hours))
}

func orderByTop(row r.RqlTerm) r.RqlTerm {
	return score(row)
}

func orderByNew(row r.RqlTerm) r.RqlTerm {
	return row.Field("Created")
}
