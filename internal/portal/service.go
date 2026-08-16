package portal

import(
 "context"
 "crypto/rand"
 "crypto/sha256"
 "encoding/base64"
 "encoding/hex"
 "fmt"
 "time"
 "github.com/jackc/pgx/v5/pgxpool"
)
type Service struct{db *pgxpool.Pool}
func NewService(db *pgxpool.Pool)*Service{return &Service{db:db}}
type CreateLinkInput struct{SchoolPublicID string `json:"school_public_id"`;StudentPublicID string `json:"student_public_id"`;GuardianPublicID string `json:"guardian_public_id"`;Amount int64 `json:"amount"`;ExpiresInHours int `json:"expires_in_hours"`}
type Link struct{PublicID string `json:"public_id"`;Token string `json:"token,omitempty"`;Amount int64 `json:"amount"`;ExpiresAt time.Time `json:"expires_at"`;Status string `json:"status"`}
type PaymentPage struct{SchoolName string `json:"school_name"`;SchoolPublicID string `json:"school_public_id"`;StudentPublicID string `json:"student_public_id"`;StudentName string `json:"student_name"`;GuardianName string `json:"guardian_name,omitempty"`;Amount int64 `json:"amount"`;Currency string `json:"currency"`;ExpiresAt time.Time `json:"expires_at"`;Status string `json:"status"`;Providers []string `json:"providers"`}
func token()string{b:=make([]byte,32);_,_=rand.Read(b);return base64.RawURLEncoding.EncodeToString(b)}
func hash(v string)string{x:=sha256.Sum256([]byte(v));return hex.EncodeToString(x[:])}
func id()string{b:=make([]byte,6);_,_=rand.Read(b);return "PL-"+time.Now().Format("060102")+"-"+hex.EncodeToString(b)}
func(s *Service)CreateLink(ctx context.Context,in CreateLinkInput)(Link,error){if in.Amount<=0{return Link{},fmt.Errorf("amount must be positive")};if in.ExpiresInHours<=0{in.ExpiresInHours=72};raw:=token();expires:=time.Now().Add(time.Duration(in.ExpiresInHours)*time.Hour);public:=id();var out Link;e:=s.db.QueryRow(ctx,`INSERT INTO payment_links(public_id,school_id,student_id,guardian_id,token_hash,amount,expires_at,status) SELECT $4,s.id,st.id,g.id,$5,$6,$7,'ACTIVE' FROM schools s JOIN students st ON st.school_id=s.id AND st.public_id=$2 LEFT JOIN guardians g ON g.public_id=NULLIF($3,'') WHERE s.public_id=$1 AND (NULLIF($3,'') IS NULL OR EXISTS(SELECT 1 FROM student_guardians sg WHERE sg.student_id=st.id AND sg.guardian_id=g.id AND sg.can_pay=true)) AND $6 <= (SELECT COALESCE(sum(balance),0) FROM student_charges sc JOIN enrollments e ON e.id=sc.enrollment_id WHERE e.student_id=st.id AND e.status='ACTIVE' AND sc.status<>'CANCELLED') RETURNING public_id,amount,expires_at,status`,in.SchoolPublicID,in.StudentPublicID,in.GuardianPublicID,public,hash(raw),in.Amount,expires).Scan(&out.PublicID,&out.Amount,&out.ExpiresAt,&out.Status);if e!=nil{return Link{},e};out.Token=raw;return out,nil}
func(s *Service)GetPage(ctx context.Context,raw string)(PaymentPage,error){var o PaymentPage;e:=s.db.QueryRow(ctx,`SELECT s.name,s.public_id,st.public_id,concat_ws(' ',st.first_name,st.last_name),COALESCE(concat_ws(' ',g.first_name,g.last_name),''),pl.amount,'XOF',pl.expires_at,CASE WHEN pl.status='ACTIVE' AND pl.expires_at<=now() THEN 'EXPIRED' ELSE pl.status END FROM payment_links pl JOIN schools s ON s.id=pl.school_id JOIN students st ON st.id=pl.student_id LEFT JOIN guardians g ON g.id=pl.guardian_id WHERE pl.token_hash=$1`,hash(raw)).Scan(&o.SchoolName,&o.SchoolPublicID,&o.StudentPublicID,&o.StudentName,&o.GuardianName,&o.Amount,&o.Currency,&o.ExpiresAt,&o.Status);if e!=nil{return PaymentPage{},e};o.Providers=[]string{"ORANGE_MONEY","MOOV_MONEY"};return o,nil}
