package student

import (
 "context"
 "fmt"
 "strings"
 "time"

 "github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{ db *pgxpool.Pool }
func NewService(db *pgxpool.Pool)*Service{return &Service{db:db}}

type Student struct{ ID string `json:"id"`; PublicID string `json:"public_id"`; SchoolID string `json:"school_id"`; SchoolStudentNo string `json:"school_student_no,omitempty"`; FirstName string `json:"first_name"`; LastName string `json:"last_name"`; MiddleName string `json:"middle_name,omitempty"`; BirthDate *time.Time `json:"birth_date,omitempty"`; Status string `json:"status"` }
type CreateStudentInput struct{ PublicID string `json:"public_id"`; SchoolPublicID string `json:"school_public_id"`; SchoolStudentNo string `json:"school_student_no"`; FirstName string `json:"first_name"`; LastName string `json:"last_name"`; MiddleName string `json:"middle_name"`; BirthDate string `json:"birth_date"`; BirthPlace string `json:"birth_place"`; Gender string `json:"gender"` }

type Guardian struct{ ID string `json:"id"`; PublicID string `json:"public_id"`; FirstName string `json:"first_name"`; LastName string `json:"last_name"`; Phone string `json:"phone"`; Email string `json:"email,omitempty"` }
type CreateGuardianInput struct{ PublicID string `json:"public_id"`; FirstName string `json:"first_name"`; LastName string `json:"last_name"`; Phone string `json:"phone"`; Email string `json:"email"` }
type LinkGuardianInput struct{ SchoolPublicID string `json:"school_public_id"`; StudentPublicID string `json:"student_public_id"`; GuardianPublicID string `json:"guardian_public_id"`; Relationship string `json:"relationship"`; IsPrimary bool `json:"is_primary"`; CanPay bool `json:"can_pay"`; CanReceiveNotifications bool `json:"can_receive_notifications"` }

func(s *Service)CreateStudent(ctx context.Context,in CreateStudentInput)(Student,error){
 if strings.TrimSpace(in.PublicID)==""||strings.TrimSpace(in.FirstName)==""||strings.TrimSpace(in.LastName)==""{return Student{},fmt.Errorf("public_id, first_name and last_name are required")}; var birth any
 if in.BirthDate!=""{v,err:=time.Parse("2006-01-02",in.BirthDate);if err!=nil{return Student{},fmt.Errorf("invalid birth_date")};birth=v}
 var out Student
 err:=s.db.QueryRow(ctx,`INSERT INTO students(public_id,school_id,school_student_no,first_name,last_name,middle_name,gender,birth_date,birth_place) SELECT $2,id,NULLIF($3,''),$4,$5,NULLIF($6,''),NULLIF($7,''),$8,NULLIF($9,'') FROM schools WHERE public_id=$1 RETURNING id::text,public_id,school_id::text,COALESCE(school_student_no,''),first_name,last_name,COALESCE(middle_name,''),birth_date,status`,in.SchoolPublicID,in.PublicID,in.SchoolStudentNo,strings.TrimSpace(in.FirstName),strings.TrimSpace(in.LastName),in.MiddleName,in.Gender,birth,in.BirthPlace).Scan(&out.ID,&out.PublicID,&out.SchoolID,&out.SchoolStudentNo,&out.FirstName,&out.LastName,&out.MiddleName,&out.BirthDate,&out.Status);return out,err
}
func(s *Service)CreateGuardian(ctx context.Context,in CreateGuardianInput)(Guardian,error){if in.PublicID==""||in.FirstName==""||in.LastName==""||in.Phone==""{return Guardian{},fmt.Errorf("public_id, first_name, last_name and phone are required")};var out Guardian;err:=s.db.QueryRow(ctx,`INSERT INTO guardians(public_id,first_name,last_name,phone,email) VALUES($1,$2,$3,$4,NULLIF($5,'')) RETURNING id::text,public_id,first_name,last_name,phone,COALESCE(email,'')`,in.PublicID,in.FirstName,in.LastName,in.Phone,in.Email).Scan(&out.ID,&out.PublicID,&out.FirstName,&out.LastName,&out.Phone,&out.Email);return out,err}
func(s *Service)LinkGuardian(ctx context.Context,in LinkGuardianInput)error{tag,err:=s.db.Exec(ctx,`INSERT INTO student_guardians(student_id,guardian_id,relationship,is_primary,can_pay,can_receive_notifications) SELECT st.id,g.id,$4,$5,$6,$7 FROM students st JOIN schools s ON s.id=st.school_id JOIN guardians g ON g.public_id=$3 WHERE s.public_id=$1 AND st.public_id=$2 ON CONFLICT(student_id,guardian_id) DO UPDATE SET relationship=EXCLUDED.relationship,is_primary=EXCLUDED.is_primary,can_pay=EXCLUDED.can_pay,can_receive_notifications=EXCLUDED.can_receive_notifications`,in.SchoolPublicID,in.StudentPublicID,in.GuardianPublicID,in.Relationship,in.IsPrimary,in.CanPay,in.CanReceiveNotifications);if err!=nil{return err};if tag.RowsAffected()==0{return fmt.Errorf("student or guardian not found")};return nil}
