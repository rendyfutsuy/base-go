package request

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplySearchConditionWithSubqueriesThreshold applies search condition to the query based on the search query, search columns, exists subqueries and threshold.
// It takes the query, search query, search columns, exists subqueries and threshold as parameters.
// It returns the modified query.
// The function splits the search query into individual tokens and generates relevance terms based on the tokens.
// The relevance terms are then used to construct an ORDER BY clause with the SIMILARITY function.
// The SIMILARITY function is used to determine the relevance of each row based on the tokens.
// The rows are then sorted by the relevance score in descending order.
// The search columns and exists subqueries are used to construct the WHERE clause.
// The threshold is used to determine the similarity threshold for the search condition.
// If the search query is empty or the search columns and exists subqueries are empty, the function returns the query unchanged.
// The function returns the modified query.
func ApplySearchConditionWithSubqueriesThreshold(query *gorm.DB, searchQuery string, searchColumns []string, existsSubqueries []string, threshold float64) *gorm.DB {
	if searchQuery == "" || (len(searchColumns) == 0 && len(existsSubqueries) == 0) {
		return query
	}

	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return query
	}

	searchQueryNoSpaces := strings.ReplaceAll(searchQuery, " ", "")
	words = append(words, searchQueryNoSpaces)

	wordClauses := make([]string, 0, len(words))
	args := make([]interface{}, 0)

	reILIKE := ILIKEPatternRegex

	// adjust for search per words.
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		wLower := strings.ToLower(w)
		columnConditions := make([]string, 0, len(searchColumns)+len(existsSubqueries))

		// search column for similarity without case sensitivity and without space with threshold
		for _, column := range searchColumns {
			if !ColumnNameRegex.MatchString(column) {
				continue
			}
			columnConditions = append(columnConditions, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", threshold))
			args = append(args, wLower)
		}

		// 2025/03/04: comment this to decrease scope of search, you can un-comment this to increase scope.
		// for now on default this logic is not really needed.
		// search column for exact match without case sensitivity and without space
		// for _, column := range searchColumns {
		// 	columnConditions = append(columnConditions, "LOWER(REPLACE("+column+", ' ', '')) ILIKE ?")
		// 	wLower = "%" + strings.ToLower(w) + "%"
		// 	args = append(args, wLower)
		// }

		// search column for exact match without case sensitivity
		for _, column := range searchColumns {
			columnConditions = append(columnConditions, "LOWER("+column+") ILIKE ?")
			wLower = "%" + strings.ToLower(w) + "%"
			args = append(args, wLower)
		}

		// sub query logic, if sub query exists, for similarity without case sensitivity and without space with threshold
		for _, existsSubquery := range existsSubqueries {
			modified := reILIKE.ReplaceAllString(existsSubquery, "SIMILARITY(LOWER(REPLACE($1, ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", threshold))
			columnConditions = append(columnConditions, modified)
			args = append(args, wLower)
		}

		// append all column conditions
		if len(columnConditions) > 0 {
			wordClauses = append(wordClauses, "("+strings.Join(columnConditions, " OR ")+")")
		}
	}

	if len(wordClauses) == 0 {
		return query
	}

	whereClause := "(" + strings.Join(wordClauses, " OR ") + ")"
	query = ApplyRelevanceSorting(query, searchQuery, searchColumns)
	return query.Where(whereClause, args...)
}

// ApplyRelevanceSorting applies relevance sorting to the query based on the search query and search columns.
// It takes the query, search query and search columns as parameters.
// It returns the modified query.
// The function splits the search query into individual tokens and generates relevance terms based on the tokens.
// The relevance terms are then used to construct an ORDER BY clause with the SIMILARITY function.
// The SIMILARITY function is used to determine the relevance of each row based on the tokens.
// The rows are then sorted by the relevance score in descending order.
func ApplyRelevanceSorting(query *gorm.DB, searchQuery string, searchColumns []string) *gorm.DB {
	if searchQuery == "" || len(searchColumns) == 0 {
		return query
	}
	tokens := strings.Fields(searchQuery)
	collapsed := strings.ReplaceAll(strings.TrimSpace(searchQuery), " ", "")
	relevanceTerms := make([]string, 0)
	relevanceArgs := make([]interface{}, 0)
	addTerms := func(token string) {
		t := strings.TrimSpace(token)
		if t == "" {
			return
		}
		wNorm := strings.ToLower(strings.ReplaceAll(t, " ", ""))
		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			relevanceTerms = append(relevanceTerms, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), ?)")
			relevanceArgs = append(relevanceArgs, wNorm)
		}
	}
	for _, t := range tokens {
		addTerms(t)
	}
	if collapsed != "" {
		addTerms(collapsed)
	}
	if len(relevanceTerms) == 0 {
		return query
	}
	scoreSQL := "(" + strings.Join(relevanceTerms, " + ") + ") DESC"
	return query.Order(clause.Expr{SQL: scoreSQL, Vars: relevanceArgs})
}

// BuildSearchConditionForRawSQLFromInterface builds a search condition for a raw SQL query based on the search query and the searcher.
// It takes the search query, searcher, startArgIndex and clauseType as parameters.
// It returns the search condition and the arguments.
// If the searcher is nil, the function returns an empty string and an empty args slice.
// The threshold is used to determine the similarity threshold for the search condition.
// If the searcher's SimilarityThreshold() returns a non-nil value, that value is used for the threshold. Otherwise, the default threshold is used.
// The startArgIndex is used to determine the starting index for the query arguments.
// The clauseType is used to determine the type of clause to build. If empty, HAVING is used.
func BuildSearchConditionForRawSQLFromInterface(searchQuery string, searcher NeedSearchPredefine, startArgIndex int, clauseType string) (clause string, args []interface{}) {
	if searcher == nil {
		return "", []interface{}{}
	}
	threshold := DefaultSimilarityThreshold
	if t := searcher.SimilarityThreshold(); t != nil {
		threshold = *t
	}
	return BuildSearchConditionForRawSQLWithThreshold(searchQuery, searcher.GetSearchColumns(), startArgIndex, clauseType, threshold)
}

// BuildSearchConditionForRawSQLWithThreshold builds a search condition for a raw SQL query based on the search query, search columns, startArgIndex, clauseType and threshold.
// It takes the search query, search columns, startArgIndex, clauseType and threshold as parameters.
// It returns the search condition and the arguments.
// If the searcher is nil, the function returns an empty string and an empty args slice.
// The threshold is used to determine the similarity threshold for the search condition.
// If the searcher's SimilarityThreshold() returns a non-nil value, that value is used for the threshold. Otherwise, the default threshold is used.
// The startArgIndex is used to determine the starting index for the query arguments.
// The clauseType is used to determine the type of clause to build. If empty, HAVING is used.
func BuildSearchConditionForRawSQLWithThreshold(searchQuery string, searchColumns []string, startArgIndex int, clauseType string, threshold float64) (clause string, args []interface{}) {
	if searchQuery == "" || len(searchColumns) == 0 {
		return "", []interface{}{}
	}

	if clauseType == "" {
		clauseType = "HAVING"
	}

	if startArgIndex < 1 {
		startArgIndex = 1
	}

	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return "", []interface{}{}
	}

	searchQueryNoSpaces := strings.ReplaceAll(searchQuery, " ", "")
	words = append(words, searchQueryNoSpaces)

	wordConditions := make([]string, 0, len(words))
	currentArgIndex := startArgIndex

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		wLower := strings.ToLower(w)
		columnConditions := make([]string, 0, len(searchColumns))

		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			columnConditions = append(columnConditions, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), $"+fmt.Sprintf("%d", currentArgIndex)+") >= "+fmt.Sprintf("%0.2f", threshold))
			args = append(args, wLower)
			currentArgIndex++
		}

		if len(columnConditions) > 0 {
			wordConditions = append(wordConditions, "("+strings.Join(columnConditions, " OR ")+")")
		}
	}

	if len(wordConditions) == 0 {
		return "", []interface{}{}
	}

	clause = " " + clauseType + " (" + strings.Join(wordConditions, " OR ") + ")"
	return clause, args
}

// SearchPredefineBase provides a reusable Threshold field and method
// that returns a pointer to the Threshold field
type SearchPredefineBase struct {
	Threshold *float64
}

func (b SearchPredefineBase) SimilarityThreshold() *float64 { return b.Threshold }
