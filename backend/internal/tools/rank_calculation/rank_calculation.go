package rank_calculation

func CalculateRank(xp int64) string {
    switch {
        case xp >= 500000:
            return "Atlas Director"
        case xp >= 250000:
            return "Division Commander"
        case xp >= 100000:
            return "Intelligence Officer"
        case xp >= 50000:
            return "Senior Analyst"
        case xp >= 20000:
            return "Field Analyst"
        default:
            return "Junior Agent"
    }
}
