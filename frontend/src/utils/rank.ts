export interface Rank {
    role: string;
    minScore: number;
}

export const Ranks: Rank[] = [
    { role: "Junior Agent", minScore: 0 },
    { role: "Field Analyst", minScore: 20000 },
    { role: "Senior Analyst", minScore: 50000 },
    { role: "Intelligence Officer", minScore: 100000 },
    { role: "Division Commander", minScore: 250000 },
    { role: "Atlas Director", minScore: 500000 }
];

export default function CalculateRank( totalScore: number ) {
    
    let currentRank = Ranks[0];
    let nextRank = Ranks[1];

    for (let i = 0; i < Ranks.length; i++) {
        if (totalScore >= Ranks[i].minScore) {
            currentRank = Ranks[i];
            nextRank = Ranks[i+1];
        }
    }

    if (!nextRank) {
        return ({
            nextRank: null,
            xpRemaining: 0,
            progress: 100  
        });
    }

    const xpIntoRank = totalScore - currentRank.minScore;
    const xpToNextRank = nextRank.minScore - currentRank.minScore;

    return ({
        nextRank: nextRank.role,
        xpRemaining: nextRank.minScore - totalScore,
        progress: (xpIntoRank / xpToNextRank) * 100  
    });
}