package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	DRS "github.com/fbaube/datarepo/sqlite"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Brings in global var L
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// SimpleRepo is in env.
func ImportBatchIntoDB(env *XmlAppEnv, InputContentities []*mcfile.Contentity) error {

     var jout []byte
     var jerr, e error
     
	// ===============================
	//   4. LOAD FILES INTO DATABASE
	//      See also above
	// ===============================

	var haveDB bool = (env.SimpleRepo != nil)
	L.L.Dbg("env.SimpleRepo: <%T> %#v", env.SimpleRepo, env.SimpleRepo)
	// var jout []byte
	// var jerr error
	jout, jerr = json.MarshalIndent(env.SimpleRepo, "SimpleRepo: ", "  ")
	if jerr != nil {
		println(jerr)
		panic(jerr)
	}
	L.L.Info("JSON! " + string(jout))
	println("JSON! " + string(jout))

	var pSR *DRS.SqliteRepo
	var ok bool
	pSR, ok = env.SimpleRepo.(*DRS.SqliteRepo)
	if !ok {
		panic("Exec: repo is not *SimpleSqliteRepo")
	}
	if pSR == nil {
		L.L.Error("Cannot proceed: SqliteRepo is not valid")
		os.Exit(1)
	}
	// if haveDB && env.cfg.b.DBdoImport {
	if env.cfg.b.DBdoImport {
	   	var batchIndex int
		if !haveDB {
			panic("Exec: doImport but not haveDB")
		}
		var err error
		L.L.Progress("Will loadicontent for import batch, ID:%d",
			batchIndex)
		// var pTx *sql.Tx
		// =====================
		//  START A TRANSACTION
		// =====================
		err = env.SimpleRepo.Begin()
		if err != nil {
			L.L.Error("cli.exec.BeginTx: %w", err)
		}
		L.L.Info("TRANSACTION IS STARTED")
		var timeNow = time.Now().UTC().Format(time.RFC3339)
		// ============================
		//  FOR EVERY Input Contentity
		// ============================
		for _, pMCF := range InputContentities {
			// Prepare a DB record for the File
			pMCF.T_Imp = timeNow
			// L.L.Info("exec.L463: Trying new INSERT Generic")
			/*
				var stmt string
				// OBS stmt, e = pSR.NewInsertStmt(&pMCF.ContentityRow)
				stmt, e = DRS.NewInsertStmtGnrcFunc(
				      pSR, &pMCF.ContentityRow)
				if e != nil {
					return mcfile.WrapAsContentityError(e,
					  "new insert contentity stmt (cli.exec)", pMCF)
				}
				var insID int
				insID, e = pSR.ExecInsertStmt(stmt)
			*/
			var insID int
			insID, e = DRS.DoInsertGeneric(pSR, &pMCF.ContentityRow)
			if e != nil {
				return mcfile.WrapAsContentityError(e,
					"insert contentity to DB (cli.exec)", pMCF)
			}
			L.L.Info("Added file to import batch, ID: %d", insID)
		}
		pIB := new(DRM.InbatchRow)
		pIB.FilCt = len(InputContentities)
		pIB.Descr = "CLI import"
		// pIB.RelFP =
		// pIB.AnsFP =
		pIB.T_Cre, pIB.T_Imp = timeNow, timeNow
		/*
			var stmt string
			stmt, e = pSR.NewInsertStmt(pIB)
			if e != nil {
				return fmt.Errorf("new insert inbatch stmt (cli.exec): %w", e)
			}
			var insID int
			insID, e = pSR.ExecInsertStmt(stmt)
		*/
		var insID int
		insID, e = DRS.DoInsertGeneric(pSR, pIB)

		if e != nil {
			return fmt.Errorf("new insert inbatch to DB (cli.exec): %w", e)
		}
		L.L.Okay("cli/exec: INSERT'd inbatch OK, ID: %d", insID)

		e = pSR.Commit()
		if e != nil {
			return mcfile.WrapAsContentityError(e,
				"commit txn to DB failed (cli.exec)", nil)
		}
		L.L.Okay("Batch imported OK: TRANSACTION SUCCEEDED")
		// env.SimpleRepo.Tx = nil
		// pp := pCA.SimpleRepo.GetFileAll()
		// fmt.Printf("    DD:Files len %d id[last] %d \n", len(pp), fileIndex)
	}
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