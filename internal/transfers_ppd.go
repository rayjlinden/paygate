// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"time"

	"github.com/moov-io/ach"
	"github.com/moov-io/base"
)

func createPPDBatch(id, userId string, transfer *Transfer, receiver *Receiver, receiverDep *Depository, orig *Originator, origDep *Depository) (ach.Batcher, error) {
	batchHeader := ach.NewBatchHeader()
	batchHeader.ID = id
	batchHeader.ServiceClassCode = determineServiceClassCode(transfer)
	batchHeader.CompanyName = orig.Metadata
	batchHeader.StandardEntryClassCode = ach.PPD
	batchHeader.CompanyIdentification = orig.Identification
	batchHeader.CompanyEntryDescription = transfer.Description
	batchHeader.CompanyDescriptiveDate = time.Now().Format("060102")
	batchHeader.EffectiveEntryDate = base.Now().AddBankingDay(1).Format("060102") // Date to be posted, YYMMDD
	batchHeader.ODFIIdentification = aba8(origDep.RoutingNumber)

	// Add EntryDetail to PPD batch
	entryDetail := ach.NewEntryDetail()
	entryDetail.ID = id
	entryDetail.TransactionCode = determineTransactionCode(transfer, origDep)
	entryDetail.RDFIIdentification = aba8(receiverDep.RoutingNumber)
	entryDetail.CheckDigit = abaCheckDigit(receiverDep.RoutingNumber)
	entryDetail.DFIAccountNumber = receiverDep.AccountNumber
	entryDetail.Amount = transfer.Amount.Int()
	entryDetail.IdentificationNumber = createIdentificationNumber()
	entryDetail.IndividualName = receiver.Metadata
	entryDetail.DiscretionaryData = transfer.Description
	entryDetail.TraceNumber = createTraceNumber(origDep.RoutingNumber)

	// Add Addenda05
	addenda05 := ach.NewAddenda05()
	addenda05.ID = id
	addenda05.PaymentRelatedInformation = "paygate transaction"
	addenda05.SequenceNumber = 1
	addenda05.EntryDetailSequenceNumber = 1
	entryDetail.AddAddenda05(addenda05)
	entryDetail.AddendaRecordIndicator = 1

	// For now just create PPD
	batch, err := ach.NewBatch(batchHeader)
	if err != nil {
		return nil, fmt.Errorf("ACH file %s (userId=%s): failed to create batch: %v", id, userId, err)
	}
	batch.AddEntry(entryDetail)
	batch.SetControl(ach.NewBatchControl())

	if err := batch.Create(); err != nil {
		return batch, err
	}
	return batch, nil
}
