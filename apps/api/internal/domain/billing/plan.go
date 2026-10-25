package billing

import (
	"fmt"

	"uvico-api/internal/domain/shared"
)

type PlanID string

const (
	PlanFounder      PlanID = "founder"
	PlanEssential    PlanID = "essential"
	PlanProfessional PlanID = "professional"
	PlanEnterprise   PlanID = "enterprise"
)

type BillingInterval string

const (
	IntervalMonthly BillingInterval = "monthly"
	IntervalAnnual  BillingInterval = "annual"
)

// Plan es inmutable y se construye desde config. Nunca se muta en runtime.
type Plan struct {
	ID                      PlanID
	MonthlyCents            int64
	AnnualCents             int64
	TrialDays               int
	FounderPromoMonths      int
	FounderPromoMonthlyCents int64
}

var (
	ErrInvalidPlan     = fmt.Errorf("billing: plan inválido")
	ErrInvalidInterval = fmt.Errorf("billing: intervalo inválido")
)

// Plans construye el catálogo inmutable de planes desde config.
// El 17% de descuento anual se calcula como Monthly * 12 * bps / 10000
// y se aplica floor (siempre a favor del cliente).
func Plans(founderMonthly, founderPromo, essentialMonthly, professionalMonthly, enterpriseMonthly, annualBps int64) map[PlanID]Plan {
	return map[PlanID]Plan{
		PlanFounder: {
			ID:                       PlanFounder,
			MonthlyCents:             founderMonthly,
			AnnualCents:              annual(founderMonthly, 0, annualBps),
			TrialDays:                7,
			FounderPromoMonths:       3,
			FounderPromoMonthlyCents: founderPromo,
		},
		PlanEssential: {
			ID:           PlanEssential,
			MonthlyCents: essentialMonthly,
			AnnualCents:  annual(essentialMonthly, 0, annualBps),
			TrialDays:    7,
		},
		PlanProfessional: {
			ID:           PlanProfessional,
			MonthlyCents: professionalMonthly,
			AnnualCents:  annual(professionalMonthly, 0, annualBps),
			TrialDays:    7,
		},
		PlanEnterprise: {
			ID:           PlanEnterprise,
			MonthlyCents: enterpriseMonthly,
			AnnualCents:  annual(enterpriseMonthly, 0, annualBps),
			TrialDays:    7,
		},
	}
}

func annual(monthly, _ int64, bps int64) int64 {
	return (monthly * 12 * bps) / 10000
}

func (p Plan) PriceFor(interval BillingInterval) (shared.Money, error) {
	switch interval {
	case IntervalMonthly:
		return shared.USD(p.MonthlyCents), nil
	case IntervalAnnual:
		return shared.USD(p.AnnualCents), nil
	default:
		return shared.Money{}, fmt.Errorf("%w: %q", ErrInvalidInterval, interval)
	}
}