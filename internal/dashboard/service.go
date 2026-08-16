package dashboard

import (
 "context"
 "time"
 "github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{ db *pgxpool.Pool }
func NewService(db *pgxpool.Pool)*Service{return &Service{db:db}}

type KPIs struct{Expected int64 `json:"expected"`;Collected int64 `json:"collected"`;Outstanding int64 `json:"outstanding"`;Overdue int64 `json:"overdue"`;CollectionRate float64 `json:"collection_rate"`;ActiveStudents int64 `json:"active_students"`;StudentsUpToDate int64 `json:"students_up_to_date"`;StudentsOverdue int64 `json:"students_overdue"`}
type Month struct{Month string `json:"month"`;Amount int64 `json:"amount"`}
type Channel struct{Channel string `json:"channel"`;Amount int64 `json:"amount"`}
type OverdueAccount struct{StudentPublicID string `json:"student_public_id"`;StudentName string `json:"student_name"`;ClassName string `json:"class_name"`;Amount int64 `json:"amount"`;DaysLate int `json:"days_late"`}
type RecentPayment struct{PaymentPublicID string `json:"payment_public_id"`;StudentName string `json:"student_name"`;ClassName string `json:"class_name"`;Method string `json:"method"`;Amount int64 `json:"amount"`;PaidAt time.Time `json:"paid_at"`}
type Dashboard struct{SchoolName string `json:"school_name"`;SchoolYear string `json:"school_year"`;KPIs KPIs `json:"kpis"`;Monthly []Month `json:"monthly"`;Channels []Channel `json:"channels"`;OverdueAccounts []OverdueAccount `json:"overdue_accounts"`;RecentPayments []RecentPayment `json:"recent_payments"`}

func(s *Service)Get(ctx context.Context,schoolPublicID string)(Dashboard,error){var o Dashboard
 err:=s.db.QueryRow(ctx,`SELECT name FROM schools WHERE public_id=$1`,schoolPublicID).Scan(&o.SchoolName);if err!=nil{return o,err}
 _=s.db.QueryRow(ctx,`SELECT COALESCE(name,'') FROM school_years sy JOIN schools sc ON sc.id=sy.school_id WHERE sc.public_id=$1 AND sy.is_current=true LIMIT 1`,schoolPublicID).Scan(&o.SchoolYear)
 err=s.db.QueryRow(ctx,`WITH x AS (SELECT sc.enrollment_id,sum(sc.amount_due) expected,sum(sc.amount_paid) collected,sum(sc.balance) outstanding,sum(CASE WHEN sc.due_date<CURRENT_DATE AND sc.balance>0 THEN sc.balance ELSE 0 END) overdue FROM student_charges sc JOIN schools s ON s.id=sc.school_id WHERE s.public_id=$1 AND sc.status<>'CANCELLED' GROUP BY sc.enrollment_id) SELECT COALESCE(sum(expected),0),COALESCE(sum(collected),0),COALESCE(sum(outstanding),0),COALESCE(sum(overdue),0),COALESCE(count(*),0),COALESCE(count(*) FILTER(WHERE outstanding=0),0),COALESCE(count(*) FILTER(WHERE overdue>0),0) FROM x`,schoolPublicID).Scan(&o.KPIs.Expected,&o.KPIs.Collected,&o.KPIs.Outstanding,&o.KPIs.Overdue,&o.KPIs.ActiveStudents,&o.KPIs.StudentsUpToDate,&o.KPIs.StudentsOverdue);if err!=nil{return o,err};if o.KPIs.Expected>0{o.KPIs.CollectionRate=float64(o.KPIs.Collected)*100/float64(o.KPIs.Expected)}
 rows,err:=s.db.Query(ctx,`SELECT to_char(date_trunc('month',p.paid_at),'YYYY-MM'),sum(p.amount) FROM payments p JOIN schools s ON s.id=p.school_id WHERE s.public_id=$1 AND p.status='SUCCESS' AND p.paid_at>=date_trunc('month',CURRENT_DATE)-interval '5 months' GROUP BY 1 ORDER BY 1`,schoolPublicID);if err!=nil{return o,err};defer rows.Close();for rows.Next(){var x Month;if err=rows.Scan(&x.Month,&x.Amount);err!=nil{return o,err};o.Monthly=append(o.Monthly,x)}
 rows2,err:=s.db.Query(ctx,`SELECT payment_method,sum(amount) FROM payments p JOIN schools s ON s.id=p.school_id WHERE s.public_id=$1 AND p.status='SUCCESS' GROUP BY payment_method ORDER BY sum(amount) DESC`,schoolPublicID);if err!=nil{return o,err};defer rows2.Close();for rows2.Next(){var x Channel;if err=rows2.Scan(&x.Channel,&x.Amount);err!=nil{return o,err};o.Channels=append(o.Channels,x)}
 rows3,err:=s.db.Query(ctx,`SELECT st.public_id,concat(st.first_name,' ',st.last_name),COALESCE(c.name,''),sum(sc.balance),GREATEST(0,CURRENT_DATE-min(sc.due_date))::int FROM student_charges sc JOIN schools s ON s.id=sc.school_id JOIN enrollments e ON e.id=sc.enrollment_id JOIN students st ON st.id=e.student_id LEFT JOIN classes c ON c.id=e.class_id WHERE s.public_id=$1 AND sc.balance>0 AND sc.due_date<CURRENT_DATE AND sc.status<>'CANCELLED' GROUP BY st.public_id,st.first_name,st.last_name,c.name ORDER BY sum(sc.balance) DESC LIMIT 5`,schoolPublicID);if err!=nil{return o,err};defer rows3.Close();for rows3.Next(){var x OverdueAccount;if err=rows3.Scan(&x.StudentPublicID,&x.StudentName,&x.ClassName,&x.Amount,&x.DaysLate);err!=nil{return o,err};o.OverdueAccounts=append(o.OverdueAccounts,x)}
 rows4,err:=s.db.Query(ctx,`SELECT p.public_id,concat(st.first_name,' ',st.last_name),COALESCE(c.name,''),p.payment_method,p.amount,p.paid_at FROM payments p JOIN schools s ON s.id=p.school_id JOIN enrollments e ON e.id=p.enrollment_id JOIN students st ON st.id=p.student_id LEFT JOIN classes c ON c.id=e.class_id WHERE s.public_id=$1 AND p.status='SUCCESS' ORDER BY p.paid_at DESC LIMIT 8`,schoolPublicID);if err!=nil{return o,err};defer rows4.Close();for rows4.Next(){var x RecentPayment;if err=rows4.Scan(&x.PaymentPublicID,&x.StudentName,&x.ClassName,&x.Method,&x.Amount,&x.PaidAt);err!=nil{return o,err};o.RecentPayments=append(o.RecentPayments,x)}
 return o,nil
}
