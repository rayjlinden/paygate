// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package internal

import (
	"testing"

	"github.com/moov-io/base"
)

func TestCCD__createCCDBatch(t *testing.T) {
	id, userID := base.ID(), base.ID()
	receiverDep := &Depository{
		ID:            DepositoryID(base.ID()),
		BankName:      "foo bank",
		Holder:        "jane doe",
		HolderType:    Individual,
		Type:          Checking,
		RoutingNumber: "121042882",
		AccountNumber: "2",
		Status:        DepositoryVerified,
		Metadata:      "jane doe checking",
	}
	receiver := &Receiver{
		ID:                ReceiverID(base.ID()),
		Email:             "jane.doe@example.com",
		DefaultDepository: receiverDep.ID,
		Status:            ReceiverVerified,
		Metadata:          "jane doe",
	}
	origDep := &Depository{
		ID:            DepositoryID(base.ID()),
		BankName:      "foo bank",
		Holder:        "john doe",
		HolderType:    Individual,
		Type:          Savings,
		RoutingNumber: "231380104",
		AccountNumber: "2",
		Status:        DepositoryVerified,
		Metadata:      "john doe savings",
	}
	orig := &Originator{
		ID:                OriginatorID(base.ID()),
		DefaultDepository: origDep.ID,
		Identification:    "dddd",
		Metadata:          "john doe",
	}
	amt, _ := NewAmount("USD", "100.00")
	transfer := &Transfer{
		ID:                     TransferID(base.ID()),
		Type:                   PushTransfer,
		Amount:                 *amt,
		Originator:             orig.ID,
		OriginatorDepository:   origDep.ID,
		Receiver:               receiver.ID,
		ReceiverDepository:     receiverDep.ID,
		Description:            "sending money",
		StandardEntryClassCode: "CCD",
		Status:                 TransferPending,
		CCDDetail: &CCDDetail{
			PaymentInformation: "test payment",
		},
	}

	batch, err := createCCDBatch(id, userID, transfer, receiver, receiverDep, orig, origDep)
	if err != nil {
		t.Fatal(err)
	}
	if batch == nil {
		t.Error("nil CCD Batch")
	}

	file, err := constructACHFile(id, "", userID, transfer, receiver, receiverDep, orig, origDep)
	if err != nil {
		t.Fatal(err)
	}
	if file == nil {
		t.Error("nil CCD ach.File")
	}

	// sad path, empty CCDDetail.PaymentInformation
	transfer.CCDDetail.PaymentInformation = ""
	batch, err = createCCDBatch(id, userID, transfer, receiver, receiverDep, orig, origDep)
	if err == nil || batch != nil {
		t.Fatalf("expected error: batch=%#v", batch)
	}
}
