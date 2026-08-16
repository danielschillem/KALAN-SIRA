package admission

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{ db *pgxpool.Pool }

func NewService(db *pgxpool.Pool) *Service { return &Service{db: db} }

type Option struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClassOption struct {
	ID           string `json:"id"`
	PublicID     string `json:"public_id"`
	Name         string `json:"name"`
	LevelID      string `json:"level_id"`
	LevelName    string `json:"level_name"`
	SchoolYearID string `json:"school_year_id"`
}

type FeePreview struct {
	FeeScheduleID    string `json:"fee_schedule_id"`
	FeeScheduleName  string `json:"fee_schedule_name"`
	InstallmentPlanID string `json:"installment_plan_id"`
	InstallmentPlanName string `json:"installment_plan_name"`
	TotalAmount      int64  `json:"total_amount"`
}

type Catalog struct {
	SchoolYears []Option      `json:"school_years"`
	Classes     []ClassOption `json:"classes"`
}

type CreateInput struct {
	SchoolYearID       string `json:"school_year_id"`
	ClassPublicID      string `json:"class_public_id"`
	EnrollmentType     string `json:"enrollment_type"`
	StudentPublicID    string `json:"student_public_id"`
	SchoolStudentNo    string `json:"school_student_no"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	MiddleName         string `json:"middle_name"`
	BirthDate          string `json:"birth_date"`
	BirthPlace         string `json:"birth_place"`
	Gender             string `json:"gender"`
	EnrollmentPublicID string `json:"enrollment_public_id"`
	FeeScheduleID      string `json:"fee_schedule_id"`
	InstallmentPlanID  string `json:"installment_plan_id"`
}

type Result struct {
	StudentPublicID    string `json:"student_public_id"`
	EnrollmentPublicID string `json:"enrollment_public_id"`
	Status             string `json:"status"`
	ChargesCreated     int64  `json:"charges_created"`
	TotalAmount        int64  `json:"total_amount"`
}

func (s *Service) Catalog(ctx context.Context, schoolPublicID string) (Catalog, error) {
	out := Catalog{SchoolYears: []Option{}, Classes: []ClassOption{}}
	rows, err := s.db.Query(ctx, `SELECT sy.id::text,sy.name FROM school_years sy JOIN schools s ON s.id=sy.school_id WHERE s.public_id=$1 AND sy.status IN ('OPEN','ACTIVE') ORDER BY sy.starts_at DESC`, schoolPublicID)
	if err != nil { return out, err }
	for rows.Next() { var v Option; if err=rows.Scan(&v.ID,&v.Name); err!=nil { rows.Close(); return out,err }; out.SchoolYears=append(out.SchoolYears,v) }
	rows.Close()
	rows, err = s.db.Query(ctx, `SELECT c.id::text,c.public_id,c.name,l.id::text,l.name,c.school_year_id::text FROM classes c JOIN levels l ON l.id=c.level_id JOIN schools s ON s.id=c.school_id WHERE s.public_id=$1 AND c.status='ACTIVE' ORDER BY l.display_order,c.name`, schoolPublicID)
	if err != nil { return out, err }
	defer rows.Close()
	for rows.Next() { var v ClassOption; if err=rows.Scan(&v.ID,&v.PublicID,&v.Name,&v.LevelID,&v.LevelName,&v.SchoolYearID); err!=nil{return out,err}; out.Classes=append(out.Classes,v) }
	return out, rows.Err()
}

func (s *Service) Preview(ctx context.Context, schoolPublicID, classPublicID string) (FeePreview, error) {
	var out FeePreview
	err:=s.db.QueryRow(ctx, `SELECT fs.id::text,fs.name,ip.id::text,ip.name,COALESCE((SELECT sum(fi.amount) FROM fee_items fi WHERE fi.fee_schedule_id=fs.id AND fi.mandatory=true),0) FROM classes c JOIN schools s ON s.id=c.school_id JOIN fee_schedules fs ON fs.school_id=s.id AND fs.school_year_id=c.school_year_id AND (fs.class_id=c.id OR (fs.class_id IS NULL AND fs.level_id=c.level_id)) JOIN installment_plans ip ON ip.fee_schedule_id=fs.id WHERE s.public_id=$1 AND c.public_id=$2 AND fs.status IN ('DRAFT','ACTIVE') ORDER BY CASE WHEN fs.class_id=c.id THEN 0 ELSE 1 END,fs.created_at DESC LIMIT 1`,schoolPublicID,classPublicID).Scan(&out.FeeScheduleID,&out.FeeScheduleName,&out.InstallmentPlanID,&out.InstallmentPlanName,&out.TotalAmount)
	return out,err
}

func (s *Service) CreateAndActivate(ctx context.Context, schoolPublicID string, in CreateInput) (Result,error) {
	if strings.TrimSpace(in.StudentPublicID)==""||strings.TrimSpace(in.EnrollmentPublicID)==""||strings.TrimSpace(in.ClassPublicID)=="" { return Result{},fmt.Errorf("student_public_id, enrollment_public_id and class_public_id are required") }
	if strings.TrimSpace(in.FirstName)==""||strings.TrimSpace(in.LastName)=="" { return Result{},fmt.Errorf("first_name and last_name are required") }
	if in.EnrollmentType=="" { in.EnrollmentType="NEW" }
	tx,err:=s.db.Begin(ctx); if err!=nil{return Result{},err}; defer tx.Rollback(ctx)
	var schoolID,studentID,classID,levelID string
	err=tx.QueryRow(ctx,`SELECT s.id::text,c.id::text,c.level_id::text FROM schools s JOIN classes c ON c.school_id=s.id AND c.public_id=$2 AND c.school_year_id=$3::uuid WHERE s.public_id=$1`,schoolPublicID,in.ClassPublicID,in.SchoolYearID).Scan(&schoolID,&classID,&levelID); if err!=nil{return Result{},fmt.Errorf("class not found for school/year: %w",err)}
	err=tx.QueryRow(ctx,`INSERT INTO students(public_id,school_id,school_student_no,first_name,last_name,middle_name,gender,birth_date,birth_place) VALUES($1,$2::uuid,NULLIF($3,''),$4,$5,NULLIF($6,''),NULLIF($7,''),NULLIF($8,'')::date,NULLIF($9,'')) ON CONFLICT(school_id,public_id) DO UPDATE SET school_student_no=COALESCE(NULLIF(EXCLUDED.school_student_no,''),students.school_student_no),first_name=EXCLUDED.first_name,last_name=EXCLUDED.last_name,middle_name=EXCLUDED.middle_name,gender=EXCLUDED.gender,birth_date=EXCLUDED.birth_date,birth_place=EXCLUDED.birth_place,updated_at=now() RETURNING id::text`,in.StudentPublicID,schoolID,in.SchoolStudentNo,strings.TrimSpace(in.FirstName),strings.TrimSpace(in.LastName),in.MiddleName,in.Gender,in.BirthDate,in.BirthPlace).Scan(&studentID); if err!=nil{return Result{},err}
	var enrollmentID string
	err=tx.QueryRow(ctx,`INSERT INTO enrollments(public_id,school_id,school_year_id,student_id,class_id,enrollment_type,status,enrolled_at) VALUES($1,$2::uuid,$3::uuid,$4::uuid,$5::uuid,$6,'PENDING',now()) RETURNING id::text`,in.EnrollmentPublicID,schoolID,in.SchoolYearID,studentID,classID,in.EnrollmentType).Scan(&enrollmentID); if err!=nil{return Result{},err}
	var scheduleSchool, scheduleYear string; var scheduleLevel,scheduleClass *string
	err=tx.QueryRow(ctx,`SELECT school_id::text,school_year_id::text,level_id::text,class_id::text FROM fee_schedules WHERE id=$1::uuid`,in.FeeScheduleID).Scan(&scheduleSchool,&scheduleYear,&scheduleLevel,&scheduleClass); if err!=nil{return Result{},err}
	if scheduleSchool!=schoolID||scheduleYear!=in.SchoolYearID||(scheduleClass!=nil&&*scheduleClass!=classID)||(scheduleClass==nil&&(scheduleLevel==nil||*scheduleLevel!=levelID)){return Result{},fmt.Errorf("fee schedule does not apply to selected class")}
	var planSchedule string; err=tx.QueryRow(ctx,`SELECT fee_schedule_id::text FROM installment_plans WHERE id=$1::uuid`,in.InstallmentPlanID).Scan(&planSchedule); if err!=nil{return Result{},err}; if planSchedule!=in.FeeScheduleID{return Result{},fmt.Errorf("installment plan does not belong to fee schedule")}
	var mandatoryTotal, installmentTotal int64
	_ = tx.QueryRow(ctx,`SELECT COALESCE(sum(amount),0) FROM fee_items WHERE fee_schedule_id=$1::uuid AND mandatory=true`,in.FeeScheduleID).Scan(&mandatoryTotal)
	_ = tx.QueryRow(ctx,`SELECT COALESCE(sum(amount),0) FROM installments WHERE installment_plan_id=$1::uuid`,in.InstallmentPlanID).Scan(&installmentTotal)
	if mandatoryTotal<=0||mandatoryTotal!=installmentTotal{return Result{},fmt.Errorf("invalid fee configuration: mandatory=%d installments=%d",mandatoryTotal,installmentTotal)}
	tag,err:=tx.Exec(ctx,`INSERT INTO student_charges(school_id,student_id,enrollment_id,installment_id,label,original_amount,net_amount,balance,due_date,status) SELECT $1::uuid,$2::uuid,$3::uuid,i.id,i.label,i.amount,i.amount,i.amount,i.due_date,CASE WHEN i.due_date<CURRENT_DATE THEN 'OVERDUE' WHEN i.due_date=CURRENT_DATE THEN 'DUE' ELSE 'UPCOMING' END FROM installments i WHERE i.installment_plan_id=$4::uuid`,schoolID,studentID,enrollmentID,in.InstallmentPlanID); if err!=nil{return Result{},err}; if tag.RowsAffected()==0{return Result{},fmt.Errorf("installment plan has no installments")}
	_,err=tx.Exec(ctx,`UPDATE enrollments SET status='ACTIVE',validated_at=now(),updated_at=now() WHERE id=$1::uuid AND school_id=$2::uuid`,enrollmentID,schoolID); if err!=nil{return Result{},err}
	if err=tx.Commit(ctx);err!=nil{return Result{},err}
	return Result{StudentPublicID:in.StudentPublicID,EnrollmentPublicID:in.EnrollmentPublicID,Status:"ACTIVE",ChargesCreated:tag.RowsAffected(),TotalAmount:installmentTotal},nil
}
