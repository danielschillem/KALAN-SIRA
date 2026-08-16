package payment

import (
 "context"
 "crypto/rand"
 "encoding/hex"
 "fmt"
 "time"

 "github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{db *pgxpool.Pool}
func NewService(db *pgxpool.Pool)*Service{return &Service{db:db}}

type Intent struct{PublicID string `json:"public_id"`;Amount int64 `json:"amount"`;Currency string `json:"currency"`;Provider string `json:"provider"`;Status string `json:"status"`;ExpiresAt time.Time `json:"expires_at"`}
type CreateIntentInput struct{SchoolPublicID string `json:"school_public_id"`;EnrollmentPublicID string `json:"enrollment_public_id"`;GuardianPublicID string `json:"guardian_public_id"`;Amount int64 `json:"amount"`;Provider string `json:"provider"`}
type CashInput struct{SchoolPublicID string `json:"school_public_id"`;EnrollmentPublicID string `json:"enrollment_public_id"`;Amount int64 `json:"amount"`}
type PaymentResult struct{PaymentPublicID string `json:"payment_public_id"`;ReceiptNumber string `json:"receipt_number"`;Amount int64 `json:"amount"`;Allocated int64 `json:"allocated"`;RemainingBalance int64 `json:"remaining_balance"`;Status string `json:"status"`}

func publicID(prefix string)string{b:=make([]byte,6);_,_=rand.Read(b);return prefix+"-"+time.Now().Format("060102")+"-"+hex.EncodeToString(b)}

func(s *Service)CreateIntent(ctx context.Context,in CreateIntentInput)(Intent,error){if in.Amount<=0{return Intent{},fmt.Errorf("amount must be positive")};if in.Provider!="ORANGE_MONEY"&&in.Provider!="MOOV_MONEY"{return Intent{},fmt.Errorf("unsupported provider")};id:=publicID("PI");expires:=time.Now().Add(20*time.Minute);var out Intent
 err:=s.db.QueryRow(ctx,`INSERT INTO payment_intents(public_id,school_id,student_id,enrollment_id,guardian_id,amount,provider,status,expires_at) SELECT $3,e.school_id,e.student_id,e.id,g.id,$4,$5,'CREATED',$6 FROM enrollments e JOIN schools s ON s.id=e.school_id LEFT JOIN guardians g ON g.public_id=NULLIF($2,'') WHERE s.public_id=$1 AND e.public_id=$7 AND e.status='ACTIVE' AND (NULLIF($2,'') IS NULL OR EXISTS(SELECT 1 FROM student_guardians sg WHERE sg.student_id=e.student_id AND sg.guardian_id=g.id AND sg.can_pay=true)) AND $4 <= (SELECT COALESCE(sum(balance),0) FROM student_charges sc WHERE sc.enrollment_id=e.id AND sc.status<>'CANCELLED') RETURNING public_id,amount,currency,provider,status,expires_at`,in.SchoolPublicID,in.GuardianPublicID,id,in.Amount,in.Provider,expires,in.EnrollmentPublicID).Scan(&out.PublicID,&out.Amount,&out.Currency,&out.Provider,&out.Status,&out.ExpiresAt);return out,err}

func(s *Service)RecordCash(ctx context.Context,in CashInput)(PaymentResult,error){if in.Amount<=0{return PaymentResult{},fmt.Errorf("amount must be positive")};tx,err:=s.db.Begin(ctx);if err!=nil{return PaymentResult{},err};defer tx.Rollback(ctx)
 var enrollmentID,schoolID,studentID string;var balance int64
 err=tx.QueryRow(ctx,`SELECT e.id::text,e.school_id::text,e.student_id::text,COALESCE(sum(sc.balance),0) FROM enrollments e JOIN schools s ON s.id=e.school_id JOIN student_charges sc ON sc.enrollment_id=e.id AND sc.status<>'CANCELLED' WHERE s.public_id=$1 AND e.public_id=$2 AND e.status='ACTIVE' GROUP BY e.id`,in.SchoolPublicID,in.EnrollmentPublicID).Scan(&enrollmentID,&schoolID,&studentID,&balance);if err!=nil{return PaymentResult{},err};if in.Amount>balance{return PaymentResult{},fmt.Errorf("payment exceeds remaining balance")}
 pid:=publicID("PAY");var paymentID string;err=tx.QueryRow(ctx,`INSERT INTO payments(public_id,school_id,student_id,enrollment_id,amount,payment_method,status,paid_at) VALUES($1,$2::uuid,$3::uuid,$4::uuid,$5,'CASH','SUCCESS',now()) RETURNING id::text`,pid,schoolID,studentID,enrollmentID,in.Amount).Scan(&paymentID);if err!=nil{return PaymentResult{},err}
 remaining:=in.Amount;rows,err:=tx.Query(ctx,`SELECT id::text,balance FROM student_charges WHERE enrollment_id=$1::uuid AND balance>0 AND status<>'CANCELLED' ORDER BY due_date NULLS LAST,created_at FOR UPDATE`,enrollmentID);if err!=nil{return PaymentResult{},err};type charge struct{id string;balance int64};var charges []charge;for rows.Next(){var c charge;if err=rows.Scan(&c.id,&c.balance);err!=nil{rows.Close();return PaymentResult{},err};charges=append(charges,c)};rows.Close()
 for _,c:=range charges{if remaining==0{break};a:=c.balance;if a>remaining{a=remaining};_,err=tx.Exec(ctx,`INSERT INTO payment_allocations(payment_id,student_charge_id,amount) VALUES($1::uuid,$2::uuid,$3)`,paymentID,c.id,a);if err!=nil{return PaymentResult{},err};_,err=tx.Exec(ctx,`UPDATE student_charges SET amount_paid=amount_paid+$2,balance=balance-$2,status=CASE WHEN balance-$2=0 THEN 'PAID' ELSE 'PARTIAL' END,updated_at=now() WHERE id=$1::uuid`,c.id,a);if err!=nil{return PaymentResult{},err};remaining-=a}
 if remaining!=0{return PaymentResult{},fmt.Errorf("unable to allocate full payment")};receipt:=publicID("REC");token:=publicID("V");_,err=tx.Exec(ctx,`INSERT INTO receipts(public_id,school_id,payment_id,receipt_number,verification_token,issued_at) VALUES($1,$2::uuid,$3::uuid,$1,$4,now())`,receipt,schoolID,paymentID,token);if err!=nil{return PaymentResult{},err};if err=tx.Commit(ctx);err!=nil{return PaymentResult{},err};return PaymentResult{PaymentPublicID:pid,ReceiptNumber:receipt,Amount:in.Amount,Allocated:in.Amount,RemainingBalance:balance-in.Amount,Status:"SUCCESS"},nil}
