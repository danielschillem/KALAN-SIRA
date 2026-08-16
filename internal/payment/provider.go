package payment

import (
 "context"
 "errors"
)

type ProviderRequest struct { IntentPublicID string; Amount int64; Currency string }
type ProviderResponse struct { ProviderReference string; Status string; Raw []byte }
type Callback struct { Provider string `json:"provider"`; IntentPublicID string `json:"intent_public_id"`; ProviderReference string `json:"provider_reference"`; Status string `json:"status"`; EventID string `json:"event_id"` }

type Provider interface {
 Name() string
 Initiate(context.Context, ProviderRequest) (ProviderResponse,error)
}

var ErrUnsupportedProvider = errors.New("unsupported payment provider")

type Registry struct{ providers map[string]Provider }
func NewRegistry(ps ...Provider)*Registry{r:=&Registry{providers:map[string]Provider{}};for _,p:=range ps{r.providers[p.Name()]=p};return r}
func(r *Registry)Get(name string)(Provider,error){p,ok:=r.providers[name];if !ok{return nil,ErrUnsupportedProvider};return p,nil}
