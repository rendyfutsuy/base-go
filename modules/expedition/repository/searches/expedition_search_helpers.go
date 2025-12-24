package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type ExpeditionSearchHelper struct{ request.SearchPredefineBase }

func (ExpeditionSearchHelper) GetSearchColumns() []string {
	return []string{
		"e.expedition_code",
		"e.expedition_name",
		"e.address",
	}
}

func (ExpeditionSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{
		"EXISTS (SELECT 1 FROM expedition_contacts ec WHERE ec.expedition_id = e.id AND ec.deleted_at IS NULL AND ec.phone_number ILIKE ?)",
	}
}

var _ request.NeedSearchPredefine = ExpeditionSearchHelper{}

func NewExpeditionSearchHelper() ExpeditionSearchHelper {
	return ExpeditionSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: nil}}
}
