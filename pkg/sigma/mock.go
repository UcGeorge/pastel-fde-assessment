package sigma

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// MockClient is a fake Sigma API client that generates realistic random
// responses. It implements the SigmaClient interface and can be swapped in
// via dependency injection to test the UI without hitting the live API.
type MockClient struct{}

// Compile-time assertion.
var _ SigmaClient = (*MockClient)(nil)

// NewMockClient creates a new MockClient instance.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Helpers

func pick[T any](items ...T) T           { return items[rand.Intn(len(items))] }
func randBetween(min, max int) int       { return min + rand.Intn(max-min+1) }
func randFloat(min, max float64) float64 { return min + rand.Float64()*(max-min) }
func randID() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", rand.Int63(), rand.Intn(0xffff), rand.Intn(0xffff), rand.Intn(0xffff), rand.Int63())
}
func ptr[T any](v T) *T { return &v }

var (
	firstNames = []string{"Vladimir", "Abdulrahman", "Chen", "Maria", "James", "Fatima", "Sergei", "Olga", "Hassan", "Amina", "Boris", "Elena", "Yusuf", "Aisha", "Dmitry"}
	lastNames  = []string{"Putin", "Okonkwo", "Wei", "Garcia", "Smith", "Al-Rashid", "Ivanov", "Petrova", "Mohammed", "Diallo", "Volkov", "Kuznetsova", "Ibrahim", "Bello", "Sorokin"}
	countries  = []string{"Russia", "Nigeria", "China", "United States", "Iran", "North Korea", "Syria", "Venezuela", "Cuba", "Belarus", "Myanmar", "Zimbabwe", "Libya", "Sudan", "Yemen"}
	sanctions  = []string{"OFAC SDN List", "UN Security Council", "EU Consolidated List", "UK HMT Sanctions", "DFAT Consolidated List", "SECO Sanctions", "Japan MOF List"}
	parties    = []string{"United Russia", "APC", "CPC", "Democratic Party", "PSUV", "Workers Party", "Baath Party", "ZANU-PF"}
	positions  = []string{"President", "Prime Minister", "Minister of Finance", "Head of State Security", "Central Bank Governor", "Ambassador", "Senator", "Governor", "Chief of Staff", "Defence Minister"}
	datasets   = []Dataset{
		{Name: "OpenSanctions", URL: "https://opensanctions.org"},
		{Name: "World-Check", URL: "https://www.refinitiv.com/en/products/world-check"},
		{Name: "Dow Jones Risk", URL: "https://www.dowjones.com/professional/risk"},
		{Name: "OFAC SDN", URL: "https://sanctionssearch.ofac.treas.gov"},
		{Name: "EU Sanctions Map", URL: "https://sanctionsmap.eu"},
	}
	ruleNames = []string{
		"High-Value Transaction Alert",
		"Unusual Transaction Pattern",
		"Cross-Border Transfer Limit",
		"Velocity Check Exceeded",
		"New Account Large Transfer",
		"Dormant Account Reactivation",
		"Structuring Detection",
		"Geographic Risk Assessment",
	}
	severities    = []SeverityLevel{SeverityLow, SeverityMedium, SeverityHigh}
	entityTypes   = []string{"person", "organization"}
	genders       = []string{"male", "female"}
	mediaStatuses = []string{"pending", "processing", "completed"}
)

// SubmitTransaction

func (m *MockClient) SubmitTransaction(_ context.Context, req *SubmitTransactionRequest) (*TransactionResponse, error) {
	severity := pick(severities...)
	var action ActionType
	var code ActionCode
	switch severity {
	case SeverityHigh:
		action, code = ActionDeclined, ActionCodeRejected
	default:
		action, code = ActionApproved, ActionCodeApproved
	}

	riskScore := fmt.Sprintf("%.1f", randFloat(0, 100))
	ruleID := fmt.Sprintf("rule_%04d", randBetween(1, 9999))

	return &TransactionResponse{
		Message: "Transaction processed successfully",
		Data: &TransactionResponseData{
			TransactionID: fmt.Sprintf("txn_%s", randID()),
			RiskScore:     riskScore,
			Action: TransactionAction{
				Result: action,
				Code:   code,
			},
			RuleResult: pick("triggered", "passed", "flagged"),
			Reason: TransactionReason{
				Code:     fmt.Sprintf("RC-%03d", randBetween(100, 999)),
				Message:  fmt.Sprintf("Transaction flagged: %s on reference %s", pick(ruleNames...), req.TransactionData.Reference),
				Severity: severity,
				Rule: TransactionRule{
					ID:   ruleID,
					Name: pick(ruleNames...),
				},
			},
			Screening: ScreeningResult{
				Sender: ScreeningParty{
					PEP:      pick("clear", "match_found", "pending"),
					Sanction: pick("clear", "match_found", "pending"),
				},
				Receiver: ScreeningParty{
					PEP:      pick("clear", "match_found", "pending"),
					Sanction: pick("clear", "match_found", "pending"),
				},
			},
		},
	}, nil
}

// CheckPEP / CheckSanction

func (m *MockClient) CheckPEP(_ context.Context, req *ScreeningRequest) (*ScreeningResponse, error) {
	return m.mockScreening(req, "PEP"), nil
}

func (m *MockClient) CheckSanction(_ context.Context, req *ScreeningRequest) (*ScreeningResponse, error) {
	return m.mockScreening(req, "Sanction"), nil
}

func (m *MockClient) mockScreening(req *ScreeningRequest, screenType string) *ScreeningResponse {
	count := randBetween(1, 5)
	entities := make([]MatchedEntity, count)

	for i := range entities {
		numAliases := randBetween(1, 4)
		aliases := make([]string, numAliases)
		for j := range aliases {
			aliases[j] = fmt.Sprintf("%s %s", pick(firstNames...), pick(lastNames...))
		}

		numPositions := randBetween(1, 3)
		positionList := make([]Position, numPositions)
		for j := range positionList {
			startYear := randBetween(1990, 2020)
			positionList[j] = Position{
				ID:        randID(),
				Title:     pick(positions...),
				StartDate: ptr(fmt.Sprintf("%d-01-01", startYear)),
				EndDate:   ptr(fmt.Sprintf("%d-12-31", startYear+randBetween(1, 15))),
			}
		}

		numSanctions := 0
		var sanctionList []string
		if screenType == "Sanction" {
			numSanctions = randBetween(1, 3)
			sanctionList = make([]string, numSanctions)
			for j := range sanctionList {
				sanctionList[j] = pick(sanctions...)
			}
		}

		numDatasets := randBetween(1, 3)
		datasetList := make([]Dataset, numDatasets)
		for j := range datasetList {
			datasetList[j] = pick(datasets...)
		}

		numCountries := randBetween(1, 3)
		countryList := make([]string, numCountries)
		for j := range countryList {
			countryList[j] = pick(countries...)
		}

		birthYear := randBetween(1940, 1990)

		entities[i] = MatchedEntity{
			ID:              randID(),
			Type:            strings.ToLower(screenType),
			ExternalID:      fmt.Sprintf("ext-%s", randID()[:8]),
			Addresses:       []string{fmt.Sprintf("%d %s Street, %s", randBetween(1, 500), pick("Main", "Oak", "Kremlin", "Palace", "Government"), pick(countries...))},
			Aliases:         aliases,
			BirthDate:       fmt.Sprintf("%d-%02d-%02d", birthYear, randBetween(1, 12), randBetween(1, 28)),
			Countries:       countryList,
			CreatedAt:       time.Now().Add(-time.Duration(randBetween(1, 365*5)) * 24 * time.Hour).Format(time.RFC3339),
			UpdatedAt:       time.Now().Add(-time.Duration(randBetween(1, 30)) * 24 * time.Hour).Format(time.RFC3339),
			Dataset:         datasetList,
			EntityType:      pick(entityTypes...),
			Name:            fmt.Sprintf("%s %s", pick(firstNames...), pick(lastNames...)),
			Sanctions:       sanctionList,
			PoliticalParty:  []string{pick(parties...)},
			Gender:          pick(genders...),
			Positions:       positionList,
			Tags:            []string{pick("high-risk", "pep", "sanctioned", "watchlist", "monitored")},
			SearchScore:     randFloat(0.3, 1.0),
			ConfidenceScore: randFloat(req.Threshold, 1.0),
		}
	}

	return &ScreeningResponse{
		Message:    fmt.Sprintf("%s screening completed successfully", screenType),
		Count:      count,
		PageNumber: 1,
		Data:       entities,
	}
}

// CheckAdverseMedia

func (m *MockClient) CheckAdverseMedia(_ context.Context, req *AdverseMediaRequest) (*AdverseMediaResponse, error) {
	now := time.Now().UTC()
	return &AdverseMediaResponse{
		Message: "Adverse media check request submitted successfully",
		Data: &AdverseMediaData{
			ID:               randID(),
			Query:            req.Query,
			BusinessProfile:  randID()[:8],
			Status:           "completed",
			CreatedAt:        now.Format(time.RFC3339),
			UpdatedAt:        now.Add(time.Duration(randBetween(1, 60)) * time.Second).Format(time.RFC3339),
			FindingsReturned: ptr(true),
			RiskCategory:     ptr(pick("High", "Medium", "Low")),
			Sources:          []string{pick("Global News", "Local Media", "Financial Times"), pick("Reuters", "Bloomberg")},
		},
	}, nil
}
