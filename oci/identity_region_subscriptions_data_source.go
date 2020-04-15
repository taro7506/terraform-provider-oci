// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	oci_identity "github.com/oracle/oci-go-sdk/identity"
)

func init() {
	RegisterDatasource("oci_identity_region_subscriptions", IdentityRegionSubscriptionsDataSource())
}

func IdentityRegionSubscriptionsDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readIdentityRegionSubscriptions,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"tenancy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"region_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tenancy_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						// Optional

						// Computed
						"is_home_region": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readIdentityRegionSubscriptions(d *schema.ResourceData, m interface{}) error {
	sync := &IdentityRegionSubscriptionsDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).identityClient

	return ReadResource(sync)
}

type IdentityRegionSubscriptionsDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_identity.IdentityClient
	Res    *oci_identity.ListRegionSubscriptionsResponse
}

func (s *IdentityRegionSubscriptionsDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *IdentityRegionSubscriptionsDataSourceCrud) Get() error {
	request := oci_identity.ListRegionSubscriptionsRequest{}

	if tenancyId, ok := s.D.GetOkExists("tenancy_id"); ok {
		tmp := tenancyId.(string)
		request.TenancyId = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "identity")

	response, err := s.Client.ListRegionSubscriptions(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response
	return nil
}

func (s *IdentityRegionSubscriptionsDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(GenerateDataSourceID())
	resources := []map[string]interface{}{}

	for _, r := range s.Res.Items {
		regionSubscription := map[string]interface{}{}

		if r.IsHomeRegion != nil {
			regionSubscription["is_home_region"] = *r.IsHomeRegion
		}

		if r.RegionKey != nil {
			regionSubscription["region_key"] = *r.RegionKey
		}

		if r.RegionName != nil {
			regionSubscription["region_name"] = *r.RegionName
		}

		regionSubscription["state"] = r.Status

		resources = append(resources, regionSubscription)
	}

	if f, fOk := s.D.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, IdentityRegionSubscriptionsDataSource().Schema["region_subscriptions"].Elem.(*schema.Resource).Schema)
	}

	if err := s.D.Set("region_subscriptions", resources); err != nil {
		return err
	}

	return nil
}
