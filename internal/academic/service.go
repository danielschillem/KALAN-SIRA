package academic

import (
 "context"
 "fmt"
 "strings"
 "time"

 "github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{ db *pgxpool.Pool }
func NewService(db *pgxpool.Pool) *Service { return &Service{db: db} }

type SchoolYear struct { ID string `json:"id"`; SchoolID string `json:"school_id"`; Name string `json:"name"`; StartsAt time.Time `json:"starts_at"`; EndsAt time.Time `json:"ends_at"`; Status string `json:"status"` }
type CreateSchoolYearInput struct { SchoolPublicID string `json:"school_public_id"`; Name string `json:"name"`; StartsAt string `json:"starts_at"`; EndsAt string `json:"ends_at"` }

type Level struct { ID string `json:"id"`; SchoolID string `json:"school_id"`; Code string `json:"code"`; Name string `json:"name"`; Cycle string `json:"cycle,omitempty"`; DisplayOrder int `json:"display_order"` }
type CreateLevelInput struct { SchoolPublicID string `json:"school_public_id"`; Code string `json:"code"`; Name string `json:"name"`; Cycle string `json:"cycle"`; DisplayOrder int `json:"display_order"` }

type Class struct { ID string `json:"id"`; PublicID string `json:"public_id"`; SchoolID string `json:"school_id"`; SchoolYearID string `json:"school_year_id"`; LevelID string `json:"level_id"`; Name string `json:"name"`; Capacity *int `json:"capacity,omitempty"`; Status string `json:"status"` }
type CreateClassInput struct { PublicID string `json:"public_id"`; SchoolPublicID string `json:"school_public_id"`; SchoolYearID string `json:"school_year_id"`; LevelID string `json:"level_id"`; Name string `json:"name"`; Capacity *int `json:"capacity"` }

func (s *Service) CreateSchoolYear(ctx context.Context, in CreateSchoolYearInput) (SchoolYear,error) {
 start,err:=time.Parse("2006-01-02",in.StartsAt); if err!=nil{return SchoolYear{},fmt.Errorf("invalid starts_at")}; end,err:=time.Parse("2006-01-02",in.EndsAt); if err!=nil{return SchoolYear{},fmt.Errorf("invalid ends_at")}; if !end.After(start){return SchoolYear{},fmt.Errorf("ends_at must be after starts_at")}
 var out SchoolYear
 err=s.db.QueryRow(ctx,`INSERT INTO school_years(school_id,name,starts_at,ends_at) SELECT id,$2,$3,$4 FROM schools WHERE public_id=$1 RETURNING id::text,school_id::text,name,starts_at,ends_at,status`,in.SchoolPublicID,strings.TrimSpace(in.Name),start,end).Scan(&out.ID,&out.SchoolID,&out.Name,&out.StartsAt,&out.EndsAt,&out.Status)
 return out,err
}

func (s *Service) CreateLevel(ctx context.Context,in CreateLevelInput)(Level,error){
 if strings.TrimSpace(in.Code)==""||strings.TrimSpace(in.Name)==""{return Level{},fmt.Errorf("code and name are required")}; var out Level
 err:=s.db.QueryRow(ctx,`INSERT INTO levels(school_id,code,name,cycle,display_order) SELECT id,$2,$3,NULLIF($4,''),$5 FROM schools WHERE public_id=$1 RETURNING id::text,school_id::text,code,name,COALESCE(cycle,''),display_order`,in.SchoolPublicID,strings.TrimSpace(in.Code),strings.TrimSpace(in.Name),strings.TrimSpace(in.Cycle),in.DisplayOrder).Scan(&out.ID,&out.SchoolID,&out.Code,&out.Name,&out.Cycle,&out.DisplayOrder); return out,err
}

func (s *Service) CreateClass(ctx context.Context,in CreateClassInput)(Class,error){
 if strings.TrimSpace(in.PublicID)==""||strings.TrimSpace(in.Name)==""{return Class{},fmt.Errorf("public_id and name are required")}; var out Class
 err:=s.db.QueryRow(ctx,`INSERT INTO classes(public_id,school_id,school_year_id,level_id,name,capacity) SELECT $2,s.id,$3::uuid,$4::uuid,$5,$6 FROM schools s JOIN school_years sy ON sy.id=$3::uuid AND sy.school_id=s.id JOIN levels l ON l.id=$4::uuid AND l.school_id=s.id WHERE s.public_id=$1 RETURNING id::text,public_id,school_id::text,school_year_id::text,level_id::text,name,capacity,status`,in.SchoolPublicID,in.PublicID,in.SchoolYearID,in.LevelID,strings.TrimSpace(in.Name),in.Capacity).Scan(&out.ID,&out.PublicID,&out.SchoolID,&out.SchoolYearID,&out.LevelID,&out.Name,&out.Capacity,&out.Status); return out,err
}
