package tokenstore


// Service represents the repository where tokens are stored / retrieved.
type Service struct {
	TokenManager
}

// NewService create a Service to store / retrieve tokens.
func NewService(r TokenManager) Service {
	return Service{r}
}
