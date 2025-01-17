package resolvers

import (
	"artion-api-graphql/internal/repository"
	"artion-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// Token object is constructed from query, data from db are loaded on demand into "loaded" field.
type Token struct {
	Address common.Address
	TokenId hexutil.Big
	dbToken *types.Token // data for token loaded from Mongo
	jsonMeta *types.JsonMetadata // data loaded from URI
}

type TokenEdge struct {
	Node *Token
}

func (edge TokenEdge) Cursor() (types.Cursor, error) {
	return types.CursorFromId(types.TokenIdFromAddress(&edge.Node.Address, (*big.Int)(&edge.Node.TokenId))), nil
}

type TokenConnection struct {
	Edges      []TokenEdge
	TotalCount hexutil.Big
	PageInfo   PageInfo
}

func NewTokenConnection(list *types.TokenList) (con *TokenConnection, err error) {
	con = new(TokenConnection)
	con.TotalCount = (hexutil.Big)(*big.NewInt(list.TotalCount))
	con.Edges = make([]TokenEdge, len(list.Collection))
	for i := 0; i < len(list.Collection); i++ {
		resolverToken := Token{
			Address: list.Collection[i].Nft,
			TokenId: list.Collection[i].TokenId,
			dbToken: list.Collection[i],
		}
		con.Edges[i].Node = &resolverToken
	}
	con.PageInfo.HasNextPage = list.HasNext
	con.PageInfo.HasPreviousPage = list.HasPrev
	if len(list.Collection) > 0 {
		startCur, err := con.Edges[0].Cursor()
		if err != nil {
			return nil, err
		}
		endCur, err := con.Edges[len(con.Edges)-1].Cursor()
		if err != nil {
			return nil, err
		}
		con.PageInfo.StartCursor = &startCur
		con.PageInfo.EndCursor = &endCur
	}
	return con, err
}

func (t *Token) load() error {
	if t.dbToken == nil {
		tok, err := repository.R().GetToken(t.Address, t.TokenId)
		if err != nil {
			return err
		}
		t.dbToken = tok
	}
	if t.jsonMeta == nil {
		jsonMeta, err := repository.R().GetTokenJsonMetadata(t.dbToken.Uri)
		if err != nil {
			return err
		}
		t.jsonMeta = jsonMeta
	}
	return nil
}

func (t Token) Name() (string, error) {
	err := t.load()
	if err != nil {
		return "", err
	}
	return t.jsonMeta.Name, nil
}

func (t Token) Description() (string, error) {
	err := t.load()
	if err != nil {
		return "", err
	}
	return t.jsonMeta.Description, nil
}

func (t Token) Image() (*string, error) {
	err := t.load()
	if err != nil {
		return nil, err
	}
	return t.jsonMeta.Image, nil
}

func (t Token) ImageProxy() (*string, error) {
	err := t.load()
	if err != nil || t.jsonMeta.Image == nil {
		return nil, err
	}
	url := "/token-image/" + t.Address.String() + "/" + t.TokenId.String()
	return &url, nil
}

func (t Token) Events(args struct{ PaginationInput }) (con *TokenEventConnection, err error) {
	cursor, count, backward, err := args.ToRepositoryInput()
	if err != nil {
		return nil, err
	}
	list, err := repository.R().ListTokenEvents(&t.Address, &t.TokenId, nil, cursor, count, backward)
	if err != nil {
		return nil, err
	}
	return NewTokenEventConnection(list)
}

func (t Token) Listings(args struct{ PaginationInput }) (con *ListingConnection, err error) {
	cursor, count, backward, err := args.ToRepositoryInput()
	if err != nil {
		return nil, err
	}
	list, err := repository.R().ListListings(&t.Address, &t.TokenId, nil, cursor, count, backward)
	if err != nil {
		return nil, err
	}
	return NewListingConnection(list)
}

func (t Token) Offers(args struct{ PaginationInput }) (con *OfferConnection, err error) {
	cursor, count, backward, err := args.ToRepositoryInput()
	if err != nil {
		return nil, err
	}
	list, err := repository.R().ListOffers(&t.Address, &t.TokenId, nil, cursor, count, backward)
	if err != nil {
		return nil, err
	}
	return NewOfferConnection(list)
}
