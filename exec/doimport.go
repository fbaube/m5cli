package exec

import (
	"fmt"
	"time"
	"database/sql"
	DRS "github.com/fbaube/datarepo/sqlite"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Brings in global var L
	"github.com/fbaube/m5db"
)

// SimpleRepo is stored in struct [appenv].
func ImportBatchIntoDB(pSR *DRS.SqliteRepo, InputContentities []*mcfile.Contentity) error {

	var err, e error
	L.L.Info("Starting import batch...")
	// =====================
	//  START A TRANSACTION
	// =====================
	var Tx *sql.Tx
	Tx, err = pSR.Begin()
	if err != nil {
		L.L.Error("Exec.Begin.Tx failed: %w", err)
	}
	L.L.Info("TRANSACTION IS STARTED")
	var timeNow = time.Now().UTC().Format(time.RFC3339)
	// ===============================
	//    FIRST THE INBATCH
	//   So the batch number can be
	//  plugged into the Contentities 
	// ===============================
	pINB := new(m5db.InbatchRow)
	pINB.FilCt = len(InputContentities)
	pINB.Descr = "CLI import"
	// pINB.RelFP =
	// pINB.AnsFP =
	pINB.T_Cre, pINB.T_Imp = timeNow, timeNow
	var newInbatchID int
	newInbatchID, e = DRS.DoInsertGeneric(pSR, pINB)
	if e != nil {
		return fmt.Errorf("Exec.DoImport.inbatch failed: %w", e)
	}
	L.L.Okay("Exec.DoImport.inbatch: OK, ID: %d", newInbatchID)
	// ======================
	//  NOW THE CONTENTITIES
	// ======================
	for _, pMCF := range InputContentities {
		// Prepare a DB record for the File
		pMCF.T_Imp = timeNow
		pMCF.Idx_Inbatch = newInbatchID
		// L.L.Info("Exec.DoImport.L50: Trying new INSERT Generic")
		var newCtyID int
		newCtyID, e = DRS.DoInsertGeneric(pSR, &pMCF.ContentityRow)
		if e != nil {
			return mcfile.WrapAsContentityError(e,
				"Exec.DoImport.InsCty", pMCF)
		}
		L.L.Info("Added file to import batch, ID: %d", newCtyID)
	}
	// =====================
	//  END THE TRANSACTION
	// =====================
	e = Tx.Commit() // pSR.Commit()
	
	if e != nil {
		return mcfile.WrapAsContentityError(e,
			"commit txn to DB failed (cli.exec)", nil)
	}
	L.L.Okay("Batch imported OK: TRANSACTION SUCCEEDED")
	L.L.Okay("Exec.DoImport: insert'ed inbatch OK, ID:%d", newInbatchID)

	var wasFound bool
	wasFound, e = DRS.DoSelectByIdGeneric(
		  pSR, newInbatchID, new(m5db.InbatchRow))
	L.L.Info("Found the just-added new Inbatch?: %t", wasFound)
	
	return nil
}