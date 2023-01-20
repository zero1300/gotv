package resp

type Pager struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}
