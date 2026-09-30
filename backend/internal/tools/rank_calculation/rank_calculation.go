package rank_calculation

func CalculateRank(xp int64) string {
    switch {
        case xp >= 50000:
            return "Atlas Director"
        case xp >= 30000:
            return "Division Commander"
        case xp >= 20000:
            return "Intelligence Officer"
        case xp >= 12500:
            return "Senior Analyst"
        case xp >= 5000:
            return "Field Analyst"
        default:
            return "Junior Agent"
    }
}
