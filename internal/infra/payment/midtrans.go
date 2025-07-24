package payment

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func New(serverKey string) (snap.Client, coreapi.Client) {
	var snap snap.Client
	snap.New(serverKey, midtrans.Sandbox)

	var coreapi coreapi.Client
	coreapi.New(serverKey, midtrans.Sandbox)

	return snap, coreapi
}
