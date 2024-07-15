package shipyard

import "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"

// Repository is an interface that defines the methods a shipyard repository must implement.
type Repository interface {
	CreateShipyard(data *models.Shipyard) (err error)                                                    // CreateShipyard creates a new shipyard record.
	UpdateShipyard(data *models.Shipyard) (err error)                                                    // UpdateShipyard updates an existing shipyard record.
	DeleteShipyard(data *models.Shipyard) (err error)                                                    // DeleteShipyard deletes a shipyard record.
	FindShipyardByUUIDOrCode(request interface{}) (result models.Shipyard, err error)                    // FindShipyardByUUIDOrCode finds a shipyard record by its UUID or code.
	FetchShipyards(request QueryRequest) (result []*models.Shipyard, total int, lastPage int, err error) // FetchShipyards fetches all shipyard records that match the provided query request.
}

// QueryRequest is an interface that defines the methods a query request must implement.
type QueryRequest interface {
	GetCondition() (result string)                    // GetCondition gets the condition of the query request.
	GetParam() (result []any)                         // GetParam gets the parameters of the query request.
	IsPaginate() (isPaginate bool, limit, offset int) //IsPaginate gets flag to decide the request is using pagination or no and also get the page size and page number is needed
	GetSort() (field, order string)                   // GetSort get field to sort the database result and sort order
}
