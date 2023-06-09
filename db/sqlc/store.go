package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute queries and transactions
// which Queries cannot
type Store struct {
	*Queries // composition to expend a struct
	db       *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// start a db transaction, create new Queries obj with that transaction
// call the callback func with the created Queries and commit or rollback
// this methods cannot be exported because it starts with lower e
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	// default isolation lvl is Read Commited
	// BeginTx return a transaction obj or err
	if err != nil {
		return err
	}
	// call New func with the transaction
	// New() accept DBTX both DB and TX as param
	query := New(tx)
	err = fn(query)
	if err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			// Calling Errorf() function with verb %v which is used
			// for printing structs
			return fmt.Errorf("tx error: %v, toll back err: %v", err, rollBackErr)
		}
		return err
	}
	return tx.Commit()
}

// thực hiện 1 giao dịch CK sẽ cần 5 bước
// tạo giao dịch chuyển tiền VD $10
// tạo entry với -$10 của tk chuyển tiền đi
// tạo entry +$10 vs tk được nhận
// cập nhật balance của tk chuyển
// cập nhật balance của tk nhận

// TransferTxParams sẽ đưa ra những param cần thiết để thực hiện chuyển khoản.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// The TransferTxResult that need to return
// in order to know the result of the transaction
// kết quả của 5 bước chuyển khoản bên trên
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`   // see the whole transfer
	FromAccount Account  `json:from_account` //see the balance of from account
	ToAccount   Account  `json:"to_account"` //see the balance of to account
	FromEntry   Entry    `json:"from_entry"` // money moving out
	ToEntry     Entry    `json:"to_entry"`   //money in
}

// TransferTx sẽ làm nhiệm vụ như 1 transaction khi cta chuyển tiền
// tạo ra record những giao dịch (transfer), tạo entry trừ và cộng tiền, update số dư (balance)
// trong 1 transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		// bên trong này là 1 closure
		// chúng ta đang access result và arg
		// nắm ngoài scope của func để sử dụng
		// go không support generics type nên là callback func
		// giúp cta biết đc kiểu trả về
		var err error
		// tạo giao dịch (transfer) lưu lại quá trình chuyển tiền
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		// Ta có CreateTransderParams vs TransferTxParams có y hệt các fields
		// vậy nên có thể làm ntn: CreateTransferParams(arg)

		// trả về kqua err
		// nếu như transfer không thành công
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // tk chuyển phải trừ tiền
		})
		// tạo from entry không thành công
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		// tạo from entry không thành công
		if err != nil {
			return err
		}
		//TODO: update the balance of 2 account
		//! require locking and preventing deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}
		// trả về kết quả của hàm vô danh
		return nil
	})
	return result, err
}
func AddMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // automaticly return all the result n
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return // automaticly return all the result n
	}
	return
}
