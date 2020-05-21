package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.appointy.com/jaal"
	"go.appointy.com/jaal/introspection"
	"go.appointy.com/jaal/schemabuilder"
)

// Server provides our server interface for the API.
type Server struct {
	Experiences []*Experience
}

// Experience defines our schema for the GraphQL queries.
type Experience struct {
	ID                   string
	CompanyName          string
	StartDate            string
	EndDate              string
	JobTitle             string
	JobSummary           string
	JobDescription       string
	SkillsRequired       string
	DetailedAchievements string
	Type                 Type
}

// Type is a 32 bit integer.
type Type int32

// The following constants are the names of the companies for which we will
// fill out experience for.
const (
	IBM Type = iota
	TRIMBLE
	MRSGREENSREMEDIES
	APOLLODAE
	IGADI
	SMARTECH
)

// CreateExperienceRequest provides an interface to generate an experience
// document.
type CreateExperienceRequest struct {
	CompanyName          string
	StartDate            string
	EndDate              string
	JobTitle             string
	JobSummary           string
	JobDescription       string
	SkillsRequired       string
	DetailedAchievements string
	Type                 Type
}

// RegisterPayload initializes the payload for our GraphQL schema.
func RegisterPayload(schema *schemabuilder.Schema) {
	payload := schema.Object("Experience", Experience{})
	payload.FieldFunc("id", func(ctx context.Context, in *Experience) *schemabuilder.ID {
		return &schemabuilder.ID{Value: in.ID}
	})
	payload.FieldFunc("companyname", func(ctx context.Context, in *Experience) string {
		return in.CompanyName
	})
	payload.FieldFunc("type", func(ctx context.Context, in *Experience) Type {
		return in.Type
	})
}

// RegisterInput registers the user's input and generates an experience
// document.
func RegisterInput(schema *schemabuilder.Schema) {
	input := schema.InputObject("CreateExperienceRequest", CreateExperienceRequest{})
	input.FieldFunc("companyname", func(target *Experience, source string) {
		target.CompanyName = source
	})
	input.FieldFunc("type", func(target *Experience, source Type) {
		target.Type = source
	})
}

// RegisterEnum iterates the names of the companies we have experience data for.
func RegisterEnum(schema *schemabuilder.Schema) {
	schema.Enum(Type(0), map[string]interface{}{
		"IBM":       Type(0),
		"TRIMBLE":   Type(0),
		"421HEMP":   Type(0),
		"APOLLODAE": Type(0),
		"IGADI":     Type(0),
		"SMARTECH":  Type(0),
	})
}

// RegisterOperations registers the commands and actions which can be performed
// on our experience documents.
func (s *Server) RegisterOperations(schema *schemabuilder.Schema) {
	schema.Query().FieldFunc("experience", func(ctx context.Context, args struct {
		ID *schemabuilder.ID
	}) *Experience {
		for _, ex := range s.Experiences {
			if ex.ID == args.ID.Value {
				return ex
			}
		}
		return nil
	})

	schema.Query().FieldFunc("experiences", func(ctx context.Context, args struct{}) []*Experience {
		return s.Experiences
	})

	schema.Mutation().FieldFunc("createExperience", func(ctx context.Context, args struct {
		Input *CreateExperienceRequest
	}) *Experience {
		ex := &Experience{
			ID:                   uuid.Must(uuid.NewUUID()).String(),
			CompanyName:          args.Input.CompanyName,
			StartDate:            args.Input.StartDate,
			EndDate:              args.Input.EndDate,
			JobTitle:             args.Input.JobTitle,
			JobSummary:           args.Input.JobSummary,
			JobDescription:       args.Input.JobDescription,
			SkillsRequired:       args.Input.SkillsRequired,
			DetailedAchievements: args.Input.DetailedAchievements,
		}
		s.Experiences = append(s.Experiences, ex)

		return ex
	})
}

func run() {
	sb := schemabuilder.NewSchema()
	RegisterPayload(sb)
	RegisterInput(sb)
	RegisterEnum(sb)

	s := &Server{
		Experiences: []*Experience{{
			ID:                   "1",
			CompanyName:          "IBM",
			StartDate:            "2017-07-01",
			EndDate:              "2018-03-01",
			JobTitle:             "Global Technology Services: Solutions, Del, & Transf.",
			JobSummary:           "Anthem Healthcare z/OS Mainframe Administrator and State Street Investment Bank Help Desk Admin",
			JobDescription:       "Level 1.5 CTS Agent, administrator rights to remote into State Street virtual desktops and physical machines after initial troubleshooting as failed on other teams. Mainly services the U.S. and India. Support Active Directory and configuration of multiple server types. Support IBM mainframe z/OS systems. Write troubleshooting and DevOps documentation and SLA status templates. Deploy and maintain Windows 10 and Ent. Server 2012-16 Deploy Microsoft Office 365 via local source and configuration manager. Development of applications for z/OS utilizing various databases",
			SkillsRequired:       "z/OS, MVS, JCL, SDSF, VSAM, Endevor, SyncSort, MQSeries, SQL, SPUFI, FileAid, Xpediter",
			DetailedAchievements: "Master the Mainframe Part 2 Competitor",
		}},
	}

	s.RegisterOperations(sb)
	schema, err := sb.Build()
	if err != nil {
		log.Fatalln(err)
	}

	introspection.AddIntrospectionToSchema(schema)

	http.Handle("/graphql", jaal.HTTPHandler(schema))
	log.Println("Running")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic(err)
	}
}
