package exec

import (
	"fmt"
	"time"
	DRS "github.com/fbaube/datarepo/sqlite"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Brings in global var L
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// SimpleRepo is in env.
func ImportBatchIntoDB(pSR *DRS.SqliteRepo, InputContentities []*mcfile.Contentity) error {

	var err, e error
	L.L.Progress("Starting import batch...")
	// =====================
	//  START A TRANSACTION
	// =====================
	err = pSR.Begin()
	if err != nil {
		L.L.Error("Exec.BeginTx failed: %w", err)
	}
	L.L.Info("TRANSACTION IS STARTED")
	var timeNow = time.Now().UTC().Format(time.RFC3339)
	// ===============================
	//    FIRST THE INBATCH
	//   So the batch number can be
	//  plugged into the Contentities 
	// ===============================
	pIB := new(DRM.InbatchRow)
	pIB.FilCt = len(InputContentities)
	pIB.Descr = "CLI import"
	// pIB.RelFP =
	// pIB.AnsFP =
	pIB.T_Cre, pIB.T_Imp = timeNow, timeNow
	var newInbatchID int
	newInbatchID, e = DRS.DoInsertGeneric(pSR, pIB)
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
		// L.L.Info("doImport.L57: Trying new INSERT Generic")
		var newCtyID int
		newCtyID, e = DRS.DoInsertGeneric(pSR, &pMCF.ContentityRow)
		if e != nil {
			return mcfile.WrapAsContentityError(e,
				"insert contentity to DB (cli.exec)", pMCF)
		}
		L.L.Info("Added file to import batch, ID: %d", newCtyID)
	}
	e = pSR.Commit()
	if e != nil {
		return mcfile.WrapAsContentityError(e,
			"commit txn to DB failed (cli.exec)", nil)
	}
	L.L.Okay("Batch imported OK: TRANSACTION SUCCEEDED")
	// pp := pCA.SimpleRepo.GetFileAll()
	// fmt.Printf("    DD:Files len %d id[last] %d \n", len(pp), fileIndex)

	L.L.Info("TRYING SELECT BY ID")
	stmtS, eS := pSR.NewSelectByIdStmt(&DRM.TableDetailsCNT, 1)
	if eS != nil {
		return fmt.Errorf("new select contentity by id=1 stmt (cli.exec): %w", eS)
	}
	result, e3 := DRS.ExecSelectOneStmt[*DRM.ContentityRow](pSR, stmtS)
	if e3 != nil {
		return fmt.Errorf("new select contentity by id=1 from DB (cli.exec): %w", e3)
	}
	L.L.Warning("exec.go: INSERT'd inbatch OK, ID:%d", result)
	return nil
}