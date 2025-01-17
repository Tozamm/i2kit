package route53

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Create creates a CNAME resolving to a service
func Create(s *service.Service, e *environment.Environment, target string) error {
	return update(s, e, target, "UPSERT")
}

//Destroy destroys a CNAME resolving to a service
func Destroy(s *service.Service, e *environment.Environment, target string) error {
	return update(s, e, target, "DELETE")
}

func update(s *service.Service, e *environment.Environment, target, action string) error {
	svc := route53.New(session.New(), e.DNSProvider.GetConfig())
	recordName := fmt.Sprintf("%s.%s", s.GetFullName(e, "."), e.DNSProvider.HostedZone)
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL: aws.Int64(60),
					},
				},
			},
		},
		HostedZoneId: aws.String(e.DNSProvider.HostedZoneID),
	}
	_, err := svc.ChangeResourceRecordSets(params)
	return err
}
