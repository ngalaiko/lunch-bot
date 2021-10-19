package handler

import (
	"context"
	"sort"
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/response"

	"github.com/aws/aws-lambda-go/events"
)

func handleList(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	placeChances, err := roller.ListChances(ctx, time.Now())
	if err != nil {
		return response.InternalServerError(err)
	}

	type placeChance struct {
		Name   places.Name
		Chance float64
	}

	pp := make([]placeChance, 0, len(placeChances))
	for place, chance := range placeChances {
		pp = append(pp, placeChance{
			Name:   place,
			Chance: chance,
		})
	}

	sort.SliceStable(pp, func(i, j int) bool {
		return pp[i].Chance > pp[j].Chance
	})

	blocks := []*response.Block{
		response.Section(nil, response.Markdown("*Title*"), response.Markdown("*Odds*")),
		response.Divider(),
	}

	for _, p := range pp {
		blocks = append(blocks, response.SectionFields(
			[]*response.TextBlock{
				response.PlainText("%s", p.Name),
				response.PlainText("%.2f%%", p.Chance),
			},
			response.WithButton(response.PlainText("Boost"), "boost", string(p.Name)),
		))
	}

	return response.Ephemral(blocks...)
}
