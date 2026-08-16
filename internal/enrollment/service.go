package enrollment

import (
 "context"
 "fmt"
 "strings"

 "github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{db *pgxpool.Pool}
func NewService(db *pgxpool.Pool)*Service{return &Service{db:db}}
type Enrollment struct{ID string `json:"id"`;PublicID string `json:"public_id"`;SchoolID string `json:"school_id"`;SchoolYearID string `json:"school_year_id"`;StudentID string `json:"student_id"`;ClassID string `json:"class_id"`;EnrollmentType string `json:"enrollment_type"`;Status string `json:"status"`}
type CreateInput struct{PublicID string `json:"public_id"`;SchoolPublicID string `json:"school_public_id"`;SchoolYearID string `json:"school_year_id"`;StudentPublicID string `json:"student_public_id"`;ClassPublicID string `json:"class_public_id"`;EnrollmentType string `json:"enrollment_type"`}
func(s *Service)Create(ctx context.Context,in CreateInput)(Enrollment,error){if strings.TrimSpace(in.PublicID)==""{return Enrollment{},fmt.Errorf("public_id is required")};if in.EnrollmentType==""{in.EnrollmentType="NEW"};var out Enrollment;err:=s.db.QueryRow(ctx,`INSERT INTO enrollments(public_id,school_id,school_year_id,student_id,class_id,enrollment_type,status,enrolled_at) SELECT $2,s.id,sy.id,st.id,c.id,$6,'PENDING',now() FROM schools s JOIN school_years sy ON sy.id=$3::uuid AND sy.school_id=s.id JOIN students st ON st.public_id=$4 AND st.school_id=s.id JOIN classes c ON c.public_id=$5 AND c.school_id=s.id AND c.school_year_id=sy.id WHERE s.public_id=$1 RETURNING id::text,public_id,school_id::text,school_year_id::text,student_id::text,class_id::text,enrollment_type,status`,in.SchoolPublicID,in.PublicID,in.SchoolYearID,in.StudentPublicID,in.ClassPublicID,in.EnrollmentType).Scan(&out.ID,&out.PublicID,&out.SchoolID,&out.SchoolYearID,&out.StudentID,&out.ClassID,&out.EnrollmentType,&out.Status);return out,err}
