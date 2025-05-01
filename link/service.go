package link

import "github.com/pgvanniekerk/ezapp/internal/primitive"

func Service[Service primitive.Service[Params], Params any]() {}
