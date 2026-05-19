package main

import "fmt"

type CampaignSiteTargeting struct {
	CampaignID string
	Sites      []string
}

func GetCampaigns(campaigns []CampaignSiteTargeting, site string) []string {
	mp := map[string][]string{}
	for _, campaign := range campaigns {
		for _, s := range campaign.Sites {
			mp[s] = append(mp[s], campaign.CampaignID)
		}
	}
	return mp[site]
}

func BuildLargeCampaigns(numCampaigns int, sitesPerCampaign int, sitePoolSize int) []CampaignSiteTargeting {
	campaigns := make([]CampaignSiteTargeting, 0, numCampaigns)

	for i := 1; i <= numCampaigns; i++ {
		sites := make([]string, 0, sitesPerCampaign)

		for j := 0; j < sitesPerCampaign; j++ {
			siteNumber := ((i + j) % sitePoolSize) + 1
			sites = append(sites, fmt.Sprintf("site_%05d.com", siteNumber))
		}

		campaigns = append(campaigns, CampaignSiteTargeting{
			CampaignID: fmt.Sprintf("campaign_%06d", i),
			Sites:      sites,
		})
	}

	return campaigns
}

func main() {
	// 100,000 campaigns
	// 50 sites per campaign
	// 10,000 unique sites shared across campaigns
	campaigns := BuildLargeCampaigns(100000, 50, 10000)

	fmt.Println(GetCampaigns(campaigns, "siteB")) // [c1 c2]
	fmt.Println(GetCampaigns(campaigns, "siteA")) // [c1]
	fmt.Println(GetCampaigns(campaigns, "siteX")) // []
}
